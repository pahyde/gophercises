package main

import (
    "fmt"
    "os"
    "net/http"
    "container/list"
    "net/url"
    "errors"
    "strings"
    "flag"
    "encoding/xml"

    "sitemap/link"
)

func main() {
    root  := flag.String("url", "https://calhoun.io", "root url to start building the sitemap")
    depth := flag.Int("depth", 0, "Maximum depth (links traversed) while building the sitemap")
    flag.Parse()

    s := NewSiteMap(*depth)
    s.Bfs(*root)
    urls, err := s.ToXmlUrlSet()
    if err != nil {
        exit(err)
    }
    if err := os.WriteFile("sitemap.xml", urls, 0644); err != nil {
        exit(err)
    }
}

func NewQueue[T any](values ...T) *Queue[T] {
    l := list.New()
    q := &Queue[T]{l}
    for _, v := range values {
        q.Enqueue(v)
    }
    return q
}

type Queue[T any] struct {
    l  *list.List
}

func (q *Queue[T]) Len() int {
    return q.l.Len()
}

func (q *Queue[T]) Empty() bool {
    return q.Len() == 0
}

func (q *Queue[T]) Enqueue(v T) {
    q.l.PushBack(v)
}

func (q *Queue[T]) Dequeue() T {
    e := q.l.Front().Value.(T)
    q.l.Remove(q.l.Front())
    return e
}


type SiteMap struct {
    MaxDepth int
    UrlMap   map[string]bool
}

func (s *SiteMap) Contains(u string) bool {
    return s.UrlMap[u]
}

func (s *SiteMap) Add(u string) {
    s.UrlMap[u] = true
}


// XML urlset data structure
type UrlSet struct {
    XMLName  xml.Name `xml:"urlset"`
    Urls     []Url    `xml:"url"`
}
type Url struct {
    XMLName  xml.Name `xml:"url"`
    Loc      string   `xml:"loc"`
}

func (s *SiteMap) ToXmlUrlSet() ([]byte, error) {
    var set UrlSet
    set.Urls = make([]Url, 0, len(s.UrlMap))
    for l, _ := range s.UrlMap {
        set.Urls = append(set.Urls, Url{Loc: l})
    }
    b, err := xml.Marshal(&set)
    if err != nil {
        return nil, err
    }
    return append([]byte(xml.Header), b...), nil
}

func NewSiteMap(d int) *SiteMap {
    return &SiteMap{
        MaxDepth: d,
        UrlMap: make(map[string]bool),
    }
}

type SiteVertex struct {
    Edges  []string
}

func NewSiteVertex(urlstr string) (*SiteVertex, error) {
    u, err := url.Parse(urlstr)
    if err != nil {
        return nil, err
    }
    links, err := getLinksFromUrl(urlstr)
    if err != nil {
        return nil, err
    }
    var edges []string
    for _, l := range links {
        e, err := makeEdge(l.Href, u);
        if err != nil {
            // l is not a valid internal site link
            continue
        }
        edges = append(edges, e)
    }
    return &SiteVertex{edges}, nil
}

// error if href has invalid format or points to external host
func makeEdge(href string, from *url.URL) (string, error) {
    to, err := url.Parse(href)
    if err != nil {
        return "", err
    }
    if to.Path == "" {
        return "", errors.New("recursive path")
    }
    // absolute path: match
    if strings.TrimPrefix(to.Host, "www.") == strings.TrimPrefix(from.Host, "www.") {
        return to.String(), nil
    }
    // absolute path: mismatch
    if to.Host != "" {
        return "", errors.New("external host")
    }
    to.Scheme = from.Scheme
    to.Host   = from.Host
    // root relative path
    if to.Path[0] == '/' {
        return to.String(), nil
    }
    // page relative path
    to.Path, err = url.JoinPath(from.Path, to.Path)
    if err != nil {
        return "", nil
    }
    return to.String(), nil
}

func getLinksFromUrl(urlstr string) ([]link.Link , error) {
    r, err := http.Get(urlstr)
    if err != nil {
        return nil, err
    }
    defer func() {
        if err := r.Body.Close(); err != nil {
            exit(err)
        }
    }()
    links, err := link.Parse(r.Body)
    if err != nil {
        return nil, err
    }
    return links, nil
}

// breadth-first-search links from given root url -> SiteMap s
func (s *SiteMap) Bfs(root string) {
    // enqueue and add first url to site map
    q := NewQueue[string](root)
    s.Add(root)
    // add linked pages up to max depth
    for d := 0; d < s.MaxDepth && !q.Empty(); d++ {
        qLen := q.l.Len()
        // iterate pages at current depth
        for i := 0; i < qLen; i++ {
            // construct site vertex corresponding to dequeued URL
            // valid internal links (expressed as absolute paths) are availble by v.Edges
            v, err := NewSiteVertex(q.Dequeue())
            if err != nil {
                exit(err)
            }
            // if we haven't seen a page before, enqueue it and add to the site map
            for _, e := range v.Edges {
                if !s.Contains(e) {
                    s.Add(e)
                    q.Enqueue(e)
                }
            }
        }
    }
}

func check(err error) {
    if err != nil {
        exit(err)
    }
}

func exit(m any) {
    fmt.Println(m)
    os.Exit(1)
}

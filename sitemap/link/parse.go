package link

import (
    "fmt"
    "strings"
    "golang.org/x/net/html"
    "io"
)

type Link struct {
    Href string
    Text string
}

func Parse(r io.Reader) ([]Link, error) {
    root, err := html.Parse(r)
    if err != nil {
        return nil, err
    }

    // depth-first-search for anchor tags
    anchors := dfsAnchors(root)

    // get Link data from anchors
    links := make([]Link, 0)
    for _, a := range anchors {
        link, err := parseLink(a)
        if err != nil {
            // ignore anchors without href attr
            continue
        }
        links = append(links, link)
    }
    return links, nil
}

func dfsAnchors(n *html.Node) []*html.Node {
    anchors := make([]*html.Node, 0)
    var f func(n *html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            anchors = append(anchors, n)
            // ignore nested anchors
            return
        } 
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(n)
    return anchors
}

func parseLink(n *html.Node) (Link, error) {
    text := dfsText(n)
    if href, ok := getHref(n); ok {
        return Link{href, text}, nil
    }
    return Link{"", text}, fmt.Errorf("No href attr for anchor")
}

func getHref(a *html.Node) (string, bool) {
    for _, attr := range a.Attr {
        if attr.Key == "href" {
            return attr.Val, true
        }
    }
    return "", false
}

func dfsText(n *html.Node) string {
    chunks := make([]string, 0)
    var f func(n *html.Node)
    f = func(n *html.Node) {
        if n.Type == html.TextNode {
            chunk := strings.TrimSpace(n.Data)
            chunks = append(chunks, chunk)
        } 
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(n)
    return strings.Join(chunks, " ")
}


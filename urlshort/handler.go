package urlshort

import (
	"net/http"
	"net/url"
	"log"
	"gopkg.in/yaml.v3"
	"encoding/json"
	"encoding/xml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        urlstr, ok := pathsToUrls[r.URL.Path]
        if !ok {
            // short url not found
            fallback.ServeHTTP(w,r)
            return
        }
        dest, err := url.Parse(urlstr)
        if err != nil {
            // error parsing long url
            log.Printf("error parsing long url %s. possibly invalid format.\n", urlstr)
            fallback.ServeHTTP(w,r)
            return
        }
        // update url host and path, preserve all else (query params etc)
        r.URL.Host = dest.Host
        r.URL.Path = dest.Path
        http.Redirect(w, r, r.URL.String(), 302)
    }
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

// supported serialization formats
type fileType int
const (
    jsonf   fileType = iota
    yamlf
    xmlf
)

// PathToUrl: data structure representing a short path to long url kv-pair
// yaml, json, and xml formats must have underlying []PathToUrl structure
type PathToUrl struct {
    Path string
    Url  string
}

// PathToUrl wrapper for unmarshaling xml
type PathsToUrlsXML struct {
    List []PathToUrl `xml:"PathToUrl"`
}

func parseYAML(yml []byte) ([]PathToUrl, error) {
    var l []PathToUrl
    if err := yaml.Unmarshal(yml, &l); err != nil {
        return l, err
    }
    return l, nil
}

func parseJSON(jsn []byte) ([]PathToUrl, error) {
    var l []PathToUrl
    if err := json.Unmarshal(jsn, &l); err != nil {
        return l, err
    }
    return l, nil
}

func parseXML(xm []byte) ([]PathToUrl, error) {
    var l PathsToUrlsXML
    if err := xml.Unmarshal(xm, &l); err != nil {
        return nil, err
    }
    return l.List, nil
}

func buildMap(l []PathToUrl) map[string]string {
    m := make(map[string]string)
    for _, entry := range l {
        m[entry.Path] = entry.Url
    }
    return m
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
    return convertToMapHandler(yml, fallback, yamlf)
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
    return convertToMapHandler(jsn, fallback, jsonf)
}

func XMLHandler(xm []byte, fallback http.Handler) (http.HandlerFunc, error) {
    return convertToMapHandler(xm, fallback, xmlf)
}

func convertToMapHandler(dat []byte, fallback http.Handler, f fileType) (http.HandlerFunc, error) {
    var parsed []PathToUrl
    var err error
    switch f {
    case yamlf:
        parsed, err = parseYAML(dat)
    case jsonf:
        parsed, err = parseJSON(dat)
    case xmlf:
        parsed, err = parseXML(dat)
    }
    if err != nil {
        return nil, err
    }
    pathMap := buildMap(parsed)
    log.Println(pathMap)
    return MapHandler(pathMap, fallback), nil
}

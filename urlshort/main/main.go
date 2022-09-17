package main

import (
	"fmt"
	"os"
	"net/http"
	"path/filepath"
	"flag"

	"urlshort"
)

func main() {
	mux := defaultMux()

    // file flag: Allows user to submit addition path-url redirection pairs
    // in either yaml, json, or xml format
    m := `Input yaml, json, or xml file denoting (path-url)
redirection pairs for url shortener`
    inputFile := flag.String("file", "", m)
    flag.Parse()

    // yaml, json, and xml default values
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
    json := `
[
    {
        "path": "/fall",
        "url": "https://youtu.be/QhBoFEq0lU0"
    },
    {
        "path": "/feels",
        "url": "https://soundcloud.com/makzo/shanti"
    }
]
    `
    xml := `
<PathsToUrls>
    <PathToUrl>
        <Path>/my-gophercises</Path>
        <Url>https://github.com/pahyde/gophercises</Url>
    </PathToUrl>
    <PathToUrl>
        <Path>/my-urlshort</Path>
        <Url>https://github.com/pahyde/gophercises/urlshort</Url>
    </PathToUrl>
</PathsToUrls>
    `

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)


	// Build the XMLHandler -> JSONHandler -> YAMLHandler -> mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
    jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
	    fmt.Println(json)
		panic(err)
	}
	xmlHandler, err := urlshort.XMLHandler([]byte(xml), jsonHandler)
	if err != nil {
		panic(err)
	}

    entryPoint := xmlHandler
    if *inputFile != "" {
        h, err := getHandlerFromFile(*inputFile, xmlHandler)
        if err != nil {
            exit(err)
        }
        entryPoint = h
    }

	// start the server
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", entryPoint)
}

func getHandlerFromFile(filename string, fallback http.Handler) (http.HandlerFunc, error) {
    b, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("Unable to read file %s", filename)
    }
    ext := filepath.Ext(filename)
    switch ext {
    case ".yaml":
        return urlshort.YAMLHandler(b, fallback)
    case ".json":
        return urlshort.JSONHandler(b, fallback)
    case ".xml":
        return urlshort.XMLHandler(b, fallback)
    default:
        return nil, fmt.Errorf("Invalid file extention: %s\nMust be json, yaml, or xml.", ext)
    }

}

func exit(m any) {
    fmt.Println(m)
    os.Exit(1)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

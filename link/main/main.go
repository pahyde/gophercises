package main

import (
    "fmt"
    "log"
    "path/filepath"
    "os"
    "io"

    "link"
)

func main() {

    examples := []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}

    for _, fname := range examples {
        path := filepath.Join("../", fname)
        f, err := os.Open(path)
        check(err)
        defer func() {
            if err := f.Close(); err != nil {
                log.Fatal(err)
            }
        }()
        testLinkParser(f,fname)
    }
}

func testLinkParser(r io.Reader, name string) {
    links, err := link.Parse(r)
    check(err)
    fmt.Printf("links for: %s\n\n", name)
    for i, l := range links {
        fmt.Printf("link %d:\n", i+1)
        fmt.Printf("href: %s\n", l.Href)
        fmt.Printf("text: %s\n", l.Text)
        fmt.Println()
    }
    fmt.Println()
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

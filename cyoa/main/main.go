package main

import (
    "fmt"
    "log"
    "os"
    "net/http"
    "flag"
    
    "cyoa"
)

func main() {
    filename := flag.String(
        "story", 
        "gopher.json", 
        "path to JSON file containing the create your own adventure story",
    )

    f, err := os.Open(*filename)
    if err != nil {
        exit(fmt.Sprintf("Unable to read input file %s", *filename))
    }
    defer func() {
        if err := f.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    s, err := cyoa.JsonStory(f)
    if err != nil {
        exit(fmt.Sprintf("Unable to parse JSON in file %s", *filename))
    }
    storyMux, err := cyoa.NewStoryHandler(s,
        cyoa.WithArcToPathFn(func(arc string) string {
            if arc == "intro" {
                return "/"
            }
            return "/story/" + arc
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Starting the server on :8080")
    log.Fatal(http.ListenAndServe(":8080", storyMux))
}

func exit(m any) {
    fmt.Println(m)
    os.Exit(1)
}

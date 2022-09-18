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
    b, err := os.ReadFile(*filename)
    if err != nil {
        exit(fmt.Sprintf("Unable to read input file %s", *filename))
    }
    s, err := cyoa.NewStory(b)
    if err != nil {
        exit(fmt.Sprintf("Unable to parse JSON in file %s", *filename))
    }
    storyMux := cyoa.NewStoryMux(s)
    fmt.Println("Starting the server on :8080")
    log.Fatal(http.ListenAndServe(":8080", storyMux))
}

func exit(m any) {
    fmt.Println(m)
    os.Exit(1)
}

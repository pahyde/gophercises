package main

import (
    "fmt"
    "log"
    "os"
    "encoding/json"
    "net/http"
    "flag"
    "html/template"
)

type Arc struct {
    Title    string
    Story    []string
    Options  []Option
}

type Option struct {
    Text  string
    Arc   string
}

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
    s, err := NewStory(b)
    if err != nil {
        exit(fmt.Sprintf("Unable to parse JSON in file %s", *filename))
    }
    storyMux := NewStoryMux(s)
    fmt.Println("Starting the server on :8080")
    log.Fatal(http.ListenAndServe(":8080", storyMux))
}

func NewStory(b []byte) (map[string]Arc, error) {
    var story map[string]Arc
    if err := json.Unmarshal(b, &story); err != nil {
        log.Fatal(err)
        return nil, err
    }
    return story, nil
}

func NewStoryMux(story map[string]Arc) http.Handler {
    mux := http.NewServeMux()
    // root arc handler
    // responds with error message if exact path != "/"
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            arcHandler("intro", story).ServeHTTP(w, r)
            return
        }
        fmt.Fprintln(w, "Story arc not found!")
    })
    // descendent arc handlers
    mux.Handle("/home", arcHandler("home", story))
    mux.Handle("/new-york", arcHandler("new-york", story))
    mux.Handle("/denver", arcHandler("denver", story))
    mux.Handle("/debate", arcHandler("debate", story))
    mux.Handle("/sean-kelly", arcHandler("sean-kelly", story))
    mux.Handle("/mark-bates", arcHandler("mark-bates", story))
    return mux
}


func arcHandler(arcstr string, story map[string]Arc) http.HandlerFunc {
    arc := story[arcstr]
    return func(w http.ResponseWriter, r *http.Request) {
        // template logic
        tmpl, err := template.ParseFiles("tmpl.html")
        if err != nil {
            exit(err)
        }
        tmpl.Execute(w, arc)
        return
    }
}

func exit(m any) {
    fmt.Println(m)
    os.Exit(1)
}

package cyoa

import (
    "fmt"
    "encoding/json"
    "net/http"
    "html/template"
)

type Arc struct {
    Title    string
    Story    []string
    Options  []Option
    IsEnd    bool
}

type Option struct {
    Text  string
    Arc   string
}

func NewStory(b []byte) (map[string]Arc, error) {
    var story map[string]Arc
    if err := json.Unmarshal(b, &story); err != nil {
        return nil, err
    }
    for arcstr, arc := range story {
        if len(arc.Options) == 0 {
            arc.IsEnd = true
            story[arcstr] = arc
            break
        }
    }
    return story, nil
}

func NewStoryMux(story map[string]Arc) http.Handler {
    mux := http.NewServeMux()
    // serve static assets from static dir (e.g. /static/style.css)
    //staticServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
    //mux.Handle("/static/", staticServer) 
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
        t := template.Must(template.New("cyoa").Parse(tmpl))
        t.Execute(w, arc)
        return
    }
}

var tmpl string = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			html {
				box-sizing: border-box;
			}
			
			* {
				box-sizing: inherit;
			}
			
			body {
				background-color: rgb(23, 154, 187);
				padding: 10px 0px;
				font-family: sans-serif;
			}
			
			.arc {
				margin: 0px auto;
				background-color: #fff;
				border-radius: 10px;
				width: 500px;
				padding: 50px;
			}
			
			@media only screen and (max-width: 500px) {
			.arc {
				top: 0;
				left: 0;
				border-radius: 0px;
				width: 90%;
				transform: none;
			}
			}
			
			.title {
				display: block;
				font-size: 22px;
				font-weight: bold;
				text-align: center;
				margin-bottom: 40px;
			}
			
			.paragraphs {
				font-size: 15px;
			}
			
			.paragraph {
				text-indent: 20px;
			}
			
			.the-end {
				display: block;
				text-align: center;
				font-weight: bold;
				font-size: 18px;
				margin: auto;
			}
			
			.options {
				border-top: solid 1px;
				padding-top: 10px;
				margin-top: 10px;
			}
			
			.option {
				margin-bottom: 10px;
			}
			
			.option a {
				font-size: 15px;
				text-decoration: none;
				color: rgb(23, 154, 187);
			}		
		</style>
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <div class="arc">
            <span class="title">{{.Title}}</span>
            <div class="paragraphs">
                {{range .Story}} 
                    <p class="paragraph">
                        {{.}}
                    </p>
                {{end}}
            </div>
            {{if .IsEnd}}
                <span class="the-end">The End</span>
            {{end}}
            <div class="options">
                {{range .Options}}
                    <div class="option">
                        <a href="/{{.Arc}}">
                            > {{.Text}}
                        </a>
                    </div>
                {{end}}
            </div>
        </div>
    </body>
</html>
`

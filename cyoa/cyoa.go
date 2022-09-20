package cyoa

import (
    "errors"
    "os"
    "encoding/json"
    "net/http"
    "html/template"
)


var tmplStr string = `
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
            {{if .Options}}
                <div class="options">
                    {{range .Options}}
                        <div class="option">
                            <a href="{{atop .Arc}}">
                                > {{.Text}}
                            </a>
                        </div>
                    {{end}}
                </div>
            {{else}}
                <span class="the-end">The End</span>
            {{end}}
        </div>
    </body>
</html>
`

type Story map[string]Arc

type Arc struct {
    Title    string
    Story    []string
    Options  []Option
}

type Option struct {
    Text  string
    Arc   string
}

func JsonStory(f *os.File) (Story, error) {
    var story Story
    d := json.NewDecoder(f)
    if err := d.Decode(&story); err != nil {
        return nil, err
    }
    return story, nil
}

// provided by default or given by user as HandlerOption 
type ArcToPathFn  func(arc  string) string 
// generated automatically as inverse of ArcToPathFn
type PathToArcFn  func(path string) string 

type handler struct {
    story    Story
    t        *template.Template
    atop     ArcToPathFn            // arc  -> path
    ptoa     PathToArcFn            // path -> arc
}

type HandlerOption func(h *handler) error

func WithTemplate(t *template.Template) HandlerOption {
    return func(h *handler) error {
        h.t = t
        return nil
    }
}

func WithArcToPathFn(f ArcToPathFn) HandlerOption {
    return func(h *handler) error {
        fInv, err := invert(f, h.story)
        if err != nil {
            return err
        }
        h.atop = f
        h.ptoa = fInv
        return nil
    }
}

// construct PathToArcFn from ArcToPathFn. Error if provided fn is not one-to-one
func invert(f ArcToPathFn, s Story) (PathToArcFn, error) {
    ptoaMap := make(map[string]string)
    for arc, _ := range s {
        path := f(arc)
        ptoaMap[path] = arc
    }
    if len(ptoaMap) < len(s) {
        return nil, errors.New("Provided ArcToPath fn is not one-to-one. One or more arcs map to the same path.")
    }
    fInv := func(path string) string {
        return ptoaMap[path]
    }
    return fInv, nil
}

// NewStoryHandler: 
// accepts handler options for arcToPath function and template
// arcToPath must be injective (one-to-one)
func NewStoryHandler(s Story, opts ...HandlerOption) (http.Handler, error) {
    h := defaultStoryHandler(s)
    for _, o := range opts {
        // update h with option
        if err := o(&h); err != nil {
            return nil, err
        }
    }
    if h.t == nil {
        // optional template not provided
        // generate new template from ArcToPathFn h.atop
        // h.atop maps arc options back to vaild href paths in template
        h.t = DefaultTemplate(h.atop)
    }
    return h, nil
}

func defaultStoryHandler(s Story) handler {
    var atop ArcToPathFn = func(arc string) string {
        if arc == "intro" {
            return "/"
        }
        return "/" + arc
    }
    ptoa := must(invert(atop, s))
    return handler{s, nil, atop, ptoa}
}

func DefaultTemplate(atop ArcToPathFn) *template.Template {
    fmap := template.FuncMap{"atop": atop}
    return template.Must(template.New("cyoa").Funcs(fmap).Parse(tmplStr)) 
}


func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    if arc, ok:= h.story[h.ptoa(path)]; ok {
        if err := h.t.Execute(w, arc); err != nil {
            http.Error(w, "Something went wrong.", http.StatusInternalServerError)
        }
        return
    }
    http.Error(w, "Story arc not found.", http.StatusNotFound)
}

func must(fInv PathToArcFn, err error) PathToArcFn {
    if err != nil {
        panic(err)
    }
    return fInv
}


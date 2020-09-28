package web

import (
    "html/template"
    "net/http"
    "os"
    "regexp"
)

type Connection interface{
    SavePage(*Page) error
    GetPage(string) (*Page, error)
}

type Page struct {
    Title string `json:”title,omitempty”`
    Body  []byte `json:”body,omitempty”`
}

type TemplateData struct {
    Title string
    Body  []byte
    Host  string
}

var templates *template.Template
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var hostname, _ = os.Hostname()

func InitTemplates(filenames ...string) error {
    files, err := template.ParseFiles(filenames...)
    if err != nil {
        return err
    }
    templates = template.Must(files, nil)
    return nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", TemplateData{p.Title, p.Body, hostname})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ViewHandler(w http.ResponseWriter, r *http.Request, con Connection, title string) {
    p, err := con.GetPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, con Connection, title string) {
    p, err := con.GetPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, con Connection, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := con.SavePage(p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func MakeHandler(fn func(w http.ResponseWriter, r *http.Request, con Connection, title string), con Connection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, con, m[2])
    }
}

func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/frontPage", http.StatusFound)
}


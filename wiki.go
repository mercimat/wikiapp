package main

import (
    "context"
    "html/template"
    "log"
    "net/http"
    "os"
    "regexp"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
    Title string `json:”title,omitempty”`
    Body  []byte `json:”body,omitempty”`
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// mongodb global vars
var collection *mongo.Collection
var ctx = context.TODO()

func InitDB() {
    server := "mongodb://localhost:27017/"
    if _, ok := os.LookupEnv("MONGODB_HOST"); ok {
        server = os.ExpandEnv("mongodb://${MONGODB_HOST}:27017/")
    }
    clientOptions := options.Client().ApplyURI(server)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("wikiapp").Collection("pages")
    log.Println("connection established")
}

func (p *Page) save() error {
    var page Page
    filter := bson.D{{"title", p.Title}}
    update := bson.D{{"$set", bson.D{{"body", p.Body}}}}
    opts := options.FindOneAndUpdate().SetUpsert(true)
    err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&page)
    if err != nil && err != mongo.ErrNoDocuments {
        return err
    }
    return nil
}

func loadPage(title string) (*Page, error) {
    var p Page
    filter := bson.D{{"title", title}}
    err := collection.FindOne(context.TODO(), filter).Decode(&p)
    return &p, err
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/frontPage", http.StatusFound)
}

func main() {
    InitDB()
    http.HandleFunc("/", frontPageHandler)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    log.Fatal(http.ListenAndServe(":8090", nil))
}

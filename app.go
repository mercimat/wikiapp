package main

import (
    "log"
    "net/http"
    "os"
    "github.com/mercimat/wikiapp/web"
    "github.com/mercimat/wikiapp/db"
)

func main() {
    mongoServer := "mongodb://localhost:27017/"
    if _, ok := os.LookupEnv("MONGODB_HOST"); ok {
        mongoServer = os.ExpandEnv("mongodb://${MONGODB_HOST}:27017/")
    }

    dbcon := db.NewMongoDB(mongoServer, "wikiapp", "pages")
    log.Println("connection established")

    url := ":8090"
    if _, ok := os.LookupEnv("WIKIAPP_URL"); ok {
        url = os.Getenv("WIKIAPP_URL")
    }

    err := web.InitTemplates("tmpl/edit.html", "tmpl/view.html")
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", web.FrontPageHandler)
    http.HandleFunc("/view/", web.MakeHandler(web.ViewHandler, dbcon))
    http.HandleFunc("/edit/", web.MakeHandler(web.EditHandler, dbcon))
    http.HandleFunc("/save/", web.MakeHandler(web.SaveHandler, dbcon))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    log.Fatal(http.ListenAndServe(url, nil))
}

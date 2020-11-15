package server

import (
	"itWiki/router"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Run - start a main thread server
func Run() {
	log.Println("Starting itWiki server")
	r := mux.NewRouter()
	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("/app/assets/")))
	r.PathPrefix("/assets/").Handler(assets)
	r.HandleFunc("/", router.IndexHandler)
	r.HandleFunc("/add/article/", router.AddArticleHandler)
	r.HandleFunc("/api/articles/{count}", router.APIGetArticlesHandler)
	r.HandleFunc("/api/add/article/", router.APICreateArticle).Methods(http.MethodPost)
	r.HandleFunc("/settings/initializedbtable", router.InitializeDBTable)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

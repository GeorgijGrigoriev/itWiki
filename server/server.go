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
	r.HandleFunc("/read/article/{id}", router.GenerateArticlePage)
	r.HandleFunc("/api/articles/{count}", router.APIGetArticlesHandler)
	r.HandleFunc("/api/articles/add/", router.APICreateArticle).Methods("POST")
	r.HandleFunc("/api/category/add", router.APICreateCategory).Methods("POST")
	r.HandleFunc("/settings/initializedbtable", router.InitializeDBTable)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

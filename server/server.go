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
	r.HandleFunc("/edit/article/{id}", router.APIGetUpdateHandler).Methods("GET")
	r.HandleFunc("/settings/", router.SettingsPageHandler)
	r.HandleFunc("/api/articles/{count}", router.APIGetArticlesHandler)
	r.HandleFunc("/api/articles/add/", router.APICreateArticle).Methods("POST")
	r.HandleFunc("/api/categories/add", router.APICreateCategory).Methods("POST")
	r.HandleFunc("/api/categories/get", router.APIGetCategories)
	r.HandleFunc("/settings/initializedbtable", router.InitializeDBTable)
	r.HandleFunc("/api/articles/delete/{id}", router.APIDeleteArticle).Methods("DELETE")
	r.HandleFunc("/api/articles/update/{id}", router.APIPostUpdateHandler).Methods("POST")
	r.HandleFunc("/read/article/mdtest/{id}", router.MDTest)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

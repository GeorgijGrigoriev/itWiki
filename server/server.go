package server

import (
	"itWiki/auth"
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
	r.Use(auth.JWTAuth)
	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("/app/assets/")))
	r.PathPrefix("/assets/").Handler(assets)
	r.HandleFunc("/auth/createaccount", router.APIMakeAccount).Methods("POST")
	r.HandleFunc("/auth/login", router.APIMakeLogin).Methods("POST")
	r.HandleFunc("/auth/logout", router.APIMakeLogout)
	r.HandleFunc("/auth/refresh", router.APIMakeRefresh)
	r.HandleFunc("/", router.LoginHandler)
	r.HandleFunc("/app", router.IndexHandler)
	r.HandleFunc("/app/article/add/", router.AddArticleHandler)
	r.HandleFunc("/app/article/read/", router.ReadArticleHandler)
	r.HandleFunc("/app/article/edit/", router.EditArticleHandler)
	r.HandleFunc("/app/settings/", router.SettingsPageHandler)
	r.HandleFunc("/api/articles/{count}/", router.APIGetArticlesHandler)
	r.HandleFunc("/api/article/add/", router.APICreateArticle).Methods("POST")
	r.HandleFunc("/api/article/get/", router.GenerateArticlePage).Methods("POST")
	r.HandleFunc("/api/article/delete/", router.APIDeleteArticle).Methods("DELETE")
	r.HandleFunc("/api/article/update/", router.APIPostUpdateHandler).Methods("POST")
	r.HandleFunc("/api/categories/add/", router.APICreateCategory).Methods("POST")
	r.HandleFunc("/api/categories/get/", router.APIGetCategories)
	r.HandleFunc("/settings/initializedbtable", router.InitializeDBTable)
	r.HandleFunc("/app/article/upload/", router.UploadFilesHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

package router

import (
	"fmt"
	"html/template"
	"itWiki/db"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

//Article - default wiki page structure
type Article struct {
	ID       int64 `gorm:"primaryKey"`
	Category int64
	Title    string
	Post     string
	Comment  string
}

//Message - structure for generating supporting json responses
type Message struct {
	Message string `json:"message"`
}

//InitializeDBTable - initial script for create table if not exist
func InitializeDBTable(w http.ResponseWriter, r *http.Request) {
	initCommand := `
	CREATE TABLE IF NOT EXISTS articles (
		id INT NOT NULL AUTO_INCREMENT,
		creation_date TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		category INT NOT NULL,
		title TINYTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		post TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		comment TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		PRIMARY KEY (id)
	) ENGINE=InnoDB;
	`
	db := db.Driver()
	err := db.Exec(initCommand)
	fmt.Fprint(w, err)
	fmt.Fprint(w, "Database initialized OK")
}

//IndexHandler - index page handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
}

func AddArticleHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_article.html")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	templateDir := filepath.Join("/app/templates", tmpl)
	t := template.Must(template.ParseFiles(templateDir))
	err := t.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//APIGetArticlesHandler - GET method handler for api
func APIGetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var article Article
	db := db.Driver()
	if vars["count"] == "all" {
		result := db.Find(&article)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result.RowsAffected)
		//result.Error
	}
	if vars["count"] == "first" {
		result := db.First(&article)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result)
	}
}

//APICreateArticle - CREATE method handler for add new article
func APICreateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprint(w, "OK")
	} else {
		fmt.Fprint(w, http.StatusMethodNotAllowed)
	}
}

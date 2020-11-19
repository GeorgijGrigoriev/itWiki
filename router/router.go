package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"itWiki/db"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

//Article - default wiki page structure
type Article struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	CreationDate time.Time `json:"creation_date"`
	Category     int64     `json:"category"`
	Title        string    `json:"title"`
	Post         string    `json:"post"`
	Comment      string    `json:"comment"`
}

//ArticleNew - new struct
type ArticleNew struct {
	ArticleID    int64     `json:"article_id" gorm:"primaryKey"`
	CreationDate time.Time `json:"creation_date"`
	Category     int64     `json:"category"`
	Title        string    `json:"title"`
	Post         string    `json:"post"`
	Comment      string    `json:"comment"`
	CategoryName string    `json:"category_name"`
}

//MessageResults - structure for multiple rows
type MessageResults struct {
	Message string       `json:"message"`
	Data    []ArticleNew `json:"data"`
	Total   int64        `json:"total"`
}

//MessageResult - structure for one row
type MessageResult struct {
	Message string  `json:"message"`
	Data    Article `json:"data"`
	Total   int64   `json:"total"`
}

//Page - default wiki page
type Page struct {
	Title string
}

//Categories - struct for categories
type Categories struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
}

//InitializeDBTable - initial script for create table if not exist
func InitializeDBTable(w http.ResponseWriter, r *http.Request) {
	mainTable := `
	CREATE TABLE IF NOT EXISTS articles (
		article_id INT NOT NULL AUTO_INCREMENT,
		creation_date TIMESTAMP default CURRENT_TIMESTAMP,
		category INT NOT NULL,
		title TINYTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		post TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		comment TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		PRIMARY KEY (id)
	) ENGINE=InnoDB;
	`
	categoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INT NOT NULL AUTO_INCREMENT,
		category_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB;
	`
	db := db.Driver()
	err := db.Exec(mainTable)
	err = db.Exec(categoriesTable)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, "Database initialized OK")
	}
}

//IndexHandler - index page handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
}

//AddArticleHandler - handler for new article page
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
	db := db.Driver()
	w.Header().Set("Content-Type", "application/json")
	if vars["count"] == "all" {
		var articles []ArticleNew
		raw := db.Table("articles").Select("articles.article_id, articles.creation_date, articles.category, articles.title, articles.post, articles.comment, categories.id, categories.category_name").Joins("left join categories on categories.id = articles.category").Scan(&articles)
		totalRows := raw.RowsAffected
		json.NewEncoder(w).Encode(MessageResults{Message: "OK", Data: articles, Total: totalRows})
	}
}

//APICreateArticle - CREATE method handler for add new article
func APICreateArticle(w http.ResponseWriter, r *http.Request) {
	var data Article
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewDecoder(r.Body).Decode(&data)
		db := db.Driver()
		res := db.Select("Category", "Title", "Post").Create(&data)
		fmt.Fprint(w, res.RowsAffected)
	} else {
		fmt.Fprint(w, http.StatusMethodNotAllowed)
	}
}

//GenerateArticlePage - generate article full page
func GenerateArticlePage(w http.ResponseWriter, r *http.Request) {
	var page ArticleNew
	vars := mux.Vars(r)
	db := db.Driver()
	log.Println(vars)
	_ = db.Table("articles").Select("articles.article_id, articles.creation_date, articles.category, articles.title, articles.post, articles.comment, categories.id, categories.category_name").Joins("left join categories on categories.id = articles.category").Where("articles.article_id=?", vars["id"]).Scan(&page)
	t := template.Must(template.ParseFiles("/app/templates/article.html"))
	err := t.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

//APICreateCategory - CREATE method for add new category
func APICreateCategory(w http.ResponseWriter, r *http.Request) {
	var data Categories
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewDecoder(r.Body).Decode(&data)
		db := db.Driver()
		res := db.Table("categories").Select("Category").Create(&data)
		fmt.Fprint(w, res.RowsAffected)
	} else {
		fmt.Fprint(w, http.StatusMethodNotAllowed)
	}
}

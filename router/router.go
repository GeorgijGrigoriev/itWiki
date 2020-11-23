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
	"github.com/russross/blackfriday"
)

//Article - default wiki page structure
type Article struct {
	ArticleID int64  `json:"article_id"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Post      string `json:"post"`
	Comment   string `json:"comment"`
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

//MessageResultsCategories - struct for categories
type MessageResultsCategories struct {
	Message string       `json:"message"`
	Data    []Categories `json:"data"`
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
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
}

//Category - struct for add category
type Category struct {
	CategoryName string `json:"category_name"`
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

//IndexHandler - handler for index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
}

//AddArticleHandler - handler for new article page
func AddArticleHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_article.html")
}

//SettingsPageHandler - handler for settings page
func SettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "settings.html")
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
		log.Println(data)
		res := db.Table("articles").Select("Category", "Title", "Post").Create(&data)
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
	_ = db.Table("articles").Select("articles.article_id, articles.creation_date, articles.category, articles.title, articles.post, articles.comment, categories.id, categories.category_name").Joins("left join categories on categories.id = articles.category").Where("articles.article_id=?", vars["id"]).Scan(&page)
	t := template.Must(template.ParseFiles("/app/templates/article.html"))
	err := t.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

//APICreateCategory - CREATE method for add new category
func APICreateCategory(w http.ResponseWriter, r *http.Request) {
	var data Category
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewDecoder(r.Body).Decode(&data)
		log.Println(data)
		db := db.Driver()
		res := db.Select("CategoryName").Create(&data)
		fmt.Fprint(w, res.RowsAffected)
	} else {
		fmt.Fprint(w, http.StatusMethodNotAllowed)
	}
}

//APIGetCategories - get all categories
func APIGetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Categories
	db := db.Driver()
	res := db.Table("categories").Find(&categories)
	totalRows := res.RowsAffected
	json.NewEncoder(w).Encode(MessageResultsCategories{Message: "OK", Data: categories, Total: totalRows})
}

//APIDeleteArticle - delete article by id
func APIDeleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		vars := mux.Vars(r)
		var article ArticleNew
		db := db.Driver()
		res := db.Table("articles").Where("article_id like ?", vars["id"]).Delete(&article)
		totalRows := res.RowsAffected
		fmt.Fprint(w, totalRows)
	}
}

//APIGetUpdateHandler - update article
func APIGetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		var article Article
		db := db.Driver()
		db.Table("articles").Where("article_id like ?", vars["id"]).Scan(&article)
		t := template.Must(template.ParseFiles("/app/templates/edit_article.html"))
		err := t.Execute(w, article)
		if err != nil {
			log.Println(err)
		}
	}
}

//APIPostUpdateHandler - update article
func APIPostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		vars := mux.Vars(r)
		var article Article
		_ = json.NewDecoder(r.Body).Decode(&article)
		db := db.Driver()
		db.Model(&article).Where("article_id like ?", vars["id"]).Updates(&article)
		fmt.Fprint(w, vars["id"])
	}
}

func MDTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var article Article
	type page struct {
		Content template.HTML
	}
	db := db.Driver()
	res := db.Table("articles").Where("article_id like ?", vars["id"]).Scan(&article)
	totalRows := res.RowsAffected
	json.NewEncoder(w).Encode(MessageResult{Message: "OK", Data: article, Total: totalRows})
	tmp := article.Post
	md := []byte(tmp)
	html := blackfriday.MarkdownCommon(md)
	t := template.Must(template.ParseFiles("/app/templates/article.html"))
	t.Execute(w, page{Content: template.HTML(html)})
}

package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	a "itWiki/auth"
	"itWiki/db"
	u "itWiki/utils"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//Account - struct for user account
type Account struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

//Article - default wiki page structure
type Article struct {
	ArticleID string `json:"article_id"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Post      string `json:"post"`
	Comment   string `json:"comment"`
}

//ArticleNew - new struct
type ArticleNew struct {
	ArticleID    int64         `json:"article_id" gorm:"primaryKey"`
	CreationDate time.Time     `json:"creation_date"`
	Category     int64         `json:"category"`
	Title        string        `json:"title"`
	Post         string        `json:"post"`
	Comment      string        `json:"comment"`
	CategoryName string        `json:"category_name"`
	Content      template.HTML `json:"content"`
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
	Message string     `json:"message"`
	Data    ArticleNew `json:"data"`
	Total   int64      `json:"total"`
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
		PRIMARY KEY (article_id)
	) ENGINE=InnoDB;
	`
	categoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INT NOT NULL AUTO_INCREMENT,
		category_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB;
	`
	accoutsTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		password VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		created_at TIMESTAMP default NULL,
		updated_at TIMESTAMP default NULL,
		deleted_at TIMESTAMP default NULL,
		token TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
		PRIMARY KEY (id)
	) ENGINE=InnoDB;
	`
	db := db.Driver()
	err := db.Exec(mainTable)
	err = db.Exec(categoriesTable)
	err = db.Exec(accoutsTable)
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

//ReadArticleHandler - handler for article page
func ReadArticleHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "article.html")
}

//EditArticleHandler - handler for edit articles
func EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "edit_article.html")
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
	type ArticleID struct {
		ID string `json:"id"`
	}
	var id ArticleID
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&id)
	db := db.Driver()
	res := db.Table("articles").Select("articles.article_id, articles.creation_date, articles.category, articles.title, articles.post, articles.comment, categories.id, categories.category_name").Joins("left join categories on categories.id = articles.category").Where("articles.article_id=?", id.ID).Scan(&page)
	totalRows := res.RowsAffected
	tmp := page.Post
	md := []byte(tmp)
	html := blackfriday.MarkdownCommon(md)
	page.Content = template.HTML(html)
	json.NewEncoder(w).Encode(MessageResult{Message: "OK", Data: page, Total: totalRows})
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
		type ArticleID struct {
			ID string `json:"id"`
		}
		var id ArticleID
		var article ArticleNew
		w.Header().Set("Content-Type", "application/json")
		json.NewDecoder(r.Body).Decode(&id)
		log.Println(id.ID)
		db := db.Driver()
		res := db.Table("articles").Where("article_id = ?", id.ID).Delete(&article)
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
		db.Table("articles").Where("article_id = ?", vars["id"]).Scan(&article)
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
		var article Article
		log.Println(article.ArticleID)
		_ = json.NewDecoder(r.Body).Decode(&article)
		db := db.Driver()
		db.Model(&article).Where("article_id = ?", article.ArticleID).Updates(&article)
		fmt.Fprint(w, article.ArticleID)
	}
}

//ValidateAccount - validate incoming data
func (account *Account) ValidateAccount() (map[string]interface{}, bool) {
	if len(account.Password) < 6 {
		return u.Message(false, "Password must be more then 6 symbols."), false
	}

	temp := &Account{}
	db := db.Driver()
	err := db.Table("accounts").Where("username = ?", account.Username).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Contact support."), false
	}
	if temp.Username != "" {
		return u.Message(false, "User already exists."), false
	}

	return u.Message(false, "Requirments passed"), true
}

//CreateAccount - create account handler
func (account *Account) CreateAccount() map[string]interface{} {
	if resp, ok := account.ValidateAccount(); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	db := db.Driver()
	db.Table("accounts").Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account")
	}

	tk := &a.Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("api-secret"))
	account.Token = tokenString
	account.Password = ""
	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

//LoginAccount - login account handler
func LoginAccount(username, password string) map[string]interface{} {
	account := &Account{}
	db := db.Driver()
	err := db.Table("accounts").Where("username = ?", username).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "User not found.")
		}
		return u.Message(false, "Connection error.")
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials.")
	}
	account.Password = ""
	tk := &a.Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("api-secret"))
	account.Token = tokenString

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetUser - get user handler
func GetUser(u uint) *Account {
	acc := &Account{}
	db := db.Driver()
	db.Table("accounts").Where("id = ?", u).First(acc)
	if acc.Username == "" {
		return nil
	}
	acc.Password = ""
	return acc
}

//APIMakeAccount - make account handler
func APIMakeAccount(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := account.CreateAccount()
	u.Respond(w, resp)
}

//APIMakeLogin - make login handler
func APIMakeLogin(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := LoginAccount(account.Username, account.Password)
	u.Respond(w, resp)
}

//APIMakeLogout - make logout handler
func APIMakeLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not implemented")
}

//APIMakeRefresh - make refresh token handler
func APIMakeRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not implemented")
}

//LoginHandler - login page handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth.html")
}

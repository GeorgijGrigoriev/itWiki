package db

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Driver - create instance of connections driver
func Driver() *gorm.DB {
	dsn := "root:secret@tcp(mysql:3306)/itwiki?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := sql.Open("mysql", dsn)
	mysqlDB.SetMaxIdleConns(10)
	mysqlDB.SetMaxOpenConns(100)
	mysqlDB.SetConnMaxLifetime(5 * time.Minute)
	if err != nil {
		log.Println(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: mysqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}
	return db
}

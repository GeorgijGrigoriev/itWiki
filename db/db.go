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
	if err != nil {
		log.Println(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: mysqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}

	mysqlDBConfig, err := db.DB()
	mysqlDBConfig.SetMaxIdleConns(10)
	mysqlDBConfig.SetConnMaxLifetime(24 * time.Hour)
	mysqlDBConfig.SetMaxOpenConns(500)

	return db
}

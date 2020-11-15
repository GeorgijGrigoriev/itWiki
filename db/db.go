package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Driver - create instance of connections driver
func Driver() *gorm.DB {
	dsn := "root:secret@tcp(mysql:3306)/itwiki?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return db
}

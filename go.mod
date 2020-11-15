module itWiki/modules

go 1.15

require (
	itWiki/db v0.0.0-00010101000000-000000000000 // indirect
	itWiki/router v0.0.0-00010101000000-000000000000 // indirect
	itWiki/server v0.0.0
	itWiki/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	gorm.io/driver/mysql v1.0.3 // indirect
	gorm.io/gorm v1.20.6 // indirect
)

replace (
	itWiki/db => ./db
	itWiki/router => ./router
	itWiki/server => ./server
	itWiki/utils => ./utils
)

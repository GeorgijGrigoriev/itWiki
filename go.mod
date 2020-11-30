module itWiki/modules

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/russross/blackfriday v1.5.2 // indirect
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	gorm.io/gorm v1.20.7 // indirect
	itWiki/auth v0.0.0-00010101000000-000000000000 // indirect
	itWiki/db v0.0.0-00010101000000-000000000000 // indirect
	itWiki/router v0.0.0-00010101000000-000000000000 // indirect
	itWiki/server v0.0.0
	itWiki/utils v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	itWiki/auth => ./auth
	itWiki/db => ./db
	itWiki/router => ./router
	itWiki/server => ./server
	itWiki/utils => ./utils
)

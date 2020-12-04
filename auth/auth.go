package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	u "itWiki/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

//Token - rights struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

func HasSuffix(path string, parts []string) bool {
	for _, part := range parts {
		if strings.HasSuffix(path, part) == true {
			return true
		}
	}
	return false
}

//JWTAuth - midleware for JWT auth
var JWTAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/auth/login", "/app/article/upload/", "/app/settings/", "/auth/createaccount", "/", "/app", "/app/article/add/", "/app/article/read/", "/app/article/edit/"}
		requestPath := r.URL.Path
		if HasSuffix(requestPath, []string{"png", "js", "css"}) == true {
			next.ServeHTTP(w, r)
			return
		}
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1]
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("api-secret"), nil
		})

		if err != nil {
			response = u.Message(false, "Corrupted auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Invalid auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		fmt.Sprintf("User %", tk.UserID)
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

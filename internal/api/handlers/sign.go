package handlers

import (
	"encoding/json"
	"finalProject/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"os"
)

// Метод для аутентификации в сервис
func Signin(r *http.Request) []byte {
	newResponse := &db.Response{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = json.Unmarshal(body, &newResponse)
	if err != nil {
		log.Println(err)
		return nil
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(newResponse.Password))
	if err != nil {
		log.Println(err)
		return nil
	}

	return []byte("{\"token\":\"" + tokenString + "\"}")
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(os.Getenv("TODO_PASSWORD")) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Println(err)
				http.Error(w, "Error Parse", http.StatusUnauthorized)
				return
			}

			parse, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("TODO_PASSWORD")), nil
			})
			if err != nil {
				log.Println("Error Parse", err)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if !parse.Valid {
				log.Println("Token is invalid")
				http.Error(w, "Token is invalid", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}

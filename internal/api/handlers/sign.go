package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"finalProject/internal/db"

	"github.com/golang-jwt/jwt/v5"
)

// Метод для аутентификации в сервис с помощью пароля
func Signin(r *http.Request) []byte {

	response := &db.Response{}
	request := &db.Request{}

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Преобразование в структуру
	err = json.Unmarshal(body, &request)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Проверка переданного пароля
	if request.Password != os.Getenv("TODO_PASSWORD") {
		return response.LogResponseError("Неверный пароль")
	}

	// Формирование токена
	token := jwt.New(jwt.SigningMethodHS256)
	response.Token, err = token.SignedString([]byte(request.Password))
	if err != nil {
		return response.LogResponseError(err.Error())
	}
	
	// Сериализация JSON
	return response.Marshal()
}

// Метод для авторизации с помощью токена
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Проверка наличия пароля
		if len(os.Getenv("TODO_PASSWORD")) > 0 {

			// Проверка куки
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Println("No cookie found", err)
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Парсинг токена
			parse, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("TODO_PASSWORD")), nil
			})
			if err != nil {
				log.Println("Error parse token", err)
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Проверка на валидность токена
			if !parse.Valid {
				log.Println("Token is invalid")
				http.Error(w, "Token is invalid", http.StatusUnauthorized)
				return
			}
		}

		// Запуск следующей обработки
		next(w, r)
	})
}

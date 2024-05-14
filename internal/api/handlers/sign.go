package handlers

import (
	"crypto/sha256"
	"encoding/base64"
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
	h := hashPass(request.Password)
	claims := &jwt.MapClaims{"passSHA256": h}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("TODO_PASSWORD")), nil
			})
			if err != nil {
				log.Println("Error parse token", err)
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Парсинг хеша
			h := hashPass(os.Getenv("TODO_PASSWORD"))
			sha := base64.StdEncoding.EncodeToString(h)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Error parse hash", err)
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Проверка хеша пароля в токене с хешем пароля из переменной окружения
			if sha != claims["passSHA256"] {
				log.Println("Token pass and system pass hashes not match")
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Проверка на валидность токена
			if !token.Valid {
				log.Println("Token is invalid")
				http.Error(w, "Token is invalid", http.StatusUnauthorized)
				return
			}
		}

		// Запуск следующей обработки
		next(w, r)
	})
}

// Метод для шифрования пароля
func hashPass(pass string) []byte {
	h := sha256.New()
	h.Write([]byte(pass))

	return h.Sum(nil)
}

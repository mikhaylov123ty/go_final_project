package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"finalProject/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// Метод для аутентификации в сервисе с помощью пароля
func Signin(r *http.Request) []byte {

	response := &models.Response{}
	request := &models.Request{}

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

	// Создание хеша пароля и добавление в claims часть токена
	h, err := models.HashPass(request.Password)
	if err != nil {
		return response.LogResponseError(err.Error())
	}
	claims := &jwt.MapClaims{"passSHA256": h}

	// Формирование и подпись токена
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
				models.LogAuthError("No cookie found:", err, w)
				return
			}

			// Верификация токена с ключом подписи - паролем
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("TODO_PASSWORD")), nil
			})
			if err != nil {
				models.LogAuthError("Error parse token:", err, w)
				return
			}

			// Парсинг хеша
			h, err := models.HashPass(os.Getenv("TODO_PASSWORD"))
			if err != nil {
				models.LogAuthError("Error hash password:", err, w)
				return
			}
			sha := base64.StdEncoding.EncodeToString(h)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				models.LogAuthError("Error decode claims", err, w)
				return
			}

			// Дополнительная проверка хеша пароля в claims токена с хешем пароля из переменной окружения
			if sha != claims["passSHA256"] {
				models.LogAuthError("Token pass and system pass hashes not match:", err, w)
				return
			}

			// Основная проверка на валидность токена
			if !token.Valid {
				models.LogAuthError("Token is invalid:", err, w)
				return
			}
		}

		// Запуск следующей обработки
		next(w, r)
	})
}

package handlers

import "net/http"

// Метод для аутентификации в сервис
func Auth(r *http.Request) []byte {

	return []byte("{\"token\":\"\"}")
}

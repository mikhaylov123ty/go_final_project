package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"

	"finalProject/internal/db"
)

// Метод для аутентификации в сервис
func Auth(r *http.Request) []byte {
	newResponse := &db.Response{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = json.Unmarshal(body, &newResponse)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(newResponse.Password)

	res := make([]byte, base64.StdEncoding.EncodedLen(len(newResponse.Password)))
	base64.StdEncoding.Encode(res, []byte(newResponse.Password))
	fmt.Println(string(res))

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(newResponse.Password))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(tokenString)

	return []byte("{\"token\":\"\"}" + tokenString)
}

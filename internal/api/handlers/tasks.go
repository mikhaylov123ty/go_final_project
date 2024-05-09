package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"finalProject/internal/db"
)

// Метод для запроса всех задач
func GetAllTasks() []byte {
	var err error
	newResponse := &db.Response{}

	// Выполнение запроса к базе
	newResponse.Tasks, err = db.DbInstance.GetAllTasks()
	if err != nil {
		log.Println("{\"error\":\"ошибка запроса в базу\"}", err.Error())
		return []byte("{\"error\":\"ошибка запроса в базу\"}")
	}

	// Сериализация JSON
	res, err := json.Marshal(newResponse)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

// Метод для поиска задачи
func GetTasksBySearch(r *http.Request) []byte {
	var err error
	newResponse := &db.Response{}

	// Выполнение запроса к базе
	newResponse.Tasks, err = db.DbInstance.GetTaskBySearch(r.URL.Query().Get("search"))
	if err != nil {
		log.Println("{\"error\":\"ошибка запроса в базу\"}", err.Error())
		return []byte("{\"error\":\"ошибка запроса в базу\"}")
	}

	// Сериализация JSON
	res, err := json.Marshal(newResponse)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

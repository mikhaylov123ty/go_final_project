package handlers

import (
	"fmt"
	"net/http"

	"finalProject/internal/db"
)

// Метод для запроса всех задач
func GetAllTasks() []byte {
	var err error
	response := &db.Response{}

	// Выполнение запроса к базе
	response.Tasks, err = db.DbInstance.GetAllTasks()
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Сериализация JSON
	return response.Marshal()
}

// Метод для поиска задачи
func GetTasksBySearch(r *http.Request) []byte {
	var err error
	response := &db.Response{}

	// Выполнение запроса к базе
	response.Tasks, err = db.DbInstance.GetTaskBySearch(r.URL.Query().Get("search"))
	fmt.Println(len(response.Tasks))
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Сериализация JSON
	return response.Marshal()
}

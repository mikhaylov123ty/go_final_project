package handlers

import (
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

	// Проверка на пустой слайс. Сериализация его не обрабатывает, т.к. установлен omitempty для этого поля
	if len(response.Tasks) == 0 {
		return []byte("{\"tasks\":[]}")
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
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Проверка на пустой слайс. Сериализация его не обрабатывает, т.к. установлен omitempty для этого поля
	if len(response.Tasks) == 0 {
		return []byte("{\"tasks\":[]}")
	}

	// Сериализация JSON
	return response.Marshal()
}

package api

import (
	"log"
	"net/http"
	"time"

	"finalProject/internal/api/handlers"
	"finalProject/internal/models"
	"finalProject/internal/tasks"
)

const (
	incorrectRequest = "Не корректный запрос"
)

// Основной обработчик для ручки api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Распределение допустимых запросов
	switch {

	// Запрос POST, создание задачи
	case r.Method == "POST":
		w.Write(handlers.AddTask(r))

	// Запрос PUT, изменение задачи
	case r.Method == "PUT":
		w.Write(handlers.ChangeTask(r))

	// Запрос DELETE, удаление задачи
	case r.Method == "DELETE":
		w.Write(handlers.DeleteTaskById(r))

	// Запрос GET по id задачи
	case r.Method == "GET":
		w.Write(handlers.GetTaskById(r))

	// По умолчанию возвращает статус с ошибкой
	default:
		resp := models.Response{Error: incorrectRequest}
		w.Write(resp.Marshal())
	}
}

// Основной обработчик для ручки api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Распределение допустимых запросов
	switch r.URL.Query().Has("search") {

	// Запрос с текстом в поле search
	case true:
		w.Write(handlers.GetTasksBySearch(r))

	// По умолчанию возвращает все задачи
	default:
		w.Write(handlers.GetAllTasks())
	}
}

// Основной обработчик для ручки api/task
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Распределение допустимых запросов
	switch {

	// Запрос POST, завершение задачи
	case r.Method == "POST":
		w.Write(handlers.DoneTask(r))

		// По умолчанию возвращает статус с ошибкой
	default:
		resp := models.Response{Error: incorrectRequest}
		w.Write(resp.Marshal())
	}
}

// Основной обработчик для ручки api/sign
func SignHandler(w http.ResponseWriter, r *http.Request) {

	// Установка типа контента json для ответа на запросы
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Распределение допустимых запросов
	switch {

	// Запрос POST, проверка пароля и формирование токена
	case r.Method == "POST":
		w.Write(handlers.Signin(r))

	default:
		resp := models.Response{Error: incorrectRequest}
		w.Write(resp.Marshal())
	}
}

// Обработчик для ручки api/nextDate
func NextDate(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	now, err := time.Parse("20060102", values.Get("now"))
	if err != nil {
		log.Println(err)
		return
	}

	res, err := tasks.NextDateHandler(now, values.Get("date"), values.Get("repeat"))
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte(res))
}

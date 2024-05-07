package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"finalProject/internal/db"
	"finalProject/internal/tasks"
)

type response struct {
	Tasks []*db.Task `json:"tasks"`
}

// Основной обработчик для ручки api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	// Распределение допустимых запросов
	switch {

	// Запрос POST, создание задачи
	case r.Method == "POST":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(addTask(r))

	// Запрос PUT, изменение задачи
	case r.Method == "PUT":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(changeTask(r))

	// Запрос по id задачи
	case r.URL.Query().Has("id") == true:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getTaskById(r))

	// По умолчанию возвращает статус с ошибкой
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"error\":\"Не указан идентификатор\"}"))

	}
}

// Основной обработчик для ручки api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Распределение допустимых запросов
	switch r.URL.Query().Has("search") {

	// Запрос с текстом в поле search
	case true:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getTasksBySearch(r))

	// По умолчанию возвращает все задачи
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getAllTasks())
	}
}

// Метод для запроса всех задач
func getAllTasks() []byte {
	var err error
	newResponse := &response{}

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
func getTasksBySearch(r *http.Request) []byte {
	var err error
	newResponse := &response{}

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

// Метод для добавления задачи в базу
func addTask(r *http.Request) []byte {
	newTask := &db.Task{}

	//TODO don't like how here callback bytes res and error

	// Проверка вводных данных задачи
	res, err := newTask.CheckTask(r)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	// Добавление задачи в базу
	id, err := db.DbInstance.AddTask(newTask)
	if err != nil {
		log.Println("{\"error\":\"Не удалось добавить в базу\"}", err.Error())
		return []byte("{\"error\":\"Не удалось добавить в базу\"}")
	}

	strID := strconv.Itoa(id)

	return []byte("{\"id\":\"" + strID + "\"}")
}

// Метод для запроса задачи по id
func getTaskById(r *http.Request) []byte {

	// Проверка аргумента id в ссылке
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("{\"error\":\"Задача не найдена\"}")
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	// Выполнение запроса в базу по аргументу из ссылки
	respTask, err := db.DbInstance.GetTaskByID(id)
	if err != nil {
		log.Println("{\"error\":\"Задача не найдена\"}", err.Error())
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	// Сериализация JSON
	res, err := json.Marshal(respTask)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

// Метод для изменения задачи
func changeTask(r *http.Request) []byte {
	newTask := &db.Task{}

	//TODO don't like how here callback bytes res and error

	// Проверка вводных данных задачи
	res, err := newTask.CheckTask(r)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	fmt.Println(newTask)

	// Добавление задачи в базу
	_, err = db.DbInstance.UpateTask(newTask)
	if err != nil {
		log.Println("{\"error\":\"Не удалось обновить запись в базе\"}", err.Error())
		return []byte("{\"error\":\"Не удалось обновить запись в базе\"}")
	}

	return []byte("{}")
}

// Обработчик для ручки api/nextDate
func NextDate(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	now, err := time.Parse("20060102", values.Get("now"))
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := tasks.NextDateHandler(now, values.Get("date"), values.Get("repeat"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ANSWER IS:", res)
	w.Write([]byte(res))
}

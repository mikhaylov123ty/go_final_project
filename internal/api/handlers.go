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

	// Запрос POST
	case r.Method == "POST":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(addTask(r))

	// Запрос по id задачи
	case r.URL.Query().Has("id") == true:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getTaskById(r))

	// По умолчанию возвращает статус с ошибкой
	default:
		http.Error(w, "", http.StatusBadRequest)
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

	// Десериализация JSON
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Println("{\"error\":\"ошибка десериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка десериализации JSON\"}")
	}

	// Обработка условия, если обязательно поле "Дата" оказалось пустым
	// Установка текущей даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	}

	// Обработка условия, если обязательно поле "Заголовок" оказалось пустым
	// Возврат ошибки
	if newTask.Title == "" {
		log.Println("{\"error\":\"Не указан заголовок задачи\"}")
		return []byte("{\"error\":\"Не указан заголовок задачи\"}")
	}

	// Проверка соответствия формата даты
	taskDate, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		log.Println("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}", err.Error())
		return []byte("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}")
	}

	// Обработка условия, если текущая дата меньше даты задачи
	if taskDate.Before(time.Now().UTC().Round(24*time.Hour).AddDate(0, 0, -1)) {

		// Обработка условия с повторением, если он не пустой и поиск следующей даты от даты повторения
		if newTask.Repeat != "" {
			newTask.Date, err = tasks.NextDateHandler(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				log.Println("{\"error\":\"" + err.Error() + "}")
				return []byte("{\"error\":\"" + err.Error() + "\"}")
			}

			// Установка текущей даты, если условия повторения нет
		} else {
			newTask.Date = time.Now().Format("20060102")
		}
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

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

// Основной обработчик для ручки task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(addTask(w, r))
		//default:
		//	getTaskHandler(w, r)
	}
}

// Основной обработчик для ручки tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Has("search") {
	case true:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getTask(r))
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getAllTasks())
	}
}

// Метод для запроса всех задач
func getAllTasks() []byte {
	var err error
	newResponse := &response{}

	newResponse.Tasks, err = db.DbInstance.GetAllTasks()
	if err != nil {
		log.Println("{\"error\":\"ошибка запроса в базу\"}", err.Error())
		return []byte("{\"error\":\"ошибка запроса в базу\"}")
	}

	res, err := json.Marshal(newResponse)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

// Метод для запроса задачи по поиску
func getTask(r *http.Request) []byte {
	var err error
	newResponse := &response{}

	searchQuery := r.URL.Query().Get("search")
	fmt.Println("search query:", searchQuery)

	newResponse.Tasks, err = db.DbInstance.GetTaskBySearch(r.URL.Query().Get("search"))
	if err != nil {
		log.Println("{\"error\":\"ошибка запроса в базу\"}", err.Error())
		return []byte("{\"error\":\"ошибка запроса в базу\"}")
	}

	res, err := json.Marshal(newResponse)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

// Метод для добавления задачи в базу
func addTask(w http.ResponseWriter, r *http.Request) []byte {
	newTask := &db.Task{}

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Println("{\"error\":\"ошибка десериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка десериализации JSON\"}")
	}

	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	}

	if newTask.Title == "" {
		log.Println("{\"error\":\"Не указан заголовок задачи\"}")
		return []byte("{\"error\":\"Не указан заголовок задачи\"}")
	}

	taskDate, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		log.Println("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}", err.Error())
		return []byte("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}")
	}

	if taskDate.Before(time.Now().UTC().Round(24*time.Hour).AddDate(0, 0, -1)) {
		if newTask.Repeat != "" {
			newTask.Date, err = tasks.NextDate(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				log.Println("{\"error\":\"" + err.Error() + "}")
				return []byte("{\"error\":\"" + err.Error() + "\"}")
			}
		} else {
			newTask.Date = time.Now().Format("20060102")
		}
	}

	id, err := db.DbInstance.AddTask(newTask)
	if err != nil {
		log.Println("{\"error\":\"Не удалось добавить в базу\"}", err.Error())
		return []byte("{\"error\":\"Не удалось добавить в базу\"}")
	}
	strID := strconv.Itoa(id)

	return []byte("{\"id\":\"" + strID + "\"}")
}

// Обработчик для nextDate запросов
func NextDate(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	now, err := time.Parse("20060102", values.Get("now"))
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := tasks.NextDate(now, values.Get("date"), values.Get("repeat"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ANSWER IS:", res)
	w.Write([]byte(res))
}

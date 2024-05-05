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
	Id    int   `json:"id,omitempty"`
	Error error `json:"error,omitempty"`
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
		getTask(w, r)
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(getAllTasks(w, r))
	}
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

// Метод для запроса всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) []byte {
	result, err := db.DbInstance.GetAllTasks()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("JSON", string(res))

	return res
}

// Метод для запроса задачи по id
func getTask(w http.ResponseWriter, r *http.Request) {

}

// Метод для добавления задачи в базу
func addTask(w http.ResponseWriter, r *http.Request) []byte {

	newTask := &db.Task{}
	//resp := response{}

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
		log.Println(err)
	}
	strID := strconv.Itoa(id)

	return []byte("{\"id\":\"" + strID + "\"}")
}

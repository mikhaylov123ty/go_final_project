package api

import (
	"encoding/json"
	"finalProject/internal/db"
	"finalProject/internal/tasks"
	"fmt"
	"log"
	"net/http"
	"time"
)

type response struct {
	Id    int   `json:"id,omitempty"`
	Error error `json:"error,omitempty"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		addTask(w, r)
		//default:
		//	getTaskHandler(w, r)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Has("search") {
	case true:
		getTask(w, r)
	default:
		getAllTasks(w, r)
	}
}

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

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	result, err := db.DbInstance.GetAllTasks()
	if err != nil {
		fmt.Println(err)
		return
	}

	json, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("JSON", string(json))
	_, err = w.Write(json)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {

}

func addTask(w http.ResponseWriter, r *http.Request) {

	newTask := &db.Task{}
	resp := response{}

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {

		log.Println("{\"error\":\"ошибка десериализации JSON\"}", err.Error())
		w.Write([]byte("{\"error\":\"ошибка десериализации JSON\"}"))
		return
	}

	if newTask.Title == "" {
		log.Println("{\"error\":\"Не указан заголовок задачи\"}")
		w.Write([]byte("{\"error\":\"Не указан заголовок задачи\"}"))
		return
	}

	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	}

	taskDate, err := time.Parse("20060102", newTask.Date)
	if err != nil {
		log.Println("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}", err.Error())
		w.Write([]byte("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}"))
		return
	}

	//TODO think about today case
	if taskDate.Before(time.Now()) || taskDate.Format("20060102") != time.Now().Format("20060102") {
		if newTask.Repeat != "" {
			newTask.Date, err = tasks.NextDate(time.Now(), newTask.Date, newTask.Repeat)
			if err != nil {
				log.Println("{\"error\":\"" + err.Error() + "}")
				w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
				return
			}
		} else {
			newTask.Date = time.Now().Format("20060102")
		}
	}

	resp.Id, resp.Error = db.DbInstance.AddTask(newTask)

	json, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshall", err)
		return
	}
	fmt.Println("RESPO", resp)
	fmt.Println("JSON", string(json))
	w.Write(json)
}

//func (r *response) constructResponse() []byte {
//
//}

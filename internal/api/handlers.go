package api

import (
	"encoding/json"
	"finalProject/internal/db"
	"finalProject/internal/tasks"
	"fmt"
	"net/http"
	"time"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		addTask(w, r)
		//default:
		//	getTaskHandler(w, r)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
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
	fmt.Println("ADD task")
	newTask := &db.Task{}

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		fmt.Println("error unmarshalling task")
		return
	}

	err = db.DbInstance.AddTask(newTask)
	fmt.Println(newTask)
}

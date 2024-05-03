package api

import (
	"finalProject/internal/db"
	"finalProject/internal/tasks"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type task struct {
	date    string `json:"date"`
	title   string `json:"title"`
	comment string `json:"comment,omitempty"`
	repeat  string `json:"repeat,omitempty"`
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "nextdate"):
		nextDate(w, r)
	case strings.Contains(r.URL.Path, "tasks"):
		taskHandler(w, r)
	default:
		http.NotFound(w, r)
	}

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		addTask(w, r)
	default:
		getTask(w, r)
	}
}

func nextDate(w http.ResponseWriter, r *http.Request) {
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

func getTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET task")
	res, err := db.DbInstance.Connection.Query("SELECT * FROM scheduler")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Close()
	result := &task{}
	for res.Next() {
		res.Scan(result)
	}

	fmt.Println(result)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ADD task")

	res, err := db.DbInstance.Connection.Exec("INSERT INTO scheduler VALUES (3,20060103,'test','testcoment','')")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.LastInsertId())
}

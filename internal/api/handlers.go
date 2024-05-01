package api

import (
	"finalProject/internal/tasks"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "nextdate"):
		nextDate(w, r)
	default:
		http.NotFound(w, r)
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

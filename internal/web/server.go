package web

import (
	"log"
	"net/http"

	"finalProject/internal/api"
)

type server struct {
	port string
}

func Init(port string) *server {
	return &server{
		port: port,
	}
}

func (s *server) Start(webDir string) error {
	log.Println("Starting server...")

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/nextdate", api.NextDate)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)

	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		return err
	}

	return nil
}

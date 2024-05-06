package web

import (
	"log"
	"net/http"

	"finalProject/internal/api"
)

// Структура веб-сервера
type server struct {
	port string
}

// Инициализация веб-сервера
func Init(port string) *server {
	return &server{
		port: port,
	}
}

// Запуск веб-сервера
func (s *server) Start(webDir string) error {
	log.Println("Starting server...")

	// Обработчик файлового сервера
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Перечисление ручек и их обработок
	http.HandleFunc("/api/nextdate", api.NextDate)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)

	// Запуск
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		return err
	}

	return nil
}

package web

import (
	"fmt"
	"net/http"

	"finalProject/internal/api"
	"finalProject/internal/api/handlers"
	"finalProject/internal/logger"
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
	logger.Slog.Json.Info("Starting server...")
	logger.Slog.Json.Info(fmt.Sprintf("Server: http://localhost:%s/", s.port))

	// Обработчик файлового сервера
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Перечисление ручек и их обработок
	http.HandleFunc("/api/nextdate", api.NextDate)

	http.HandleFunc("/api/task", handlers.Auth(api.TaskHandler))
	http.HandleFunc("/api/task/done", handlers.Auth(api.TaskDoneHandler))
	http.HandleFunc("/api/tasks", handlers.Auth(api.TasksHandler))
	http.HandleFunc("/api/signin", api.SignHandler)

	// Запуск
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		return err
	}

	return nil
}

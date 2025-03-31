package main

import (
	"log"
	"os"

	"finalProject/internal/db"
	"finalProject/internal/logger"
	"finalProject/internal/web"
)

func main() {
	var port string
	var dbFile string

	// Инициализация логгера Slog
	logger.Init()

	// Предопределение, если переменные окружения пустые
	if len(os.Getenv("TODO_DBFILE")) > 0 {
		dbFile = os.Getenv("TODO_DBFILE")
	} else {
		dbFile = "scheduler.db"
	}

	// Инициализация базы
	dbInstance, err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Закрытие подключения
	defer dbInstance.Connection.Close()

	// Предопределение, если переменные окружения пустые
	if len(os.Getenv("TODO_PORT")) > 0 {
		port = os.Getenv("TODO_PORT")

	} else {
		port = "7540"
	}

	// Инициализация веб-сервера
	server := web.Init(port)
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}
}

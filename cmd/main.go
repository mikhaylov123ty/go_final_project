package main

import (
	"fmt"
	"log"
	"os"

	"finalProject/internal/db"
	"finalProject/internal/web"
)

func main() {
	fmt.Println(os.Getenv("TODO_DBFILE"))
	// Инициализация базы
	dbInstance, err := db.Init(os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Закрытие подключения
	defer dbInstance.Connection.Close()

	// Инициализация веб-сервера
	server := web.Init(os.Getenv("TODO_PORT"))
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}
}

package main

import (
	"finalProject/internal/db"
	"finalProject/internal/web"
	"log"
	"os"
)

func main() {
	dbInstance, err := db.Init(os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbInstance.Connection.Close()

	server := web.Init(os.Getenv("TODO_PORT"))
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}

}

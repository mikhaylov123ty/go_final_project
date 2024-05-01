package main

import (
	"finalProject/internal/dataBase"
	"finalProject/internal/webServer"
	"log"
	"os"
)

func main() {
	dbConnect, err := dataBase.Init(os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbConnect.DB.Close()

	server := webServer.Init(os.Getenv("TODO_PORT"))
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}

}

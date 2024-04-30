package main

import (
	"finalProject/internal/dataBase"
	"finalProject/internal/tasks"
	"finalProject/internal/webServer"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	dbConnect, err := dataBase.Init(os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbConnect.DB.Close()

	nowTime, err := time.Parse("20060102", "20240126")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(tasks.NextDate(nowTime, "20240113", "d 7"))
	fmt.Println(tasks.NextDate(nowTime, "16890220", "y"))
	fmt.Println(tasks.NextDate(nowTime, "20240126", "k 34"))
	fmt.Println(tasks.NextDate(nowTime, "20240126", "d 2"))

	server := webServer.Init(os.Getenv("TODO_PORT"))
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}

}

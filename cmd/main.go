package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"finalProject/internal/dataBase"
	"finalProject/internal/tasks"
	"finalProject/internal/webServer"
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

	fmt.Println(tasks.NextDate(nowTime, "20000113", "d 7"))
	fmt.Println(tasks.NextDate(nowTime, "16890220", "y"))
	fmt.Println(tasks.NextDate(nowTime, "20240126", "k 34"))
	fmt.Println(tasks.NextDate(nowTime, "20240126", "d 2"))
	fmt.Println(tasks.NextDate(nowTime, "20240125", "w 2,3"))
	fmt.Println(tasks.NextDate(nowTime, "20240126", "w 7"))
	fmt.Println(tasks.NextDate(nowTime, "20230126", "w 4,5"))
	fmt.Println(tasks.NextDate(nowTime, "20240329", "m 10,17 12,8,1"))
	fmt.Println(tasks.NextDate(nowTime, "16890811", "m 1,27 12,8,1"))
	fmt.Println(tasks.NextDate(nowTime, "20241229", "m 1,27 12,8,1"))
	fmt.Println(tasks.NextDate(nowTime, "20240326", "m -1,-2"))

	fmt.Println(tasks.NextDate(nowTime, "20240409", "m 31"))
	fmt.Println(tasks.NextDate(nowTime, "20240201", "m 18,-1,-2"))

	server := webServer.Init(os.Getenv("TODO_PORT"))
	err = server.Start("web")
	if err != nil {
		log.Fatal(err.Error())
	}

}

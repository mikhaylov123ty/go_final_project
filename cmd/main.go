package main

import (
	"os"

	"finalProject/internal/webServer"
)

func main() {
	newServer := webServer.Init(os.Getenv("webServerPort"))
	newServer.Start("web/")
}

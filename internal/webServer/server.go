package webServer

import (
	"log"
	"net/http"
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
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		return err
	}
	return nil
}

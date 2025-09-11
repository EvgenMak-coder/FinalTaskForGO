package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	pathToWebDir = "./web"
	defaultPort  = "7540"
)

func Run() {
	webDir := os.Getenv("TODO_WEB_DIR")
	if webDir == "" {
		webDir = pathToWebDir
	}
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatal("dir WEB can't find: ", err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	fmt.Println("Server started at :", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server crashed err: ", err)
	}

}

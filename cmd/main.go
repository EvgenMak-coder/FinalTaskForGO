package main

import (
	"log"
	"net/http"
	"os"
)

const defaultDir = "./web"

func main() {
	webDir := os.Getenv("TODO_WEB_DIR")
	if webDir == "" {
		webDir = defaultDir
	}
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Web directory does not exist")
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal("Server crashed", err)
	}
}

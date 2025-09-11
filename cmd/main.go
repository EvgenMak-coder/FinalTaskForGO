package main

import (
	"FinalTaskForGO/pkg/api"
	"FinalTaskForGO/pkg/db"
	"FinalTaskForGO/pkg/server"
	"log"
	"os"
)

const defaultDBfile = "./scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBfile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}
	defer db.GetDB().Close()

	api.Init()
	server.Run()

}

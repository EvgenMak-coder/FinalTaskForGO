package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL,
		title VARCHAR(256) NOT NULL, 
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT ""	
	);
	CREATE INDEX IF NOT EXIST idx_date ON scheduler(date);
	`

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install {
		_, err = database.Exec(schema)
		if err != nil {
			return err
		}
	}
	DB = database

	return nil
}

func GetDB() *sql.DB {
	return DB
}

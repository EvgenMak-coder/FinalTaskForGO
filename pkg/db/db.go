package db

import (
	"database/sql"
	"fmt"
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
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	if install {
		if _, err = database.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}
	DB = database
	return nil
}

func GetDB() *sql.DB {
	return DB
}

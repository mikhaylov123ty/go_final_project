package dataBase

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type dataBase struct {
	DB *sql.DB
}

func Init(file string) (*dataBase, error) {
	log.Println("Initializing database")
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	fmt.Println("Need to create db? ", checkFile())
	if checkFile() {
		_, err = db.Exec("CREATE TABLE scheduler (" +
			"id INTEGER PRIMARY KEY AUTOINCREMENT," +
			"date INTEGER NOT NULL," +
			"title TEXT NOT NULL," +
			"comment TEXT NOT NULL," +
			"repeat TEXT(128) NOT NULL" +
			");")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE INDEX scheduler_id_IDX ON scheduler (date);")
		if err != nil {
			return nil, err
		}
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &dataBase{DB: db}, nil
}

func checkFile() bool {
	_, err := os.Stat(os.Getenv("TODO_DBFILE"))

	var install bool
	if err != nil {
		install = true
	}
	return install
}

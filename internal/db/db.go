package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type dbInstance struct {
	Connection *sql.DB
}

var DbInstance *dbInstance

func Init(file string) (*dbInstance, error) {
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
			"comment TEXT," +
			"repeat TEXT(128)" +
			");")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE INDEX scheduler_id_IDX ON scheduler (id);")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("CREATE INDEX scheduler_date_IDX ON scheduler (date);")
		if err != nil {
			return nil, err
		}
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	DbInstance = &dbInstance{Connection: db}

	return DbInstance, nil
}

func checkFile() bool {
	_, err := os.Stat(os.Getenv("TODO_DBFILE"))

	var install bool
	if err != nil {
		install = true
	}
	return install
}

func (db *dbInstance) GetAllTasks() ([]*Task, error) {
	res, err := db.Connection.Query("SELECT * FROM scheduler")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Close()

	result := []*Task{}
	for res.Next() {
		row := &Task{}
		err = res.Scan(&row.Id, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		result = append(result, row)
	}
	return result, nil
}

func (db *dbInstance) AddTask(t *Task) error {
	res, err := db.Connection.Exec(
		"INSERT INTO scheduler (date, title,comment,repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res.LastInsertId())
	return nil
}

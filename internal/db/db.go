package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type dbInstance struct {
	Connection *sql.DB
}

var DbInstance *dbInstance

// Метод инициализации файла базы данных
// file - путь к файлу с базой
func Init(file string) (*dbInstance, error) {
	log.Println("Initializing database")

	// Открываем\создаем файл с базой данных
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	// Проверка наличия файла с базой данных
	if checkFile() {
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createIdIndex)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createDateIndex)
		if err != nil {
			return nil, err
		}
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Передача инстанса в общую переменную
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

	// Выполнение запроса к базе
	res, err := db.Connection.Query("SELECT * FROM scheduler ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Упаковка результатов в слайс адресов задач
	result := []*Task{}
	for res.Next() {
		row := &Task{}
		err = res.Scan(&row.Id, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// Метод для поиска задач
// search - текст, который вводится в поисковую строку
func (db *dbInstance) GetTaskBySearch(search string) ([]*Task, error) {

	// Парсинг вероятной даты
	possibleDate, err := time.Parse("02.01.2006", search)
	fmt.Println("DATE", possibleDate)

	// Выполнение запроса в базу
	res, err := db.Connection.Query(`SELECT * FROM scheduler WHERE 
		id = :id
		OR title LIKE :search
		OR comment LIKE :search
		OR date = :possibleDate
		ORDER BY date;`,
		sql.Named("id", search),
		sql.Named("search", "%"+search+"%"),
		sql.Named("possibleDate", possibleDate.Format("20060102")),
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Упаковка результата в слайс адресов задач
	result := []*Task{}
	for res.Next() {
		row := &Task{}
		err = res.Scan(&row.Id, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// Метод для поиска задачи по id
func (db *dbInstance) GetTaskByID(id string) (*Task, error) {
	res := &Task{}

	// Выполнение запроса к базе
	row := db.Connection.QueryRow(`SELECT * FROM scheduler WHERE
		id = :id;`,
		sql.Named("id", id))

	// Сканирование строки
	err := row.Scan(&res.Id, &res.Date, &res.Title, &res.Comment, &res.Repeat)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Метод для добавления задачи
// t - экземпляр структуры Task из models
func (db *dbInstance) AddTask(t *Task) (int, error) {

	// Выполнение запроса к базе
	exec, err := db.Connection.Exec(
		"INSERT INTO scheduler (date, title,comment,repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return 0, err
	}

	// Передача последнего id записи
	res, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

func (db *dbInstance) UpateTask(t *Task) (int, error) {

	// Проверка, что такой id существует
	_, err := db.GetTaskByID(t.Id)
	if err != nil {
		return 0, err
	}

	// Выполнение запроса к базе
	exec, err := db.Connection.Exec(
		`UPDATE scheduler 
    SET date = :date, title =:title,comment = :comment,repeat = :repeat
WHERE id = :id;`,
		sql.Named("id", t.Id),
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return 0, err
	}

	// Передача последнего id записи
	res, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

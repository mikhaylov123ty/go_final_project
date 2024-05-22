package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"finalProject/internal/models"

	_ "modernc.org/sqlite"
)

// Определение констант для запросов создания таблицы и индексов
const (
	createTableQuery = `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date INTEGER NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT(128)
		);`

	createIdIndex   = `CREATE INDEX scheduler_id_IDX ON scheduler (id);`
	createDateIndex = `CREATE INDEX scheduler_date_IDX ON scheduler (date);`
)

// Структура подключения к БД
type dbInstance struct {
	Connection *sql.DB
}

// Общая переменная инстанса подключения к БД
var DbInstance *dbInstance

// Метод инициализации файла БД
// file - путь к файлу с БД
func Init(file string) (*dbInstance, error) {
	log.Println("Initializing database")

	// Открываем\создаем файл с базой данных
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	// Проверка наличия таблицы в БД, создание таблицы и индексов в случае отсутствия
	if checkDbTable(db) {
		log.Println("Table scheduler not found, creating")
		_, err := db.Exec(createTableQuery)
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

	// Логирование пути файла с БД
	exec, _ := os.Executable()
	log.Println("Database initialized, path:", exec+"/"+file)

	// Передача инстанса в общую переменную
	DbInstance = &dbInstance{Connection: db}

	return DbInstance, nil
}

// Метод для проверки наличия таблицы в файле БД
func checkDbTable(db *sql.DB) bool {

	// Производится запрос в БД на наличи таблицы
	rows, err := db.Query(`SELECT name 
		FROM sqlite_master 
		WHERE type='table' AND name='scheduler';`)

	// Закрыть строки после проверки
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	//Если строки нет - необходимо создать таблицу
	if !rows.Next() {
		return true
	}

	return false
}

// Метод для запроса в БД и вывода всех задач
func (db *dbInstance) GetAllTasks() ([]*models.Task, error) {

	// Выполнение запроса к базе
	res, err := db.Connection.Query(
		`SELECT * FROM scheduler
         		ORDER BY date;`,
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Упаковка результатов в слайс адресов задач
	result := []*models.Task{}
	for res.Next() {
		row := &models.Task{}
		err = res.Scan(&row.Id, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// Метод для поиска задач в БД
// search - текст, который вводится в поисковую строку
func (db *dbInstance) GetTaskBySearch(search string) ([]*models.Task, error) {

	// Парсинг вероятной даты
	possibleDate, err := time.Parse("02.01.2006", search)

	// Выполнение запроса в базу
	res, err := db.Connection.Query(
		`SELECT * FROM scheduler
				WHERE id = :id
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
	result := []*models.Task{}
	for res.Next() {
		row := &models.Task{}
		err = res.Scan(&row.Id, &row.Date, &row.Title, &row.Comment, &row.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// Метод для поиска задачи в БД по id
// id - идентификатор задачи
func (db *dbInstance) GetTaskByID(id string) (*models.Task, error) {
	res := &models.Task{}

	// Выполнение запроса к базе
	row := db.Connection.QueryRow(
		`SELECT * FROM scheduler 
         		WHERE id = :id;`,
		sql.Named("id", id))

	// Сканирование строки
	err := row.Scan(&res.Id, &res.Date, &res.Title, &res.Comment, &res.Repeat)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Метод для добавления задачи в БД
// t - адрес экземпляра структуры Task из models
func (db *dbInstance) AddTask(t *models.Task) (int, error) {

	// Выполнение запроса к базе
	exec, err := db.Connection.Exec(
		`INSERT INTO scheduler (date, title,comment,repeat)
				VALUES (:date, :title, :comment, :repeat)`,
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

// Метод для обновления задачи в БД
// t - адрес экземпляра структуры Task из models
func (db *dbInstance) UpateTask(t *models.Task) (int, error) {

	// Выполнение запроса к базе
	exec, err := db.Connection.Exec(
		`UPDATE scheduler 
				SET date = :date,
					title =:title,
					comment = :comment,
					repeat = :repeat
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

// Метод для удаления задачи из БД
// id - идентификатор задачи
func (db *dbInstance) DeleteTask(id string) error {

	// Выполнение запроса к базе
	_, err := db.Connection.Exec(
		`DELETE FROM scheduler
       WHERE id = :id;`,
		sql.Named("id", id),
	)

	if err != nil {
		return err
	}

	return nil
}

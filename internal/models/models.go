package models

// Константы для создания таблицы и индексов
const (
	CreateTableQuery = `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date INTEGER NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT(128)
		);`
	CreateIdIndex   = `CREATE INDEX scheduler_id_IDX ON scheduler (id);`
	CreateDateIndex = `CREATE INDEX scheduler_date_IDX ON scheduler (date);`

	authenticationRequired = "Authentication required"
	IncorrectRequest       = "Не корректный запрос"
)

// Структура для задач
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Структура для запросов
type Request struct {
	Password string `json:"password,omitempty"`
}

// Структура для ответов
type Response struct {
	Tasks []*Task `json:"tasks,omitempty"`
	Id    int     `json:"id,omitempty"`
	Error string  `json:"error,omitempty"`
	Token string  `json:"token,omitempty"`
}

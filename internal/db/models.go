package db

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

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

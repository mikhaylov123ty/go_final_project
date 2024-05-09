package db

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"finalProject/internal/tasks"
)

// Константы для создания таблицы и индексов
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

// Структура для задач
type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Структура для ответов
type Response struct {
	Tasks    []*Task `json:"tasks"`
	Password string  `json:"password,omitempty"`
}

func (t *Task) CheckTask(r *http.Request) ([]byte, error) {
	// Десериализация JSON
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return []byte("{\"error\":\"ошибка десериализации JSON\"}"), err
	}

	// Обработка условия, если обязательно поле "Дата" оказалось пустым
	// Установка текущей даты
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}

	// Обработка условия, если обязательно поле "Заголовок" оказалось пустым
	// Возврат ошибки
	if t.Title == "" {
		return []byte("{\"error\":\"Не указан заголовок задачи\"}"), errors.New("{\"error\":\"Не указан заголовок задачи\"}")
	}

	// Проверка соответствия формата даты
	taskDate, err := time.Parse("20060102", t.Date)
	if err != nil {
		return []byte("{\"error\":\"Дата представлена в формате, отличном от 20060102\"}"), err
	}

	// Обработка условия с повторением, если он не пустой и поиск следующей даты от даты повторения
	if t.Repeat != "" {
		tempDateStr, err := tasks.NextDateHandler(time.Now(), t.Date, t.Repeat)
		if err != nil {
			errStr := strings.Replace(err.Error(), "\"", "", -1)
			return []byte("{\"error\":\"" + errStr + "\"}"), err
		}
		if taskDate.Before(time.Now().AddDate(0, 0, -1)) {
			t.Date = tempDateStr
		}
	}

	// Обработка условия если дата пуста или меньше сегодня
	taskDate, err = time.Parse("20060102", t.Date)
	if t.Date == "" || taskDate.Before(time.Now()) {
		t.Date = time.Now().Format("20060102")
	}

	return nil, nil
}

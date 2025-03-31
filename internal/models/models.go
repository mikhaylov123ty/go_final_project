package models

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

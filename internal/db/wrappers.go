package db

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"finalProject/internal/tasks"
)

func (r *Response) Marshal() []byte {
	res, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshalling response:", err)
	}
	return res
}

func (r *Response) LogResponseError(s string) []byte {
	r.Error = s
	log.Println(r.Error)
	return r.Marshal()
}

func (t *Task) CheckTask(r *http.Request) *Response {
	response := &Response{}
	// Десериализация JSON
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		response.Error = err.Error()
		return response
	}

	// Обработка условия, если обязательно поле "Дата" оказалось пустым
	// Установка текущей даты
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}

	// Обработка условия, если обязательно поле "Заголовок" оказалось пустым
	// Возврат ошибки
	if t.Title == "" {
		response.Error = "Не указан заголовок задачи"
		return response
	}

	// Проверка соответствия формата даты
	taskDate, err := time.Parse("20060102", t.Date)
	if err != nil {
		response.Error = "Дата представлена в формате, отличном от 20060102"
		return response
	}

	// Обработка условия с повторением, если он не пустой и поиск следующей даты от даты повторения
	if t.Repeat != "" {
		tempDateStr, err := tasks.NextDateHandler(time.Now(), t.Date, t.Repeat)
		if err != nil {
			errStr := strings.Replace(err.Error(), "\"", "", -1)
			response.Error = errStr
			return response
		}
		if taskDate.Before(time.Now().AddDate(0, 0, -1)) {
			t.Date = tempDateStr
		}
	}

	// Обработка условия если дата пуста или меньше сегодня
	taskDate, _ = time.Parse("20060102", t.Date)
	if t.Date == "" || taskDate.Before(time.Now()) {
		t.Date = time.Now().Format("20060102")
	}

	return response
}

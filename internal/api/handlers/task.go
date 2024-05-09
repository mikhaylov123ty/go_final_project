package handlers

import (
	"encoding/json"
	"finalProject/internal/tasks"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"finalProject/internal/db"
)

// Метод для добавления задачи в базу
func AddTask(r *http.Request) []byte {
	newTask := &db.Task{}

	//TODO don't like how here callback bytes res and error

	// Проверка вводных данных задачи
	res, err := newTask.CheckTask(r)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	// Добавление задачи в базу
	id, err := db.DbInstance.AddTask(newTask)
	if err != nil {
		log.Println("{\"error\":\"Не удалось добавить в базу\"}", err.Error())
		return []byte("{\"error\":\"Не удалось добавить в базу\"}")
	}

	strID := strconv.Itoa(id)

	return []byte("{\"id\":\"" + strID + "\"}")
}

// Метод для изменения задачи
func ChangeTask(r *http.Request) []byte {
	modifiedTask := &db.Task{}

	// Проверка вводных данных задачи
	res, err := modifiedTask.CheckTask(r)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	_, err = db.DbInstance.GetTaskByID(modifiedTask.Id)
	if err != nil {
		log.Println("{\"error\":\"Задача не найдена\"}", err.Error())
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	// Добавление задачи в базу
	_, err = db.DbInstance.UpateTask(modifiedTask)
	if err != nil {
		log.Println("{\"error\":\"Не удалось обновить запись в базе\"}", err.Error())
		return []byte("{\"error\":\"Не удалось обновить запись в базе\"}")
	}

	return []byte("{}")
}

// Метод для удаления задачи
func DeleteTaskById(r *http.Request) []byte {
	taskID := r.URL.Query().Get("id")

	task, err := db.DbInstance.GetTaskByID(taskID)
	if err != nil {
		log.Println("{\"error\":\"Задача не найдена\"}", err.Error())
		return []byte("{\"error\":\"Задача не найдена\"}")
	}
	err = db.DbInstance.DeleteTask(task.Id)
	if err != nil {
		log.Println("{\"error\":\"Не удалось удалить запись в базе\"}", err.Error())
		return []byte("{\"error\":\"Не удалось удалить запись в базе\"}")
	}

	return []byte("{}")
}

// Метод для запроса задачи по id
func GetTaskById(r *http.Request) []byte {

	// Проверка аргумента id в ссылке
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("{\"error\":\"Задача не найдена\"}")
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	// Выполнение запроса в базу по аргументу из ссылки
	respTask, err := db.DbInstance.GetTaskByID(id)
	if err != nil {
		log.Println("{\"error\":\"Задача не найдена\"}", err.Error())
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	// Сериализация JSON
	res, err := json.Marshal(respTask)
	if err != nil {
		log.Println("{\"error\":\"ошибка сериализации JSON\"}", err.Error())
		return []byte("{\"error\":\"ошибка сериализации JSON\"}")
	}

	return res
}

// Метод для завершения задачи
func DoneTask(r *http.Request) []byte {
	taskID := r.URL.Query().Get("id")

	task, err := db.DbInstance.GetTaskByID(taskID)
	if err != nil {
		log.Println("{\"error\":\"Задача не найдена\"}", err.Error())
		return []byte("{\"error\":\"Задача не найдена\"}")
	}

	if task.Repeat != "" {
		task.Date, err = tasks.NextDateHandler(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errStr := strings.Replace(err.Error(), "\"", "", -1)
			return []byte("{\"error\":\"" + errStr + "\"}")
		}
		_, err = db.DbInstance.UpateTask(task)
		if err != nil {
			log.Println("{\"error\":\"Не удалось обновить запись в базе\"}", err.Error())
			return []byte("{\"error\":\"Не удалось обновить запись в базе\"}")
		}
		return []byte("{}")
	}

	err = db.DbInstance.DeleteTask(task.Id)
	if err != nil {
		log.Println("{\"error\":\"Не удалось обновить запись в базе\"}", err.Error())
		return []byte("{\"error\":\"Не удалось обновить запись в базе\"}")
	}

	return []byte("{}")
}

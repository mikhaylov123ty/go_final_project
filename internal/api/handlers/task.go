package handlers

import (
	"encoding/json"
	"finalProject/internal/models"
	"log"
	"net/http"
	"time"

	"finalProject/internal/db"
	"finalProject/internal/tasks"
)

// Метод для добавления задачи в базу
func AddTask(r *http.Request) []byte {
	var err error
	newTask := &models.Task{}

	// Проверка вводных данных задачи
	response := newTask.CheckTask(r)
	if response.Error != "" {
		log.Println(response.Error)
		return response.Marshal()
	}

	// Добавление задачи в БД
	response.Id, err = db.DbInstance.AddTask(newTask)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	return response.Marshal()
}

// Метод для изменения задачи
func ChangeTask(r *http.Request) []byte {
	var err error
	modifiedTask := &models.Task{}

	// Проверка вводных данных задачи
	response := modifiedTask.CheckTask(r)
	if response.Error != "" {
		log.Println(response.Error)
		return response.Marshal()
	}

	// Запрос в БД по id задачи
	_, err = db.DbInstance.GetTaskByID(modifiedTask.Id)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Изменение задачи в базе
	_, err = db.DbInstance.UpateTask(modifiedTask)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	return []byte("{}")
}

// Метод для удаления задачи
func DeleteTaskById(r *http.Request) []byte {
	taskID := r.URL.Query().Get("id")
	response := &models.Response{}

	// Проверка существования задачи по id в БД
	task, err := db.DbInstance.GetTaskByID(taskID)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Удаление задачи по id из БД
	err = db.DbInstance.DeleteTask(task.Id)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	return []byte("{}")
}

// Метод для запроса задачи по id
func GetTaskById(r *http.Request) []byte {
	response := &models.Response{}

	// Проверка аргумента id в ссылке
	id := r.URL.Query().Get("id")
	if id == "" {
		return response.LogResponseError("Задача не найдена")
	}

	// Выполнение запроса в базу по аргументу из ссылки
	respTask, err := db.DbInstance.GetTaskByID(id)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Сериализация JSON
	res, err := json.Marshal(respTask)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	return res
}

// Метод для завершения задачи
func DoneTask(r *http.Request) []byte {
	response := &models.Response{}
	taskID := r.URL.Query().Get("id")

	// Проверка существования задачи по id в БД
	task, err := db.DbInstance.GetTaskByID(taskID)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	// Проверка наличия повторения у задачи
	if task.Repeat != "" {

		// Вычисление следующей даты
		task.Date, err = tasks.NextDateHandler(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return response.LogResponseError(err.Error())
		}

		// Обновление задачи в БД
		_, err = db.DbInstance.UpateTask(task)
		if err != nil {
			return response.LogResponseError(err.Error())
		}

		return []byte("{}")
	}

	// Удаление задачи
	err = db.DbInstance.DeleteTask(task.Id)
	if err != nil {
		return response.LogResponseError(err.Error())
	}

	return []byte("{}")
}

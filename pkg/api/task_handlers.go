package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"FinalTaskForGO/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

// в функции getTaskHandler изменил обработку ошибок и добавил вывод статуса запроса
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "task not found" {
			status = http.StatusNotFound
		}
		writeJson(w, map[string]string{"error": err.Error()}, status)
		return
	}
	writeJson(w, task, http.StatusOK)
}

// в функции updateTaskHandler изменил обработку ошибок и добавил вывод статуса запроса
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Invalid JSON format"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(task.ID)
	if err != nil || id == 0 {
		writeJson(w, map[string]string{"error": "Task ID is required"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Task title is required"}, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid date format, expected YYYYMMDD"}, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		if afterNow(now, t) {
			task.Date = next
		}
	} else if afterNow(now, t) {
		task.Date = now.Format(dateFormat)
	}

	if err := db.UpdateTask(&task); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "task not found" {
			status = http.StatusNotFound
		}
		writeJson(w, map[string]string{"error": err.Error()}, status)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// в функции deleteTaskHandler изменил обработку ошибок и добавил вывод статуса запроса
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "task not found" {
			status = http.StatusNotFound
		}
		writeJson(w, map[string]string{"error": err.Error()}, status)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}

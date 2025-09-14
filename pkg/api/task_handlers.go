package api

import (
	"FinalTaskForGO/pkg/db"
	"encoding/json"
	"strconv"
	"time"

	"net/http"
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

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Invalid JSON format"})
		return
	}
	id, err := strconv.Atoi(task.ID)
	if err != nil || id == 0 {
		writeJson(w, map[string]string{"error": "Task ID is required"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Task title is required"})
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid date format, expected YYYYMMDD"})
		return
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		if afterNow(now, t) {
			task.Date = next
		}
	} else if afterNow(now, t) {
		task.Date = now.Format(dateFormat)
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{})
}

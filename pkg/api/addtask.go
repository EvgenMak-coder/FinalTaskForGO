package api

import (
	"FinalTaskForGO/pkg/db"
	"encoding/json"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Invalid JSON format"})
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Failed to add task to database"})
		return
	}

	writeJson(w, db.TaskResponse{ID: id})
}

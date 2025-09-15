package api

import (
	"net/http"
	"time"

	"FinalTaskForGO/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}

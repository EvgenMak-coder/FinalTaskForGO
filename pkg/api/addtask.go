package api

import (
	"FinalTaskForGO/pkg/db"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	layout := dateFormat
	if task.Date == "" {
		task.Date = now.Format(layout)
	}
	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format(layout)
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка чтения тела запроса"})
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка парсинга JSON"})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок"})
		return
	}
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"id": (fmtInt(id))})
}

func fmtInt(i int64) string {
	return strconv.FormatInt(i, 10)
}

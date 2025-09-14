package api

import (
	"FinalTaskForGO/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	lim := 50
	if tasks, err := db.Tasks(lim); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
	} else {
		writeJson(w, TasksResp{Tasks: tasks})
	}
}

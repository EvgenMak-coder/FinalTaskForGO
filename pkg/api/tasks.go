package api

import (
	"net/http"

	"FinalTaskForGO/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const limitTasks = 50

// в функции tasksHandler добавил обработку на метод Get, а также добавил вывод статуса запроса
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(limitTasks)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)

}

package api

import (
	"encoding/json"
	"net/http"
)

// в функцию writeJson добавил вывод о статусе запроса, обработку ошибки преобразования в JSON
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "не удалось преобразовать в JSON", http.StatusInternalServerError)
	}
}

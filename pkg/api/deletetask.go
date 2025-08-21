package api

import (
	"net/http"
	"todo-server/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errResp{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

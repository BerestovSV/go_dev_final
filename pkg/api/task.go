package api

import (
	"net/http"
)

func (a *API) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.addTaskHandler(w, r)
	case http.MethodGet:
		a.getTaskHandler(w, r)
	case http.MethodPut:
		a.updateTaskHandler(w, r)
	case http.MethodDelete:
		a.deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

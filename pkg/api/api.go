package api

import (
	"net/http"
	"todo-server/pkg/config"
	"todo-server/pkg/db"
)

type API struct {
	taskStore db.TaskStore
	config    *config.Config
}

func NewAPI(taskStore db.TaskStore, cfg *config.Config) *API {
	return &API{
		taskStore: taskStore,
		config:    cfg,
	}
}

func (a *API) Init() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/nextdate", a.nextDayHandler)
	router.HandleFunc("/api/task", a.authMiddleware(a.taskHandler))
	router.HandleFunc("/api/tasks", a.authMiddleware(a.tasksHandler))
	router.HandleFunc("/api/task/done", a.authMiddleware(a.doneTaskHandler))
	router.HandleFunc("/api/signin", a.signinHandler)

	return router
}

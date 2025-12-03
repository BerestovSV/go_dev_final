package api

import (
	"net/http"
	"todo-server/pkg/db"
)

func (a *API) tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		// Проверяем является ли поиск датой
		if checkSearchDate(search) {
			// Преобразуем дату в формат БД
			dbDate, err := convertSearchDateToDBFormat(search)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, errResp{Error: err.Error()})
				return
			}
			// Ищем задачи по дате
			tasks, err = a.taskStore.SearchTasksByDate(dbDate, 50)
		} else {
			// Ищем задачи по тексту
			tasks, err = a.taskStore.SearchTasksByText(search, 50)
		}
	} else {
		// Получаем все задачи
		tasks, err = a.taskStore.Tasks(50)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}

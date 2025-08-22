package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (a *API) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "Не указан идентификатор"})
		return
	}

	// Получаем задачу
	task, err := a.taskStore.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errResp{Error: err.Error()})
		return
	}

	// Если задача не повторяющаяся - удаляем
	if strings.TrimSpace(task.Repeat) == "" {
		err = a.taskStore.DeleteTask(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
			return
		}
	} else {
		// Если задача повторяющаяся - рассчитываем следующую дату
		now := time.Now()

		// Используем только дату без времени
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		nextDate, err := NextDate(today, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errResp{Error: fmt.Sprintf("Ошибка расчета следующей даты: %v", err)})
			return
		}

		// Обновляем дату задачи
		err = a.taskStore.UpdateDate(id, nextDate)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

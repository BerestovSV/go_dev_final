package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"todo-server/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errResp{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: fmt.Sprintf("read body error: %v", err)})
		return
	}
	defer r.Body.Close()

	var task db.Task
	if err = json.Unmarshal(body, &task); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: fmt.Sprintf("json decode error: %v", err)})
		return
	}

	// Проверяем, что ID указан
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "Не указан идентификатор задачи"})
		return
	}

	// 1) обязательное поле title
	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "не указан заголовок задачи"})
		return
	}

	// 2) проверяем и нормализуем дату (и правило повторения)
	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: err.Error()})
		return
	}

	// 3) обновляем задачу в БД
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errResp{Error: err.Error()})
		return
	}

	// Возвращаем пустой JSON объект
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"todo-server/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
		return
	}

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

	// 3) добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, idResp{ID: strconv.FormatInt(id, 10)})
}

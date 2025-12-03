package api

import (
	"fmt"
	"net/http"
	"time"
)

const DateFormat = "20060102"

func (a *API) nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "Параметры date и repeat обязательны", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}

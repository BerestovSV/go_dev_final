package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo-server/pkg/db"
)

type errResp struct {
	Error string `json:"error"`
}

type idResp struct {
	ID string `json:"id"`
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Вспомогательная функция для проверки, является ли дата сегодняшней
func isToday(date, now time.Time) bool {
	return date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()
}

// проверка и нормализация даты согласно условиям
func checkDate(task *db.Task) error {
	now := time.Now()

	// если дата пустая или в прошлом — используем сегодняшнюю
	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(DateFormat)
	}

	// проверяем формат даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты (ожидается %s)", DateFormat)
	}

	repeat := strings.TrimSpace(task.Repeat)
	if repeat != "" {
		// Для повторяющихся задач с интервалом "d 1" и сегодняшней датой
		// оставляем сегодняшнюю дату как дату первого выполнения
		if repeat == "d 1" || isToday(t, now) {
			// Оставляем сегодняшнюю дату
		} else if !afterNow(t, now) {
			// Для других случаев, если дата в прошлом, рассчитываем следующую
			next, err := NextDate(now, task.Date, repeat)
			if err != nil {
				return fmt.Errorf("ошибка расчета следующей даты: %v", err)
			}
			task.Date = next
		}
	} else {
		// правила нет — если дата в прошлом/сегодня, ставим сегодня
		if !afterNow(t, now) {
			task.Date = now.Format(DateFormat)
		}
	}

	return nil
}

// Вспомогательная функция для преобразования дня недели
func convertWeekdayToTargetFormat(wd time.Weekday) int {
	// Go: Sunday=0, Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6
	// Нам нужно: Monday=1, Tuesday=2, Wednesday=3, Thursday=4, Friday=5, Saturday=6, Sunday=7

	if wd == time.Sunday {
		return 7
	}
	return int(wd) // Monday=1, Tuesday=2, etc.
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("не удалось распарсить dstart: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат repeat")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("для правила d нужно указать интервал в днях")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("недопустимый интервал для d")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if !afterNow(date, now) {
			date = now
		}

		if len(parts) != 2 {
			return "", errors.New("для правила w нужно указать дни недели через запятую")
		}

		daysStr := strings.Split(parts[1], ",")
		targetDays := make(map[int]bool)
		for _, s := range daysStr {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("день недели должен быть от 1 до 7: %s", s)
			}
			targetDays[n] = true
		}

		current := date
		// Ищем в течение 400 дней
		for day := 0; day < 400; day++ {
			current = current.AddDate(0, 0, 1)

			// Правильное преобразование дня недели
			weekday := convertWeekdayToTargetFormat(current.Weekday())

			if targetDays[weekday] && afterNow(current, now) {
				return current.Format(DateFormat), nil
			}
		}
		return "", errors.New("не найдена подходящая дата")

	case "m":
		if len(parts) < 2 {
			return "", errors.New("для правила m нужно указать дни месяца")
		}

		dayStrs := strings.Split(parts[1], ",")
		targetDays := make(map[int]bool)
		for _, s := range dayStrs {
			n, err := strconv.Atoi(s)
			if err != nil || n < -2 || n > 31 || n == 0 {
				return "", fmt.Errorf("день месяца должен быть от 1 до 31 или -1, -2: %s", s)
			}
			targetDays[n] = true
		}

		targetMonths := make(map[int]bool)
		if len(parts) >= 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, s := range monthStrs {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("месяц должен быть от 1 до 12: %s", s)
				}
				targetMonths[m] = true
			}
		} else {
			for m := 1; m <= 12; m++ {
				targetMonths[m] = true
			}
		}

		// Начинаем поиск с следующего дня от исходной даты
		current := date.AddDate(0, 0, 1)

		// Ищем в течение 730 дней (2 года)
		for day := 0; day < 730; day++ {
			currentMonth := int(current.Month())
			currentDay := current.Day()
			lastDay := time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, current.Location()).Day()

			// Проверяем соответствие дню
			dayMatch := targetDays[currentDay] ||
				(targetDays[-1] && currentDay == lastDay) ||
				(targetDays[-2] && currentDay == lastDay-1)

			// Проверяем соответствие месяцу
			monthMatch := targetMonths[currentMonth]

			if dayMatch && monthMatch && afterNow(current, now) {
				return current.Format(DateFormat), nil
			}

			current = current.AddDate(0, 0, 1)
		}
		return "", errors.New("не найдена подходящая дата")

	default:
		return "", fmt.Errorf("неподдерживаемый формат repeat: %s", repeat)
	}

	return date.Format(DateFormat), nil
}

func checkSearchDate(date string) bool {
	if len(date) != 10 {
		return false
	}
	if date[2] != '.' || date[5] != '.' {
		return false
	}

	parts := strings.Split(date, ".")
	if len(parts) != 3 {
		return false
	}

	// Проверяем что все части - валидные числа
	day, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	_, err3 := strconv.Atoi(parts[2])

	// Дополнительная проверка на валидность даты
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// Проверяем что день и месяц в допустимых диапазонах
	if day < 1 || day > 31 || month < 1 || month > 12 {
		return false
	}

	return true
}

// convertSearchDateToDBFormat преобразует дату из DD.MM.YYYY в YYYYMMDD
func convertSearchDateToDBFormat(date string) (string, error) {
	parts := strings.Split(date, ".")
	if len(parts) != 3 {
		return "", errors.New("неверный формат даты")
	}

	day := parts[0]
	month := parts[1]
	year := parts[2]

	// Добавляем ведущие нули если нужно
	if len(day) == 1 {
		day = "0" + day
	}
	if len(month) == 1 {
		month = "0" + month
	}

	return year + month + day, nil
}

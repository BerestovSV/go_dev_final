package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`  // формат 20060102 (обязательное поле)
	Title   string `json:"title"` // обязательное поле
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"` // правило повторения (может быть пустым)
}

func AddTask(task *Task) (int64, error) {
	var id int64
	if DB == nil {
		return 0, errors.New("database is not initialized")
	}
	const query = `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}

	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date >= strftime('%Y%m%d', 'now')
		ORDER BY date ASC, id ASC 
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var id int64
		var task Task
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		// Конвертируем int64 в string
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}

	// Конвертируем string ID в int64 для базы данных
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.New("некорректный идентификатор задачи")
	}

	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = ?
	`

	var task Task
	var dbID int64
	err = DB.QueryRow(query, taskID).Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}

	// Конвертируем обратно в string для JSON
	task.ID = strconv.FormatInt(dbID, 10)
	return &task, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return errors.New("database is not initialized")
	}

	// Конвертируем string ID в int64 для базы данных
	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("некорректный идентификатор задачи")
	}

	const query = `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return err
	}

	// Проверяем, что запись была обновлена
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	if DB == nil {
		return errors.New("database is not initialized")
	}

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("некорректный идентификатор задачи")
	}

	const query = `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func UpdateDate(id string, newDate string) error {
	if DB == nil {
		return errors.New("database is not initialized")
	}

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("некорректный идентификатор задачи")
	}

	const query = `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, newDate, taskID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func SearchTasksByText(search string, limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}

	searchPattern := "%" + search + "%"
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE (title LIKE ? OR comment LIKE ?)
		AND date >= strftime('%Y%m%d', 'now')
		ORDER BY date ASC, id ASC 
		LIMIT ?
	`

	rows, err := DB.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasksFromRows(rows)
}

// SearchTasksByDate ищет задачи по конкретной дате
func SearchTasksByDate(date string, limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("database is not initialized")
	}

	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = ?
		ORDER BY date ASC, id ASC 
		LIMIT ?
	`

	rows, err := DB.Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasksFromRows(rows)
}

// scanTasksFromRows - вспомогательная функция для сканирования задач из rows
func scanTasksFromRows(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		var id int64
		var task Task
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

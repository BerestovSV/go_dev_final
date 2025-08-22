package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT
);
CREATE INDEX task_date ON scheduler(date);
`

type Database struct {
	db *sql.DB
}

// Интерфейс для работы с задачами
type TaskStore interface {
	AddTask(task *Task) (int64, error)
	GetTask(id string) (*Task, error)
	UpdateTask(task *Task) error
	DeleteTask(id string) error
	Tasks(limit int) ([]*Task, error)
	SearchTasksByText(search string, limit int) ([]*Task, error)
	SearchTasksByDate(date string, limit int) ([]*Task, error)
	UpdateDate(id string, newDate string) error
}

func NewDatabase(dbFile string) (*Database, error) {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return nil, err
		}
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

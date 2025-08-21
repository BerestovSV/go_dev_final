package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// глобальная переменная для соединения (оставляем как есть)
var DB *sql.DB

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

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return err
		}
	}
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			db.Close()
			return err
		}
	}

	DB = db // используем нашу глобальную переменную
	return nil
}

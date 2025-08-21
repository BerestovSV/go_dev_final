package main

import (
	"log"
	"net/http"
	"os"
	"todo-server/pkg/api"
	"todo-server/pkg/db"
)

func main() {
	// Получаем значение переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")

	// Если переменная пустая, используем путь по умолчанию
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы %v", err)
	}

	log.Println("База данных готова к работе")

	api.Init()

	webDir := "./web"

	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)

	port := ":7540"
	log.Printf("Сервер запущен на http://localhost%v", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

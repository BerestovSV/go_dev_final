package main

import (
	"log"
	"net/http"
	"todo-server/pkg/api"
	"todo-server/pkg/config"
	"todo-server/pkg/db"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем подключение к БД
	database, err := db.NewDatabase(cfg.DBFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы: %v", err)
	}
	defer database.Close()

	log.Println("База данных готова к работе")

	// Создаем API с конфигом и БД
	api := api.NewAPI(database, cfg)
	router := api.Init()

	// Статический контент
	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	router.Handle("/", fileServer)

	log.Printf("Сервер запущен на http://localhost%v", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}

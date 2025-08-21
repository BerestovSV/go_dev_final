package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
	http.HandleFunc("/api/signin", signinHandler)

	// http.HandleFunc("/api/nextdate", nextDayHandler)
	// http.HandleFunc("/api/task", taskHandler)
	// http.HandleFunc("/api/tasks", tasksHandler)
	// http.HandleFunc("/api/task/done", doneTaskHandler)
}

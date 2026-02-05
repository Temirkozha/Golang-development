package main

import (
	"log"
	"net/http"

	"task2/internal/handlers"
	"task2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	finalHandler := middleware.RequestIDMiddleware(
		middleware.AuthMiddleware(handlers.TasksHandler),
	)

	mux.HandleFunc("/tasks", finalHandler)

	log.Println("Server is running on :8080...")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
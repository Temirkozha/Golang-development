package main

import (
	"database/sql"
	"log"
	"net/http"

	"assignment3/internal/app"
	"assignment3/internal/repository/postgres/users"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "user=postgres password=2765 dbname=go-test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("База данных не отвечает:", err)
	}
	log.Println("Успешное подключение к PostgreSQL!")

	userRepo := users.NewUserRepository(db)
	userUsecase := app.NewUserUsecase(userRepo)
	userHandler := app.NewUserHandler(userUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", userHandler.HealthCheck)

	mux.Handle("/users", userHandler)
	mux.Handle("/users/", userHandler)

	protectedMux := app.AuthAndLogMiddleware(mux)

	log.Println("Сервер запущен на порту :8080...")
	if err := http.ListenAndServe(":8080", protectedMux); err != nil {
		log.Fatal(err)
	}
}

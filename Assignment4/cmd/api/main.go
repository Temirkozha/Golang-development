package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"assignment3/internal/app"
	"assignment3/internal/repository/postgres/users"

	_ "github.com/lib/pq"
)

func main() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	
	if host == "" {
		host = "localhost"
		port = "5432"
		user = "postgres"
		password = "secret"
		dbname = "go_test"
		sslmode = "disable"
	}

	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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

package app

import (
	"log"
	"os"
	"time"
	"practice-7/internal/controller/http/v1"
	"practice-7/internal/entity" // Обязательно для миграции
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Run() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "host=localhost user=user password=password dbname=practice7 port=5432 sslmode=disable"
	}

	var pg *postgres.Postgres
	var err error

	// Ждем, пока база в Docker станет доступной
	for i := 0; i < 10; i++ {
		pg, err = postgres.New(dbUrl)
		if err == nil {
			break
		}
		log.Printf("Connecting to DB... attempt %d/10", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Final connection error: %v", err)
	}

	// САМЫЙ ВАЖНЫЙ МОМЕНТ: Go сам создаст таблицу users, если её нет
	err = pg.Conn.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	handler := gin.New()
	v1.NewRouter(handler, userUseCase)

	log.Println("Server is ready on :8090")
	handler.Run(":8090")
}
package main

import (
	"database/sql"
	"log"
	"net/http"

	"assignment5/internal/handler"
	"assignment5/internal/repository"

	_ "github.com/lib/pq"
)

func main() {

	
	db, err := sql.Open("postgres", "postgres://postgres:password@db:5432/practice5?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	repo := repository.NewUserRepository(db)
	h := handler.NewUserHandler(repo)

	
	mux := http.NewServeMux()

	
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetUsers(w, r)
		}
	})

	
	mux.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.AddFriend(w, r)
		}
	})

	
	mux.HandleFunc("/common_friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetCommonFriends(w, r)
		}
	})
 
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
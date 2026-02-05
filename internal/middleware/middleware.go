package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		

		if r.Header.Get("X-API-KEY") != "secret12345" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}
		next(w, r)
	}
}

func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		reqID := rand.Intn(100000)
		idStr := strconv.Itoa(reqID)

		log.Printf("[RequestID: %s] %s %s", idStr, r.Method, r.URL.Path)
		w.Header().Set("X-Request-ID", idStr)

		next(w, r)
	}
}
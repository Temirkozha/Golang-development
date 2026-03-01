package app

import (
	"log"
	"net/http"
	"time"
)


func AuthAndLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != "secret" { 
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		
		start := time.Now()
		
		
		next.ServeHTTP(w, r)
		
		
		log.Printf("[Time: %s] | Method: %s | Endpoint: %s", start.Format(time.RFC3339), r.Method, r.URL.Path)
	})
}
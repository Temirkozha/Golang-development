package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"practice-9/repository"
	"time"
)

type ResponseData struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func IdempotencyMiddleware(store repository.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			val, err := store.Get(ctx, key)
			if err == nil {
				if val == "processing" {
					http.Error(w, "Duplicate request in progress", http.StatusConflict)
					return
				}
				var respData ResponseData
				json.Unmarshal([]byte(val), &respData)
				w.WriteHeader(respData.StatusCode)
				w.Write([]byte(respData.Body))
				return
			}

			ok, err := store.SetProcessing(ctx, key, 1*time.Minute)
			if err != nil || !ok {
				http.Error(w, "Conflict", http.StatusConflict)
				return
			}

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			respData := ResponseData{
				StatusCode: rec.Code,
				Body:       rec.Body.String(),
			}
			jsonData, _ := json.Marshal(respData)
			store.SetCompleted(ctx, key, string(jsonData), 24*time.Hour)

			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
		})
	}
}
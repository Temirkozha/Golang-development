package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"practice-9/repository"
	"practice-9/service"
	"sync"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func CalculateBackoff(attempt int, baseDelay time.Duration) time.Duration {
	backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
	return time.Duration(rand.Int63n(int64(backoff)))
}

func ExecutePayment(ctx context.Context, url string, key string) error {
	client := &http.Client{}
	maxRetries := 5
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
		req.Header.Set("Idempotency-Key", key)
		resp, err := client.Do(req)

		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Attempt %d: Success!\n", attempt+1)
			return nil
		}

		if !IsRetryable(resp, err) || attempt == maxRetries-1 {
			return fmt.Errorf("failed after all attempts")
		}

		wait := CalculateBackoff(attempt, baseDelay)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, wait)

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func main() {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	fmt.Println("--- Task 1: Resilient Client Simulation ---")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ExecutePayment(ctx, server.URL, "task1-key")

	fmt.Println("\n--- Task 2: Double-Click Attack Simulation ---")
	store := repository.NewRedisStore("localhost:6379")
	
	businessLogic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing started...")
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "paid", "amount": 1000}`))
	})

	mw := service.IdempotencyMiddleware(store)
	ts := httptest.NewServer(mw(businessLogic))
	defer ts.Close()

	var wg sync.WaitGroup
	idempotencyKey := "unique-payment-123"
	
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req, _ := http.NewRequest("POST", ts.URL, nil)
			req.Header.Set("Idempotency-Key", idempotencyKey)
			resp, _ := http.DefaultClient.Do(req)
			fmt.Printf("Request %d: Status %d\n", id, resp.StatusCode)
		}(i)
	}
	wg.Wait()
}
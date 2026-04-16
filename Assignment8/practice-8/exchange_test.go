package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"base":"USD","target":"KZT","rate":450.5}`)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "KZT")

	assert.NoError(t, err)
	assert.Equal(t, 450.5, rate)
}

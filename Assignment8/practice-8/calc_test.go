package main

import (
	"testing"
)


func TestDivide(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := Divide(10, 2)
		want := 5
		if err != nil {
			t.Errorf("Divide(10, 2) returned error: %v", err)
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("expected error for division by zero, but got nil")
		}
	})
}


func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"Both positive numbers", 10, 5, 5},     
		{"Positive minus zero", 5, 0, 5},        
		{"Negative minus positive", -2, 3, -5},  
		{"Both negative", -5, -5, 0},            
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
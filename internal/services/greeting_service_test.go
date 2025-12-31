package services

import (
	"context"
	"testing"
)

func TestGreetingService_Greet(t *testing.T) {
	t.Parallel()

	service := GreetingService{}

	got, err := service.Greet(context.Background(), "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello, Alice"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

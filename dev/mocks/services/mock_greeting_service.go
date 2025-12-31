package services

import (
	"context"
	"fmt"
)

type MockGreetingService struct{}

func (s MockGreetingService) Greet(_ context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s (mock)", name), nil
}

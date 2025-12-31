package services

import (
	"context"
	"fmt"
)

type GreetingService struct{}

func (s GreetingService) Greet(_ context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s", name), nil
}

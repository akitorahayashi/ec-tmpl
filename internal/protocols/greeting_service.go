package protocols

import "context"

type GreetingService interface {
	Greet(ctx context.Context, name string) (string, error)
}

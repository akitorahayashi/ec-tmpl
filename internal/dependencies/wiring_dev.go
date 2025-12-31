//go:build dev

package dependencies

import (
	mockservices "example.com/ec-tmpl/dev/mocks/services"
	"example.com/ec-tmpl/internal/config"
	"example.com/ec-tmpl/internal/protocols"
	"example.com/ec-tmpl/internal/services"
)

func BuildGreetingService(settings config.Settings) (protocols.GreetingService, error) {
	if settings.UseMockGreeting {
		return mockservices.MockGreetingService{}, nil
	}
	return services.GreetingService{}, nil
}

//go:build !dev

package dependencies

import (
	"fmt"

	"example.com/ec-tmpl/internal/config"
	"example.com/ec-tmpl/internal/protocols"
	"example.com/ec-tmpl/internal/services"
)

func BuildGreetingService(settings config.Settings) (protocols.GreetingService, error) {
	if settings.UseMockGreeting {
		return nil, fmt.Errorf("EC_TMPL_USE_MOCK_GREETING is supported only in dev builds")
	}
	return services.GreetingService{}, nil
}

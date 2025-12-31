package config

import (
	"fmt"
	"os"
	"strconv"
)

type Settings struct {
	AppName          string
	BindIP           string
	BindPort         int
	UseMockGreeting  bool
	ShutdownGraceSec int
}

func Load() (Settings, error) {
	appName := envString("EC_TMPL_APP_NAME", "ec-tmpl")
	bindIP := envString("EC_TMPL_BIND_IP", "0.0.0.0")

	bindPort, err := envInt("EC_TMPL_BIND_PORT", 8000)
	if err != nil {
		return Settings{}, fmt.Errorf("invalid EC_TMPL_BIND_PORT: %w", err)
	}

	useMockGreeting, err := envBool("EC_TMPL_USE_MOCK_GREETING", false)
	if err != nil {
		return Settings{}, fmt.Errorf("invalid EC_TMPL_USE_MOCK_GREETING: %w", err)
	}

	shutdownGraceSec, err := envInt("EC_TMPL_SHUTDOWN_GRACE_SEC", 10)
	if err != nil {
		return Settings{}, fmt.Errorf("invalid EC_TMPL_SHUTDOWN_GRACE_SEC: %w", err)
	}

	return Settings{
		AppName:          appName,
		BindIP:           bindIP,
		BindPort:         bindPort,
		UseMockGreeting:  useMockGreeting,
		ShutdownGraceSec: shutdownGraceSec,
	}, nil
}

func envString(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return defaultValue
	}
	return value
}

func envInt(key string, defaultValue int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func envBool(key string, defaultValue bool) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return parsed, nil
}

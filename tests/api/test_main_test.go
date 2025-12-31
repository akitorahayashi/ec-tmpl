//go:build api

package api_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"example.com/ec-tmpl/tests/testsupport"
)

var (
	apiOnce   sync.Once
	apiServer *testsupport.RunningAPIServer
	apiErr    error
	apiURL    string
)

func apiBaseURL(t *testing.T) string {
	t.Helper()

	baseURL := strings.TrimRight(os.Getenv("EC_TMPL_TEST_BASE_URL"), "/")
	if baseURL != "" {
		return baseURL
	}

	image := os.Getenv("EC_TMPL_E2E_IMAGE")
	if image == "" {
		t.Skip("EC_TMPL_E2E_IMAGE is not set")
	}

	apiOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		apiServer, apiErr = testsupport.StartAPIServer(ctx, testsupport.APIServerOptions{
			Image: image,
			Env: map[string]string{
				"EC_TMPL_USE_MOCK_GREETING": "true",
			},
			StartupTimeout: 60 * time.Second,
		})
		if apiErr == nil {
			apiURL = apiServer.BaseURL
			fmt.Printf("\nAPI tests running against: %s\n", apiURL)
		}
	})

	if apiErr != nil {
		t.Fatalf("failed to start api container: %v", apiErr)
	}
	return apiURL
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if apiServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := apiServer.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to terminate api container: %v\n", err)
		}
	}

	os.Exit(exitCode)
}

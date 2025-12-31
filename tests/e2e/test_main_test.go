//go:build e2e

package e2e_test

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
	e2eOnce   sync.Once
	e2eServer *testsupport.RunningAPIServer
	e2eErr    error
	e2eURL    string
)

func e2eBaseURL(t *testing.T) string {
	t.Helper()

	baseURL := strings.TrimRight(os.Getenv("EC_TMPL_TEST_BASE_URL"), "/")
	if baseURL != "" {
		return baseURL
	}

	image := os.Getenv("EC_TMPL_E2E_IMAGE")
	if image == "" {
		t.Skip("EC_TMPL_E2E_IMAGE is not set")
	}

	e2eOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		e2eServer, e2eErr = testsupport.StartAPIServer(ctx, testsupport.APIServerOptions{
			Image: image,
			Env: map[string]string{
				"EC_TMPL_USE_MOCK_GREETING": "false",
			},
			StartupTimeout: 60 * time.Second,
		})
		if e2eErr == nil {
			e2eURL = e2eServer.BaseURL
			fmt.Printf("\nE2E tests running against: %s\n", e2eURL)
		}
	})

	if e2eErr != nil {
		t.Fatalf("failed to start e2e container: %v", e2eErr)
	}
	return e2eURL
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if e2eServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := e2eServer.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to terminate e2e container: %v\n", err)
		}
	}

	os.Exit(exitCode)
}

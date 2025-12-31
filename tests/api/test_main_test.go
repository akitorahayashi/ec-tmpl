//go:build api

package api_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"example.com/ec-tmpl/tests/testsupport"
)

var apiManager testsupport.TestServerManager

func apiBaseURL(t *testing.T) string {
	t.Helper()
	return apiManager.GetBaseURL(t, testsupport.APIServerOptions{
		Env: map[string]string{
			"EC_TMPL_USE_MOCK_GREETING": "true",
		},
		StartupTimeout: 60 * time.Second,
	}, "API")
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := apiManager.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to terminate api container: %v\n", err)
	}

	os.Exit(exitCode)
}

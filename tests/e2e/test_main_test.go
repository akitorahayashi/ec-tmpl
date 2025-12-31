//go:build e2e

package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"example.com/ec-tmpl/tests/testsupport"
)

var e2eManager testsupport.TestServerManager

func e2eBaseURL(t *testing.T) string {
	t.Helper()
	return e2eManager.GetBaseURL(t, testsupport.APIServerOptions{
		Env: map[string]string{
			"EC_TMPL_USE_MOCK_GREETING": "false",
		},
		StartupTimeout: 60 * time.Second,
	}, "E2E")
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e2eManager.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to terminate e2e container: %v\n", err)
	}

	os.Exit(exitCode)
}

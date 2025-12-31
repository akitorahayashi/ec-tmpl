//go:build api

package api_test

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHealth_API(t *testing.T) {
	t.Parallel()

	baseURL := testBaseURL(t)
	waitForHealthy(t, baseURL, 15*time.Second)

	resp := getJSON(t, baseURL+"/health")
	if resp["status"] != "ok" {
		t.Fatalf("unexpected payload: %v", resp)
	}
}

func TestHello_API(t *testing.T) {
	t.Parallel()

	baseURL := testBaseURL(t)
	waitForHealthy(t, baseURL, 15*time.Second)

	resp := getJSON(t, baseURL+"/hello/Alice")
	if resp["message"] != "Hello, Alice" && resp["message"] != "Hello, Alice (mock)" {
		t.Fatalf("unexpected payload: %v", resp)
	}
}

func testBaseURL(t *testing.T) string {
	t.Helper()

	baseURL := strings.TrimRight(os.Getenv("EC_TMPL_TEST_BASE_URL"), "/")
	if baseURL == "" {
		t.Skip("EC_TMPL_TEST_BASE_URL is not set")
	}
	return baseURL
}

func waitForHealthy(t *testing.T, baseURL string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/health", nil)
		if err != nil {
			t.Fatalf("request error: %v", err)
		}

		client := &http.Client{Timeout: 1 * time.Second}
		res, err := client.Do(req)
		if err == nil && res.Body != nil {
			_ = res.Body.Close()
		}
		if err == nil && res.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}

	t.Fatalf("health did not become ready within %s", timeout)
}

func getJSON(t *testing.T, url string) map[string]string {
	t.Helper()

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", res.StatusCode)
	}

	var got map[string]string
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return got
}

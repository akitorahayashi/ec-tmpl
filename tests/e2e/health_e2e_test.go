//go:build e2e

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestHealth_E2E(t *testing.T) {
	t.Parallel()

	baseURL := e2eBaseURL(t)
	waitForHealthy(t, baseURL, 15*time.Second)

	resp := getJSON(t, baseURL+"/health")
	if resp["status"] != "ok" {
		t.Fatalf("unexpected payload: %v", resp)
	}
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

//go:build intg

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/ec-tmpl/internal/services"
)

func TestHealth_Intg(t *testing.T) {
	t.Parallel()

	e := NewHTTPServer(Dependencies{GreetingService: services.GreetingService{}})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d", rec.Code)
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if got["status"] != "ok" {
		t.Fatalf("unexpected payload: %v", got)
	}
}

func TestHello_Intg(t *testing.T) {
	t.Parallel()

	e := NewHTTPServer(Dependencies{GreetingService: services.GreetingService{}})

	req := httptest.NewRequest(http.MethodGet, "/hello/Bob", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d", rec.Code)
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if got["message"] != "Hello, Bob" {
		t.Fatalf("unexpected payload: %v", got)
	}
}

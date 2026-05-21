package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/edsuarezs/devsecops-pipeline-demo/config"
	"github.com/edsuarezs/devsecops-pipeline-demo/handlers"

)

func TestLiveness(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handlers.Liveness(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	assertContentType(t, rec)

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "alive" {
		t.Errorf("expected status 'alive', got '%s'", body["status"])
	}
}

func TestReadiness(t *testing.T) {
	cfg := &config.Config{
		AppVersion:  "0.1.0",
		Environment: "test",
	}

	handler := handlers.Readiness(cfg)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	assertContentType(t, rec)

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ready" {
		t.Errorf("expected status 'ready', got '%s'", body["status"])
	}
	if body["version"] != "0.1.0" {
		t.Errorf("expected version '0.1.0', got '%s'", body["version"])
	}
	if body["environment"] != "test" {
		t.Errorf("expected environment 'test', got '%s'", body["environment"])
	}
}

// assertContentType verifies the response has application/json Content-Type.
func assertContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}

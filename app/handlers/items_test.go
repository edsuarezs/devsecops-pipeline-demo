package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/edsuarezs/devsecops-pipeline-demo/handlers"
	"github.com/edsuarezs/devsecops-pipeline-demo/models"

)

// testRouter returns a Chi router with the items routes mounted.
func testRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/v1/items/", handlers.CreateItem)
	r.Get("/api/v1/items/", handlers.ListItems)
	r.Get("/api/v1/items/{id}", handlers.GetItem)
	r.Delete("/api/v1/items/{id}", handlers.DeleteItem)
	return r
}

// postItem is a helper that creates an item with correct headers.
func postItem(r *chi.Mux, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/items/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

// ── CreateItem tests ────────────────────────────────────────

func TestCreateItemSuccess(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	rec := postItem(r, `{"name":"Widget","description":"A useful widget","price":9.99}`)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}

	assertContentType(t, rec)

	var item models.Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if item.ID != 1 {
		t.Errorf("expected ID 1, got %d", item.ID)
	}
	if item.Name != "Widget" {
		t.Errorf("expected name 'Widget', got '%s'", item.Name)
	}
	if item.Price != 9.99 {
		t.Errorf("expected price 9.99, got %f", item.Price)
	}
}

func TestCreateItemMissingName(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	rec := postItem(r, `{"price":9.99}`)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}
}

func TestCreateItemNegativePrice(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	rec := postItem(r, `{"name":"Widget","price":-5.0}`)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}
}

func TestCreateItemInvalidJSON(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	rec := postItem(r, `{not valid json}`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCreateItemWrongContentType(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/items/", strings.NewReader(`{"name":"Widget","price":1.0}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415, got %d", rec.Code)
	}
}

func TestCreateItemNameTooLong(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	longName := strings.Repeat("a", 101)
	rec := postItem(r, `{"name":"`+longName+`","price":1.0}`)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}
}

// ── ListItems tests ─────────────────────────────────────────

func TestListItemsEmpty(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	assertContentType(t, rec)

	var items []models.Item
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestListItemsAfterCreate(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	postItem(r, `{"name":"A","price":1.0}`)
	postItem(r, `{"name":"B","price":2.0}`)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	var items []models.Item
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

// ── GetItem tests ───────────────────────────────────────────

func TestGetItemExisting(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	postItem(r, `{"name":"Widget","price":9.99}`)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	assertContentType(t, rec)

	var item models.Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if item.Name != "Widget" {
		t.Errorf("expected name 'Widget', got '%s'", item.Name)
	}
}

func TestGetItemNotFound(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestGetItemInvalidID(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// ── DeleteItem tests ────────────────────────────────────────

func TestDeleteItemExisting(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	postItem(r, `{"name":"Widget","price":9.99}`)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}

	// Verify it's gone
	req = httptest.NewRequest(http.MethodGet, "/api/v1/items/1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 after delete, got %d", rec.Code)
	}
}

func TestDeleteItemNotFound(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestDeleteItemInvalidID(t *testing.T) {
	handlers.ResetStore()
	r := testRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

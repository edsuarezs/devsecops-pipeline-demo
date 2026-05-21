package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"github.com/go-chi/chi/v5"
	"github.com/edsuarezs/devsecops-pipeline-demo/models"

)

// maxBodySize is the maximum allowed request body size (1MB).
// Prevents DoS attacks via large payloads.
const maxBodySize = 1 << 20

// ── In-memory store ─────────────────────────────────────────

var (
	store   = make(map[int]models.Item)
	nextID  = 1
	storeMu sync.RWMutex
)

// ResetStore clears the in-memory store. Used by tests for isolation.
func ResetStore() {
	storeMu.Lock()
	defer storeMu.Unlock()
	store = make(map[int]models.Item)
	nextID = 1
}

// ── Handlers ────────────────────────────────────────────────

// CreateItem handles POST /api/v1/items/
func CreateItem(w http.ResponseWriter, r *http.Request) {
	// Validate Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{
			"error": "Content-Type must be application/json",
		})
		return
	}

	// Limit body size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var input models.ItemCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if err := input.Validate(); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
		return
	}

	storeMu.Lock()
	item := models.Item{
		ID:          nextID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}
	store[nextID] = item
	nextID++
	storeMu.Unlock()

	writeJSON(w, http.StatusCreated, item)
}

// ListItems handles GET /api/v1/items/
func ListItems(w http.ResponseWriter, r *http.Request) {
	storeMu.RLock()
	items := make([]models.Item, 0, len(store))
	for _, item := range store {
		items = append(items, item)
	}
	storeMu.RUnlock()

	writeJSON(w, http.StatusOK, items)
}

// GetItem handles GET /api/v1/items/{id}
func GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid item ID",
		})
		return
	}

	storeMu.RLock()
	item, exists := store[id]
	storeMu.RUnlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// DeleteItem handles DELETE /api/v1/items/{id}
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid item ID",
		})
		return
	}

	storeMu.Lock()
	_, exists := store[id]
	if !exists {
		storeMu.Unlock()
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
		return
	}
	delete(store, id)
	storeMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

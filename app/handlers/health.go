package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/edsuarezs/devsecops-pipeline-demo/config"

)

// Liveness is the Kubernetes liveness probe.
// Returns 200 if the process is alive — nothing more.
// K8s restarts the pod if this fails.
func Liveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	})
}

// Readiness is the Kubernetes readiness probe.
// Returns 200 when the app is ready to serve traffic.
// K8s removes the pod from the Service if this fails (no restart).
func Readiness(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		result := "ready"

		// Future: check real dependencies before reporting ready
		// if err := db.Ping(r.Context()); err != nil {
		//     status = http.StatusServiceUnavailable
		//     result = "not_ready"
		// }

		writeJSON(w, status, map[string]string{
			"status":      result,
			"version":     cfg.AppVersion,
			"environment": cfg.Environment,
		})
	}
}

// writeJSON encodes data as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

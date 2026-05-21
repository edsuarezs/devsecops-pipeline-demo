package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/edsuarezs/devsecops-pipeline-demo/config"
	"github.com/edsuarezs/devsecops-pipeline-demo/handlers"
	"github.com/edsuarezs/devsecops-pipeline-demo/middleware"

)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel(),
	}))
	slog.SetDefault(logger)

	r := setupRouter(cfg)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	// ── Start server in background ──────────────────────────
	go func() {
		slog.Info("server starting",
			"port", cfg.Port,
			"environment", cfg.Environment,
			"version", cfg.AppVersion,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// ── Graceful shutdown ───────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}

func setupRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// ── Middleware stack (order matters) ─────────────────────
	r.Use(middleware.SecurityHeaders) // security headers first
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)

	// ── Kubernetes probes (outside /api to avoid middleware overhead) ──
	r.Get("/healthz", handlers.Liveness)
	r.Get("/readyz", handlers.Readiness(cfg))

	// ── API v1 ──────────────────────────────────────────────
	r.Route("/api/v1/items", func(r chi.Router) {
		r.Post("/", handlers.CreateItem)
		r.Get("/", handlers.ListItems)
		r.Get("/{id}", handlers.GetItem)
		r.Delete("/{id}", handlers.DeleteItem)
	})

	return r
}

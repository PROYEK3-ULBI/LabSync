package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/repository"
	"labsync-backend/internal/service"
)

func TestHealthHandler(t *testing.T) {
	app := fiber.New()

	// Inisialisasi service dengan repository tanpa DB (nil)
	repo := repository.NewHealthRepository(nil)
	svc := service.NewHealthService(repo)
	h := NewHealthHandler(svc)

	app.Get("/api/v1/health", h.Check)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body service.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json body: %v", err)
	}

	if body.Status != "healthy" {
		t.Errorf("expected status healthy, got %s", body.Status)
	}
	if body.Services["database"] != "connected" {
		t.Errorf("expected database connected, got %s", body.Services["database"])
	}
}

func TestHealthHandlerDegraded(t *testing.T) {
	app := fiber.New()

	// Inisialisasi DB dengan SQL nil sehingga Ping menghasilkan error
	mockDB := &repository.DB{SQL: nil}
	repo := repository.NewHealthRepository(mockDB)
	svc := service.NewHealthService(repo)
	h := NewHealthHandler(svc)

	app.Get("/api/v1/health", h.Check)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", resp.StatusCode)
	}

	var body service.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json body: %v", err)
	}

	if body.Status != "degraded" {
		t.Errorf("expected status degraded, got %s", body.Status)
	}
	if body.Services["database"] != "disconnected" {
		t.Errorf("expected database disconnected, got %s", body.Services["database"])
	}
}

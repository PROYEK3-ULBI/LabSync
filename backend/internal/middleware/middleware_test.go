package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/model"
)

func TestRequestIDMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(NewRequestID())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Test case 1: generate request ID jika tidak disediakan
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqID := resp.Header.Get(fiber.HeaderXRequestID)
	if reqID == "" {
		t.Errorf("expected X-Request-ID header to be set")
	}

	// Test case 2: teruskan request ID jika sudah ada
	customID := "custom-req-id-123"
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set(fiber.HeaderXRequestID, customID)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp2.Header.Get(fiber.HeaderXRequestID) != customID {
		t.Errorf("expected header %s, got %s", customID, resp2.Header.Get(fiber.HeaderXRequestID))
	}
}

func TestCentralErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: CentralErrorHandler,
	})
	app.Use(NewRequestID())

	app.Get("/bad-request", func(c *fiber.Ctx) error {
		return model.NewBadRequest("Input tidak sesuai")
	})

	app.Get("/fiber-error", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	app.Get("/unhandled", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "database connection string leaked here!")
	})

	// 1. AppError bad request
	req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var errBody model.ErrorResponse
	bodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(bodyBytes, &errBody)

	if errBody.Error.Code != "BAD_REQUEST" {
		t.Errorf("expected code BAD_REQUEST, got %s", errBody.Error.Code)
	}
	if errBody.Error.Message != "Input tidak sesuai" {
		t.Errorf("expected custom message, got %s", errBody.Error.Message)
	}
	if errBody.Error.RequestID == "" {
		t.Errorf("expected non-empty request_id")
	}

	// 2. Fiber error not found
	req2 := httptest.NewRequest(http.MethodGet, "/fiber-error", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp2.StatusCode)
	}

	var errBody2 model.ErrorResponse
	bodyBytes2, _ := io.ReadAll(resp2.Body)
	_ = json.Unmarshal(bodyBytes2, &errBody2)

	if errBody2.Error.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", errBody2.Error.Code)
	}
}

func TestRedactURLQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/api/v1/auth/callback?token=secret123&code=ok",
			expected: "/api/v1/auth/callback?code=ok&token=%5BREDACTED%5D",
		},
		{
			input:    "/api/v1/test?password=admin123&user=john",
			expected: "/api/v1/test?password=%5BREDACTED%5D&user=john",
		},
		{
			input:    "/api/v1/health?page=1",
			expected: "/api/v1/health?page=1",
		},
	}

	for _, tt := range tests {
		got := redactURLQuery(tt.input)
		if got != tt.expected {
			t.Errorf("redactURLQuery(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestRecoverMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: CentralErrorHandler,
	})
	app.Use(NewRecover())

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("kondisi darurat")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected test error: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 on recovered panic, got %d", resp.StatusCode)
	}
}

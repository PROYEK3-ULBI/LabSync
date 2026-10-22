package middleware

import (
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var sensitiveQueryKeys = map[string]bool{
	"token":      true,
	"password":   true,
	"secret":     true,
	"key":        true,
	"api_key":    true,
	"auth":       true,
	"credential": true,
}

// StructuredLogger mencatat HTTP request terstruktur dengan redaksi data sensitif
func StructuredLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			if handlerErr := c.App().ErrorHandler(c, err); handlerErr != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
			err = nil
		}
		latency := time.Since(start)

		requestID := c.GetRespHeader(fiber.HeaderXRequestID)
		if requestID == "" {
			if id, ok := c.Locals("requestid").(string); ok {
				requestID = id
			}
		}

		redactedPath := redactURLQuery(c.OriginalURL())
		status := c.Response().StatusCode()

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(c.Context(), level, "http_request",
			"method", c.Method(),
			"path", redactedPath,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"ip", c.IP(),
			"request_id", requestID,
		)

		return err
	}
}

// redactURLQuery menyamarkan nilai query params sensitif
func redactURLQuery(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	q := u.Query()
	modified := false
	for k := range q {
		if sensitiveQueryKeys[strings.ToLower(k)] {
			q.Set(k, "[REDACTED]")
			modified = true
		}
	}

	if modified {
		u.RawQuery = q.Encode()
		return u.String()
	}

	return rawURL
}

package middleware

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/model"
)

// CentralErrorHandler menangani seluruh error aplikasi dengan format JSON konsisten
func CentralErrorHandler(c *fiber.Ctx, err error) error {
	requestID := c.GetRespHeader(fiber.HeaderXRequestID)
	if requestID == "" {
		if id, ok := c.Locals("requestid").(string); ok {
			requestID = id
		}
	}
	if requestID == "" {
		requestID = c.Get(fiber.HeaderXRequestID)
	}

	statusCode := fiber.StatusInternalServerError
	code := "INTERNAL_SERVER_ERROR"
	message := "Terjadi kesalahan pada server"

	var appErr *model.AppError
	var fiberErr *fiber.Error

	if errors.As(err, &appErr) {
		statusCode = appErr.StatusCode
		code = appErr.Code
		message = appErr.Message
		if appErr.Internal != nil {
			slog.Error("app error terdeteksi", "code", code, "status", statusCode, "internal_error", appErr.Internal.Error(), "request_id", requestID)
		}
	} else if errors.As(err, &fiberErr) {
		statusCode = fiberErr.Code
		switch statusCode {
		case fiber.StatusBadRequest:
			code = "BAD_REQUEST"
			message = "Permintaan tidak valid"
		case fiber.StatusNotFound:
			code = "NOT_FOUND"
			message = "Rute atau sumber daya tidak ditemukan"
		case fiber.StatusMethodNotAllowed:
			code = "METHOD_NOT_ALLOWED"
			message = "Metode HTTP tidak diizinkan"
		case fiber.StatusRequestEntityTooLarge:
			code = "PAYLOAD_TOO_LARGE"
			message = "Ukuran berkas atau permintaan melebihi batas yang diizinkan"
		case fiber.StatusTooManyRequests:
			code = "TOO_MANY_REQUESTS"
			message = "Terlalu banyak permintaan, silakan coba lagi beberapa saat"
		default:
			code = "HTTP_ERROR"
			message = fiberErr.Message
		}
	} else {
		// Log error internal tanpa mengekspos ke klien
		slog.Error("unhandled server error", "error", err.Error(), "request_id", requestID)
	}

	return c.Status(statusCode).JSON(model.ErrorResponse{
		Error: model.ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	})
}

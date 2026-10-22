package handler

import (
	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// Check menangani endpoint GET /api/v1/health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	res, isHealthy := h.service.Check(c.Context())
	if !isHealthy {
		return c.Status(fiber.StatusServiceUnavailable).JSON(res)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List menangani GET /api/v1/audit-logs (Admin only, append-only inspection)
func (h *AuditHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	q := model.AuditLogListQuery{
		Page:       page,
		Limit:      limit,
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		UserID:     c.Query("user_id"),
	}

	res, err := h.auditService.List(c.Context(), q)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

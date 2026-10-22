package handler

import (
	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// Student menangani GET /api/v1/dashboard/student
func (h *DashboardHandler) Student(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	res, err := h.dashboardService.GetStudentDashboard(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Admin menangani GET /api/v1/dashboard/admin
func (h *DashboardHandler) Admin(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	res, err := h.dashboardService.GetAdminDashboard(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

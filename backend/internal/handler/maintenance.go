package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type MaintenanceHandler struct {
	maintService *service.MaintenanceService
}

func NewMaintenanceHandler(maintService *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{maintService: maintService}
}

// Create menangani POST /api/v1/maintenance (Admin only)
func (h *MaintenanceHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var req model.CreateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.maintService.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Update menangani PUT & PATCH /api/v1/maintenance/:id (Admin only)
func (h *MaintenanceHandler) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID pemeliharaan tidak boleh kosong")
	}

	var req model.UpdateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.maintService.Update(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID menangani GET /api/v1/maintenance/:id (Admin only)
func (h *MaintenanceHandler) GetByID(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID pemeliharaan tidak boleh kosong")
	}

	res, err := h.maintService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// List menangani GET /api/v1/maintenance (Admin only)
func (h *MaintenanceHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	q := model.MaintenanceListQuery{
		Page:     page,
		Limit:    limit,
		AssetID:  c.Query("asset_id"),
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
	}

	res, err := h.maintService.List(c.Context(), q)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

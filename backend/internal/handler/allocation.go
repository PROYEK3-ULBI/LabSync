package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type AllocationHandler struct {
	allocService *service.AllocationService
}

func NewAllocationHandler(allocService *service.AllocationService) *AllocationHandler {
	return &AllocationHandler{allocService: allocService}
}

// Create menangani POST /api/v1/allocations (Admin only)
func (h *AllocationHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var req model.CreateAllocationRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.allocService.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Return menangani POST /api/v1/allocations/:id/return (Admin only)
func (h *AllocationHandler) Return(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID alokasi tidak boleh kosong")
	}

	var req model.ReturnAllocationRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.allocService.Return(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// List menangani GET /api/v1/allocations (Admin only)
func (h *AllocationHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	q := model.AllocationListQuery{
		Page:      page,
		Limit:     limit,
		AssetID:   c.Query("asset_id"),
		UserID:    c.Query("user_id"),
		SessionID: c.Query("session_id"),
		Status:    c.Query("status"),
	}

	res, err := h.allocService.List(c.Context(), q)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// MyAllocations menangani GET /api/v1/allocations/my (Mahasiswa)
func (h *AllocationHandler) MyAllocations(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	q := model.AllocationListQuery{
		Page:      page,
		Limit:     limit,
		AssetID:   c.Query("asset_id"),
		SessionID: c.Query("session_id"),
		Status:    c.Query("status"),
	}

	res, err := h.allocService.GetMyAllocations(c.Context(), claims.UserID, q)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID menangani GET /api/v1/allocations/:id
func (h *AllocationHandler) GetByID(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID alokasi tidak boleh kosong")
	}

	res, err := h.allocService.GetByID(c.Context(), claims.UserID, claims.Role, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// AssetHistory menangani GET /api/v1/assets/:id/history
func (h *AllocationHandler) AssetHistory(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID aset tidak boleh kosong")
	}

	res, err := h.allocService.GetAssetHistory(c.Context(), claims.UserID, claims.Role, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

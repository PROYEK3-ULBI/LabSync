package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type AssetHandler struct {
	assetService *service.AssetService
}

func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

// List menangani GET /api/v1/assets
func (h *AssetHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	query := model.AssetListQuery{
		Page:     page,
		Limit:    limit,
		Status:   c.Query("status"),
		Category: c.Query("category"),
		Search:   c.Query("search"),
	}

	res, err := h.assetService.List(c.Context(), query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID menangani GET /api/v1/assets/:id
func (h *AssetHandler) GetByID(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	res, err := h.assetService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Create menangani POST /api/v1/assets (Admin only)
func (h *AssetHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var req model.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.assetService.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Update menangani PUT /api/v1/assets/:id (Admin only)
func (h *AssetHandler) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.assetService.Update(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateStatus menangani PATCH /api/v1/assets/:id/status (Admin only)
func (h *AssetHandler) UpdateStatus(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateAssetStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.assetService.UpdateStatus(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Delete menangani DELETE /api/v1/assets/:id (Admin only)
func (h *AssetHandler) Delete(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	if err := h.assetService.Delete(c.Context(), claims.UserID, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Aset berhasil dihapus",
	})
}

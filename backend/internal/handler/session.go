package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// List menangani GET /api/v1/sessions
func (h *SessionHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	date := c.Query("session_date")
	if date == "" {
		date = c.Query("date")
	}

	query := model.SessionListQuery{
		Page:        page,
		Limit:       limit,
		Status:      c.Query("status"),
		SessionDate: date,
		Search:      c.Query("search"),
	}

	res, err := h.sessionService.List(c.Context(), query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetActive menangani GET /api/v1/sessions/active
func (h *SessionHandler) GetActive(c *fiber.Ctx) error {
	res, err := h.sessionService.GetActive(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetByID menangani GET /api/v1/sessions/:id
func (h *SessionHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	res, err := h.sessionService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Create menangani POST /api/v1/sessions (Admin only)
func (h *SessionHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var req model.CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.sessionService.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Update menangani PUT /api/v1/sessions/:id (Admin only)
func (h *SessionHandler) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := c.Params("id")
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.sessionService.Update(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateStatus menangani PATCH /api/v1/sessions/:id/status (Admin only)
func (h *SessionHandler) UpdateStatus(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := c.Params("id")
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateSessionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	res, err := h.sessionService.UpdateStatus(c.Context(), claims.UserID, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Delete menangani DELETE /api/v1/sessions/:id (Admin only)
func (h *SessionHandler) Delete(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	id := c.Params("id")
	if id == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	if err := h.sessionService.Delete(c.Context(), claims.UserID, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sesi praktikum berhasil dihapus",
	})
}

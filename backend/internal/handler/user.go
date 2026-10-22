package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List menangani GET /api/v1/users
func (h *UserHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	query := model.UserListQuery{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
		Role:   c.Query("role"),
		Status: c.Query("status"),
	}

	res, err := h.userService.List(c.Context(), query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// Create menangani POST /api/v1/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	user, err := h.userService.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetByID menangani GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	targetID := c.Params("id")
	if targetID == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	user, err := h.userService.GetByID(c.Context(), claims.UserID, claims.Role, targetID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Update menangani PUT /api/v1/users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	targetID := c.Params("id")
	if targetID == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	user, err := h.userService.Update(c.Context(), claims.UserID, claims.Role, targetID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateStatus menangani PATCH /api/v1/users/:id/status
func (h *UserHandler) UpdateStatus(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	targetID := c.Params("id")
	if targetID == "" {
		return model.NewBadRequest("Parameter ID tidak boleh kosong")
	}

	var req model.UpdateUserStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	user, err := h.userService.UpdateStatus(c.Context(), claims.UserID, targetID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

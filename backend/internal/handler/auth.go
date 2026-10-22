package handler

import (
	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/config"
	"labsync-backend/internal/middleware"
	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Login menangani POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return model.NewBadRequest("Format JSON body tidak valid", err)
	}

	user, token, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return err
	}

	// Simpan JWT di dalam HttpOnly Cookie (tidak diekspos ke JavaScript)
	c.Cookie(&fiber.Cookie{
		Name:     middleware.TokenCookie,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: h.cfg.CookieSameSite,
		MaxAge:   h.cfg.JWTExpiryHours * 3600,
	})

	return c.Status(fiber.StatusOK).JSON(model.LoginResponse{
		Message: "Login berhasil",
		User:    *user,
	})
}

// Logout menangani POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	var userID string
	if claims != nil {
		userID = claims.UserID
	}

	_ = h.authService.Logout(c.Context(), userID)

	// Hapus cookie token
	c.Cookie(&fiber.Cookie{
		Name:     middleware.TokenCookie,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: h.cfg.CookieSameSite,
		MaxAge:   -1,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}

// Me menangani GET /api/v1/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		return model.NewUnauthorized("Autentikasi diperlukan")
	}

	profile, err := h.authService.Me(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

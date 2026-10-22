package middleware

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/model"
	"labsync-backend/internal/service"
)

const (
	UserContextKey = "authenticated_user"
	TokenCookie    = "token"
)

// Authenticate memvalidasi JWT dari HttpOnly cookie (atau header Authorization Bearer sebagai fallback)
func Authenticate(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies(TokenCookie)
		if tokenStr == "" {
			authHeader := c.Get(fiber.HeaderAuthorization)
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			return model.NewUnauthorized("Autentikasi diperlukan, silakan login terlebih dahulu")
		}

		claims, err := authService.ValidateToken(tokenStr)
		if err != nil {
			return model.NewUnauthorized("Sesi tidak valid atau telah kedaluwarsa")
		}

		c.Locals(UserContextKey, claims)
		return c.Next()
	}
}

// RequireRole membatasi akses endpoint berdasarkan role pengguna
func RequireRole(allowedRoles ...string) fiber.Handler {
	allowedMap := make(map[string]bool)
	for _, r := range allowedRoles {
		allowedMap[strings.ToUpper(r)] = true
	}

	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(UserContextKey).(*model.Claims)
		if !ok || claims == nil {
			return model.NewUnauthorized("Sesi pengguna tidak valid")
		}

		if !allowedMap[strings.ToUpper(claims.Role)] {
			return model.NewForbidden("Akses ditolak: peran tidak memiliki izin")
		}

		return c.Next()
	}
}

// GetUserClaims mengambil klaim pengguna dari konteks Fiber
func GetUserClaims(c *fiber.Ctx) *model.Claims {
	if claims, ok := c.Locals(UserContextKey).(*model.Claims); ok {
		return claims
	}
	return nil
}

// CSRFProtection memverifikasi Origin/Referer untuk request mutasi (POST, PUT, PATCH, DELETE) berotentikasi cookie
func CSRFProtection(allowOrigins string) fiber.Handler {
	allowedList := strings.Split(allowOrigins, ",")
	allowedMap := make(map[string]bool)
	for _, o := range allowedList {
		allowedMap[strings.TrimSpace(o)] = true
	}

	return func(c *fiber.Ctx) error {
		method := c.Method()
		if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
			return c.Next()
		}

		// Jika request membawa cookie token, pastikan Origin terpercaya
		if c.Cookies(TokenCookie) != "" {
			origin := c.Get(fiber.HeaderOrigin)
			if origin == "" {
				referer := c.Get("Referer")
				if referer != "" {
					if u, err := url.Parse(referer); err == nil {
						origin = u.Scheme + "://" + u.Host
					}
				}
			}

			if origin != "" && !allowedMap[origin] {
				return model.NewForbidden("Origin permintaan tidak diizinkan")
			}
		}

		return c.Next()
	}
}

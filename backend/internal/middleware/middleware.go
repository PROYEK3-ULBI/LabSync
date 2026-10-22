package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// NewRequestID menambahkan atau meneruskan X-Request-ID
func NewRequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: "requestid",
	})
}

// NewRecover menangani panic agar aplikasi tidak crash seketika
func NewRecover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: false,
	})
}

// NewCORS mengkonfigurasi CORS berdasarkan allowlist environment
func NewCORS(allowOrigins string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})
}

// NewRateLimiter menerapkan rate limiter sederhana berbasis IP
func NewRateLimiter(maxRequests int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.ErrTooManyRequests
		},
	})
}

package router

import (
	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/client"
	"labsync-backend/internal/config"
	"labsync-backend/internal/handler"
	"labsync-backend/internal/middleware"
	"labsync-backend/internal/repository"
	"labsync-backend/internal/service"
)

// Setup mengonfigurasi routing default menggunakan repositori PostgreSQL
func Setup(app *fiber.App, cfg *config.Config, db *repository.DB) {
	healthRepo := repository.NewHealthRepository(db)
	healthService := service.NewHealthService(healthRepo)

	userRepo := repository.NewPostgresUserRepository(db)
	auditRepo := repository.NewPostgresAuditRepository(db)
	sessionRepo := repository.NewPostgresSessionRepository(db)
	faceRepo := repository.NewPostgresFaceRepository(db)
	attendanceRepo := repository.NewPostgresAttendanceRepository(db)
	assetRepo := repository.NewPostgresAssetRepository(db)
	allocRepo := repository.NewPostgresAllocationRepository(db)
	maintRepo := repository.NewPostgresMaintenanceRepository(db)

	faceClient := client.NewHTTPFaceClient(cfg.FaceServiceURL, cfg.InternalAPIKey, cfg.FaceServiceTimeoutSec)

	authService := service.NewAuthService(userRepo, auditRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	userService := service.NewUserService(userRepo, auditRepo)
	sessionService := service.NewSessionService(sessionRepo, auditRepo)
	faceService := service.NewFaceService(userRepo, faceRepo, auditRepo, faceClient, cfg.BiometricEncryptionKey)
	attendanceService := service.NewAttendanceService(attendanceRepo, sessionRepo, faceRepo, userRepo, auditRepo, faceClient, cfg.BiometricEncryptionKey)
	assetService := service.NewAssetService(assetRepo, auditRepo)
	allocService := service.NewAllocationService(allocRepo, assetRepo, userRepo, sessionRepo, attendanceRepo, auditRepo)
	maintService := service.NewMaintenanceService(maintRepo, assetRepo, allocRepo, auditRepo)
	auditService := service.NewAuditService(auditRepo)
	dashService := service.NewDashboardService(userRepo, sessionRepo, attendanceRepo, assetRepo, allocRepo, maintRepo, auditRepo)
	reportService := service.NewReportService(attendanceRepo, assetRepo, allocRepo, maintRepo)

	SetupWithServices(app, cfg, authService, userService, sessionService, faceService, attendanceService, assetService, allocService, maintService, auditService, dashService, reportService, healthService)
}

// SetupWithServices mengonfigurasi seluruh routing dan middleware dengan dependency yang diinjeksi
func SetupWithServices(
	app *fiber.App,
	cfg *config.Config,
	authService *service.AuthService,
	userService *service.UserService,
	sessionService *service.SessionService,
	faceService *service.FaceService,
	attendanceService *service.AttendanceService,
	assetService *service.AssetService,
	allocService *service.AllocationService,
	maintService *service.MaintenanceService,
	auditService *service.AuditService,
	dashService *service.DashboardService,
	reportService *service.ReportService,
	healthService *service.HealthService,
) {
	// Middleware global
	app.Use(middleware.NewRequestID())
	app.Use(middleware.NewRecover())
	app.Use(middleware.NewCORS(cfg.CORSAllowOrigins))
	app.Use(middleware.StructuredLogger())
	app.Use(middleware.CSRFProtection(cfg.CORSAllowOrigins))

	// Handler
	healthHandler := handler.NewHealthHandler(healthService)
	docsHandler := handler.NewDocsHandler()
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	faceHandler := handler.NewFaceHandler(faceService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	assetHandler := handler.NewAssetHandler(assetService)
	allocHandler := handler.NewAllocationHandler(allocService)
	maintHandler := handler.NewMaintenanceHandler(maintService)
	auditHandler := handler.NewAuditHandler(auditService)
	dashHandler := handler.NewDashboardHandler(dashService)
	reportHandler := handler.NewReportHandler(reportService)

	// Rate limiter untuk endpoint login
	loginLimiter := middleware.NewRateLimiter(15)

	// Base API Group
	v1 := app.Group("/api/v1")

	// System & Docs
	v1.Get("/health", healthHandler.Check)
	v1.Get("/openapi.json", docsHandler.OpenAPISpec)
	v1.Get("/docs", docsHandler.SwaggerUI)

	// Auth Routes
	authGroup := v1.Group("/auth")
	authGroup.Post("/login", loginLimiter, authHandler.Login)
	authGroup.Post("/logout", middleware.Authenticate(authService), authHandler.Logout)
	authGroup.Get("/me", middleware.Authenticate(authService), authHandler.Me)

	// User Routes (Memerlukan autentikasi)
	userGroup := v1.Group("/users", middleware.Authenticate(authService))
	userGroup.Get("", middleware.RequireRole("ADMIN"), userHandler.List)
	userGroup.Post("", middleware.RequireRole("ADMIN"), userHandler.Create)
	userGroup.Get("/:id", userHandler.GetByID) // RBAC diperiksa di dalam service (Admin atau pemilik)
	userGroup.Put("/:id", userHandler.Update)  // RBAC diperiksa di dalam service (Admin atau pemilik)
	userGroup.Patch("/:id/status", middleware.RequireRole("ADMIN"), userHandler.UpdateStatus)

	// Session Routes (Memerlukan autentikasi)
	sessionGroup := v1.Group("/sessions", middleware.Authenticate(authService))
	sessionGroup.Get("/active", sessionHandler.GetActive)
	sessionGroup.Get("", sessionHandler.List)
	sessionGroup.Get("/:id", sessionHandler.GetByID)
	sessionGroup.Post("", middleware.RequireRole("ADMIN"), sessionHandler.Create)
	sessionGroup.Put("/:id", middleware.RequireRole("ADMIN"), sessionHandler.Update)
	sessionGroup.Patch("/:id/status", middleware.RequireRole("ADMIN"), sessionHandler.UpdateStatus)
	sessionGroup.Delete("/:id", middleware.RequireRole("ADMIN"), sessionHandler.Delete)

	// Face Biometrics Routes (Memerlukan autentikasi)
	faceGroup := v1.Group("/face", middleware.Authenticate(authService))
	faceGroup.Post("/enroll", middleware.RequireRole("STUDENT"), faceHandler.Enroll)
	faceGroup.Get("/status", faceHandler.GetStatus)
	faceGroup.Delete("/template/:user_id", middleware.RequireRole("ADMIN"), faceHandler.Revoke)

	// Attendance Routes (Memerlukan autentikasi)
	attGroup := v1.Group("/attendance", middleware.Authenticate(authService))
	attGroup.Post("/check-in", middleware.RequireRole("STUDENT"), attendanceHandler.CheckIn)
	attGroup.Post("/check-out", middleware.RequireRole("STUDENT"), attendanceHandler.CheckOut)
	attGroup.Get("/history", middleware.RequireRole("STUDENT"), attendanceHandler.History)
	attGroup.Get("", middleware.RequireRole("ADMIN"), attendanceHandler.List)
	attGroup.Get("/:id", attendanceHandler.GetByID)

	// Asset Routes (Memerlukan autentikasi)
	assetGroup := v1.Group("/assets", middleware.Authenticate(authService))
	assetGroup.Get("/my-allocations", allocHandler.MyAllocations)
	assetGroup.Get("", assetHandler.List)
	assetGroup.Get("/:id/history", allocHandler.AssetHistory)
	assetGroup.Get("/:id", assetHandler.GetByID)
	assetGroup.Post("", middleware.RequireRole("ADMIN"), assetHandler.Create)
	assetGroup.Put("/:id", middleware.RequireRole("ADMIN"), assetHandler.Update)
	assetGroup.Patch("/:id/status", middleware.RequireRole("ADMIN"), assetHandler.UpdateStatus)
	assetGroup.Delete("/:id", middleware.RequireRole("ADMIN"), assetHandler.Delete)

	// Allocation Routes (Memerlukan autentikasi)
	allocGroup := v1.Group("/allocations", middleware.Authenticate(authService))
	allocGroup.Post("", middleware.RequireRole("ADMIN"), allocHandler.Create)
	allocGroup.Get("", middleware.RequireRole("ADMIN"), allocHandler.List)
	allocGroup.Get("/my", allocHandler.MyAllocations)
	allocGroup.Get("/:id", allocHandler.GetByID)
	allocGroup.Post("/:id/return", middleware.RequireRole("ADMIN"), allocHandler.Return)

	// Maintenance Routes (Admin only)
	maintGroup := v1.Group("/maintenance", middleware.Authenticate(authService), middleware.RequireRole("ADMIN"))
	maintGroup.Post("", maintHandler.Create)
	maintGroup.Get("", maintHandler.List)
	maintGroup.Get("/:id", maintHandler.GetByID)
	maintGroup.Put("/:id", maintHandler.Update)
	maintGroup.Patch("/:id", maintHandler.Update)

	// Audit Logs Route (Admin only, append-only inspection)
	v1.Get("/audit-logs", middleware.Authenticate(authService), middleware.RequireRole("ADMIN"), auditHandler.List)

	// Dashboard Routes
	dashGroup := v1.Group("/dashboard", middleware.Authenticate(authService))
	dashGroup.Get("/student", middleware.RequireRole("STUDENT"), dashHandler.Student)
	dashGroup.Get("/admin", middleware.RequireRole("ADMIN"), dashHandler.Admin)

	// Report Routes (Admin only)
	repGroup := v1.Group("/reports", middleware.Authenticate(authService), middleware.RequireRole("ADMIN"))
	repGroup.Get("/attendance", reportHandler.Attendance)
	repGroup.Get("/assets", reportHandler.Assets)
	repGroup.Get("/allocations", reportHandler.Allocations)
	repGroup.Get("/maintenance", reportHandler.Maintenance)

	// Fallback 404
	app.Use(func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})
}

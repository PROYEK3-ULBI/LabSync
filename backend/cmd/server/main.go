package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"labsync-backend/internal/config"
	"labsync-backend/internal/middleware"
	"labsync-backend/internal/repository"
	"labsync-backend/internal/router"
)

func main() {
	// Logger terstruktur default menggunakan JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Muat dan validasi konfigurasi environment
	cfg, err := config.Load()
	if err != nil {
		slog.Error("konfigurasi tidak valid", "error", err.Error())
		os.Exit(1)
	}

	// Inisialisasi pool koneksi database
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		slog.Warn("koneksi database awal bermasalah", "error", err.Error())
	}
	defer func() {
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				slog.Error("gagal menutup connection pool database", "error", closeErr.Error())
			}
		}
	}()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName:      "LabSync Backend v1.0",
		BodyLimit:    cfg.BodyLimitBytes(),
		ErrorHandler: middleware.CentralErrorHandler,
	})

	// Setup routing dan middleware
	router.Setup(app, cfg, db)

	// Channel sinyal OS untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Jalankan server di goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		slog.Info("server dimulai", "port", cfg.Port, "env", cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			slog.Info("server berhenti menerima koneksi", "info", err.Error())
		}
	}()

	// Menunggu sinyal shutdown
	sig := <-quit
	slog.Info("sinyal shutdown diterima, memulai graceful shutdown", "signal", sig.String())

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("graceful shutdown gagal", "error", err.Error())
	} else {
		slog.Info("server berhasil berhenti dengan bersih")
	}
}

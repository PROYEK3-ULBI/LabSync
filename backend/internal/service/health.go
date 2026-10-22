package service

import (
	"context"
	"log/slog"
	"time"

	"labsync-backend/internal/repository"
)

// HealthResponse format respons status kesehatan layanan
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

type HealthService struct {
	repo *repository.HealthRepository
}

func NewHealthService(repo *repository.HealthRepository) *HealthService {
	return &HealthService{repo: repo}
}

// Check memeriksa kesehatan database tanpa mengekspos detail error internal
func (s *HealthService) Check(ctx context.Context) (HealthResponse, bool) {
	services := make(map[string]string)
	isHealthy := true

	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.repo.CheckDB(ctxTimeout); err != nil {
		slog.Warn("pemeriksaan health database gagal", "error", err.Error())
		services["database"] = "disconnected"
		isHealthy = false
	} else {
		services["database"] = "connected"
	}

	status := "healthy"
	if !isHealthy {
		status = "degraded"
	}

	return HealthResponse{
		Status:   status,
		Services: services,
	}, isHealthy
}

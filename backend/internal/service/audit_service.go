package service

import (
	"context"
	"time"

	"labsync-backend/internal/model"
	"labsync-backend/internal/repository"
)

type AuditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

// List mengembalikan log audit append-only untuk kebutuhan pengawasan administrator
func (s *AuditService) List(ctx context.Context, q model.AuditLogListQuery) (*model.AuditLogListResponse, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	list, total, err := s.auditRepo.List(ctx, q)
	if err != nil {
		return nil, model.NewInternalError("Gagal mengambil log audit", err)
	}

	data := make([]model.AuditLogResponse, len(list))
	for i, item := range list {
		data[i] = model.AuditLogResponse{
			ID:             item.ID,
			UserID:         item.UserID,
			Action:         item.Action,
			EntityType:     item.EntityType,
			EntityID:       item.EntityID,
			Status:         item.Status,
			IPAddress:      item.IPAddress,
			UserAgent:      item.UserAgent,
			Metadata:       item.Metadata,
			CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
			UserIdentifier: item.UserIdentifier,
			UserFullName:   item.UserFullName,
		}
	}

	return &model.AuditLogListResponse{
		Data: data,
		Meta: model.PaginationMeta{
			Page:  q.Page,
			Limit: q.Limit,
			Total: total,
		},
	}, nil
}

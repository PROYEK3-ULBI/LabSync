package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"labsync-backend/internal/model"
)

type AuditRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	List(ctx context.Context, query model.AuditLogListQuery) ([]model.AuditLog, int, error)
}

type PostgresAuditRepository struct {
	db *DB
}

func NewPostgresAuditRepository(db *DB) *PostgresAuditRepository {
	return &PostgresAuditRepository{db: db}
}

func (r *PostgresAuditRepository) Create(ctx context.Context, a *model.AuditLog) error {
	userIDVal := "nil"
	if a.UserID != nil {
		userIDVal = *a.UserID
	}
	slog.Info("security audit event",
		"action", a.Action,
		"entity_type", a.EntityType,
		"entity_id", a.EntityID,
		"user_id", userIDVal,
	)

	if r.db == nil || r.db.SQL == nil {
		return nil
	}

	query := `
		INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, status, ip_address, user_agent, metadata, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at`

	var userID sql.NullString
	if a.UserID != nil {
		userID = sql.NullString{String: *a.UserID, Valid: true}
	}

	status := a.Status
	if status == "" {
		status = "SUCCESS"
	}

	var metaBytes []byte
	if len(a.Metadata) > 0 {
		metaBytes, _ = json.Marshal(a.Metadata)
	}

	err := r.db.SQL.QueryRowContext(ctx, query,
		userID, a.Action, a.EntityType, a.EntityID, status, a.IPAddress, a.UserAgent, metaBytes,
	).Scan(&a.ID, &a.CreatedAt)

	if err != nil {
		slog.Warn("gagal menyimpan audit log ke database", "error", err.Error())
		return err
	}
	return nil
}

func (r *PostgresAuditRepository) List(ctx context.Context, q model.AuditLogListQuery) ([]model.AuditLog, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, nil
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.Action != "" {
		conditions = append(conditions, fmt.Sprintf("al.action = $%d", argIdx))
		args = append(args, q.Action)
		argIdx++
	}
	if q.EntityType != "" {
		conditions = append(conditions, fmt.Sprintf("al.entity_type = $%d", argIdx))
		args = append(args, q.EntityType)
		argIdx++
	}
	if q.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("al.user_id = $%d", argIdx))
		args = append(args, q.UserID)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM audit_logs al WHERE %s`, whereClause)
	var total int
	if err := r.db.SQL.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	offset := (q.Page - 1) * q.Limit

	dataQuery := fmt.Sprintf(`
		SELECT al.id, al.user_id, al.action, al.entity_type, al.entity_id, al.status,
		       al.ip_address, al.user_agent, al.metadata, al.created_at,
		       COALESCE(u.identifier, ''), COALESCE(u.full_name, '')
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE %s
		ORDER BY al.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.AuditLog{}
	for rows.Next() {
		var a model.AuditLog
		var uid, ip, ua sql.NullString
		var metaBytes []byte
		var ident, name string

		if err := rows.Scan(
			&a.ID, &uid, &a.Action, &a.EntityType, &a.EntityID, &a.Status,
			&ip, &ua, &metaBytes, &a.CreatedAt,
			&ident, &name,
		); err != nil {
			return nil, 0, err
		}

		if uid.Valid {
			a.UserID = &uid.String
		}
		if ip.Valid {
			a.IPAddress = &ip.String
		}
		if ua.Valid {
			a.UserAgent = &ua.String
		}
		if len(metaBytes) > 0 {
			_ = json.Unmarshal(metaBytes, &a.Metadata)
		}
		a.UserIdentifier = ident
		a.UserFullName = name
		a.CreatedAt = a.CreatedAt.UTC()

		list = append(list, a)
	}

	return list, total, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"labsync-backend/internal/model"
)

var (
	ErrMaintenanceNotFound      = errors.New("tiket pemeliharaan tidak ditemukan")
	ErrInvalidStatusTransition  = errors.New("transisi status pemeliharaan tidak diizinkan")
	ErrMaintenanceAlreadyClosed = errors.New("tiket pemeliharaan sudah selesai atau dibatalkan")
)

type MaintenanceRepository interface {
	CreateWithTransaction(ctx context.Context, m *model.Maintenance) error
	UpdateWithTransaction(ctx context.Context, m *model.Maintenance, resultAssetStatus *model.AssetStatus) error
	GetByID(ctx context.Context, id string) (*model.Maintenance, error)
	List(ctx context.Context, query model.MaintenanceListQuery) ([]model.Maintenance, int, error)
}

type PostgresMaintenanceRepository struct {
	db *DB
}

func NewPostgresMaintenanceRepository(db *DB) *PostgresMaintenanceRepository {
	return &PostgresMaintenanceRepository{db: db}
}

func (r *PostgresMaintenanceRepository) CreateWithTransaction(ctx context.Context, m *model.Maintenance) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	tx, err := r.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi pemeliharaan: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Kunci baris aset dan periksa status & alokasi aktif
	var currentStatus string
	lockAssetQuery := `SELECT status FROM assets WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, lockAssetQuery, m.AssetID).Scan(&currentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAssetNotFound
	}
	if err != nil {
		return err
	}

	if currentStatus == string(model.AssetStatusAllocated) {
		return ErrAssetHasActiveAllocation
	}
	if currentStatus == string(model.AssetStatusRetired) {
		return errors.New("aset berstatus RETIRED tidak dapat diajukan pemeliharaan")
	}

	// Periksa apakah ada alokasi yang belum dikembalikan
	var hasAlloc bool
	checkAllocQuery := `SELECT EXISTS(SELECT 1 FROM asset_allocations WHERE asset_id = $1 AND returned_at IS NULL)`
	if err := tx.QueryRowContext(ctx, checkAllocQuery, m.AssetID).Scan(&hasAlloc); err != nil {
		return err
	}
	if hasAlloc {
		return ErrAssetHasActiveAllocation
	}

	// 2. Simpan tiket pemeliharaan
	insertQuery := `
		INSERT INTO maintenance (id, asset_id, reported_by, issue_type, description, priority, status, started_at, completed_at, resolution, cost, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'OPEN', NULL, NULL, NULL, NULL, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, insertQuery,
		m.AssetID, m.ReportedBy, m.IssueType, m.Description, string(m.Priority),
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal menyimpan data pemeliharaan: %w", err)
	}

	// 3. Ubah status aset menjadi MAINTENANCE
	updateAssetQuery := `UPDATE assets SET status = 'MAINTENANCE', updated_at = NOW() WHERE id = $1`
	if _, err := tx.ExecContext(ctx, updateAssetQuery, m.AssetID); err != nil {
		return fmt.Errorf("gagal memperbarui status aset ke MAINTENANCE: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi pemeliharaan: %w", err)
	}

	m.CreatedAt = m.CreatedAt.UTC()
	m.UpdatedAt = m.UpdatedAt.UTC()
	m.Status = model.MaintenanceStatusOpen

	return nil
}

func (r *PostgresMaintenanceRepository) UpdateWithTransaction(ctx context.Context, m *model.Maintenance, resultAssetStatus *model.AssetStatus) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	tx, err := r.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi update pemeliharaan: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Kunci tiket pemeliharaan dan periksa status lama
	var oldStatusStr, assetID string
	var startedAt, completedAt sql.NullTime
	lockMaintQuery := `SELECT status, asset_id, started_at, completed_at FROM maintenance WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, lockMaintQuery, m.ID).Scan(&oldStatusStr, &assetID, &startedAt, &completedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMaintenanceNotFound
	}
	if err != nil {
		return err
	}

	oldStatus := model.MaintenanceStatus(oldStatusStr)
	newStatus := m.Status

	// 2. Validasi transisi status
	if !IsValidMaintenanceTransition(oldStatus, newStatus) {
		return ErrInvalidStatusTransition
	}

	// 3. Kunci baris aset
	lockAssetQuery := `SELECT id FROM assets WHERE id = $1 FOR UPDATE`
	if _, err := tx.ExecContext(ctx, lockAssetQuery, assetID); err != nil {
		return err
	}

	now := time.Now().UTC()
	var finalStartedAt, finalCompletedAt *time.Time
	if startedAt.Valid {
		t := startedAt.Time.UTC()
		finalStartedAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time.UTC()
		finalCompletedAt = &t
	}

	if newStatus == model.MaintenanceStatusInProgress && finalStartedAt == nil {
		finalStartedAt = &now
	}

	if newStatus == model.MaintenanceStatusCompleted {
		finalCompletedAt = &now
		// Update status aset ke AVAILABLE atau RETIRED
		targetAssetStatus := model.AssetStatusAvailable
		if resultAssetStatus != nil && *resultAssetStatus == model.AssetStatusRetired {
			targetAssetStatus = model.AssetStatusRetired
		}
		updateAssetQuery := `UPDATE assets SET status = $1, updated_at = NOW() WHERE id = $2`
		if _, err := tx.ExecContext(ctx, updateAssetQuery, string(targetAssetStatus), assetID); err != nil {
			return err
		}
	} else if newStatus == model.MaintenanceStatusCancelled {
		// Revert status aset ke AVAILABLE jika tidak ada tiket pemeliharaan lain yang OPEN/IN_PROGRESS
		var otherOpen bool
		checkOtherQuery := `SELECT EXISTS(SELECT 1 FROM maintenance WHERE asset_id = $1 AND id != $2 AND status IN ('OPEN', 'IN_PROGRESS'))`
		if err := tx.QueryRowContext(ctx, checkOtherQuery, assetID, m.ID).Scan(&otherOpen); err != nil {
			return err
		}
		if !otherOpen {
			updateAssetQuery := `UPDATE assets SET status = 'AVAILABLE', updated_at = NOW() WHERE id = $1`
			if _, err := tx.ExecContext(ctx, updateAssetQuery, assetID); err != nil {
				return err
			}
		}
	}

	// 4. Update data tiket pemeliharaan
	updateMaintQuery := `
		UPDATE maintenance
		SET status = $1, description = $2, priority = $3, started_at = $4, completed_at = $5,
		    resolution = $6, cost = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at`

	err = tx.QueryRowContext(ctx, updateMaintQuery,
		string(newStatus), m.Description, string(m.Priority), finalStartedAt, finalCompletedAt,
		m.Resolution, m.Cost, m.ID,
	).Scan(&m.UpdatedAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi pemeliharaan: %w", err)
	}

	m.StartedAt = finalStartedAt
	m.CompletedAt = finalCompletedAt
	m.UpdatedAt = m.UpdatedAt.UTC()

	return nil
}

func (r *PostgresMaintenanceRepository) GetByID(ctx context.Context, id string) (*model.Maintenance, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT mn.id, mn.asset_id, mn.reported_by, mn.issue_type, mn.description, mn.priority,
		       mn.status, mn.started_at, mn.completed_at, mn.resolution, mn.cost, mn.created_at, mn.updated_at,
		       ast.asset_code, ast.asset_name, u.identifier, u.full_name
		FROM maintenance mn
		JOIN assets ast ON mn.asset_id = ast.id
		JOIN users u ON mn.reported_by = u.id
		WHERE mn.id = $1
		LIMIT 1`

	var m model.Maintenance
	var startedAt, completedAt sql.NullTime
	var resolution sql.NullString
	var cost sql.NullFloat64
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.AssetID, &m.ReportedBy, &m.IssueType, &m.Description, &m.Priority,
		&m.Status, &startedAt, &completedAt, &resolution, &cost, &m.CreatedAt, &m.UpdatedAt,
		&m.AssetCode, &m.AssetName, &m.ReporterIdentifier, &m.ReporterName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMaintenanceNotFound
	}
	if err != nil {
		return nil, err
	}

	if startedAt.Valid {
		t := startedAt.Time.UTC()
		m.StartedAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time.UTC()
		m.CompletedAt = &t
	}
	if resolution.Valid {
		m.Resolution = &resolution.String
	}
	if cost.Valid {
		m.Cost = &cost.Float64
	}

	m.CreatedAt = m.CreatedAt.UTC()
	m.UpdatedAt = m.UpdatedAt.UTC()

	return &m, nil
}

func (r *PostgresMaintenanceRepository) List(ctx context.Context, q model.MaintenanceListQuery) ([]model.Maintenance, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, errors.New("database tidak tersedia")
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.AssetID != "" {
		conditions = append(conditions, fmt.Sprintf("mn.asset_id = $%d", argIdx))
		args = append(args, q.AssetID)
		argIdx++
	}
	if q.Status != "" {
		conditions = append(conditions, fmt.Sprintf("mn.status = $%d", argIdx))
		args = append(args, strings.ToUpper(q.Status))
		argIdx++
	}
	if q.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("mn.priority = $%d", argIdx))
		args = append(args, strings.ToUpper(q.Priority))
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM maintenance mn WHERE %s`, whereClause)
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
		SELECT mn.id, mn.asset_id, mn.reported_by, mn.issue_type, mn.description, mn.priority,
		       mn.status, mn.started_at, mn.completed_at, mn.resolution, mn.cost, mn.created_at, mn.updated_at,
		       ast.asset_code, ast.asset_name, u.identifier, u.full_name
		FROM maintenance mn
		JOIN assets ast ON mn.asset_id = ast.id
		JOIN users u ON mn.reported_by = u.id
		WHERE %s
		ORDER BY mn.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.Maintenance{}
	for rows.Next() {
		var m model.Maintenance
		var startedAt, completedAt sql.NullTime
		var resolution sql.NullString
		var cost sql.NullFloat64

		if err := rows.Scan(
			&m.ID, &m.AssetID, &m.ReportedBy, &m.IssueType, &m.Description, &m.Priority,
			&m.Status, &startedAt, &completedAt, &resolution, &cost, &m.CreatedAt, &m.UpdatedAt,
			&m.AssetCode, &m.AssetName, &m.ReporterIdentifier, &m.ReporterName,
		); err != nil {
			return nil, 0, err
		}

		if startedAt.Valid {
			t := startedAt.Time.UTC()
			m.StartedAt = &t
		}
		if completedAt.Valid {
			t := completedAt.Time.UTC()
			m.CompletedAt = &t
		}
		if resolution.Valid {
			m.Resolution = &resolution.String
		}
		if cost.Valid {
			m.Cost = &cost.Float64
		}

		m.CreatedAt = m.CreatedAt.UTC()
		m.UpdatedAt = m.UpdatedAt.UTC()

		list = append(list, m)
	}

	return list, total, nil
}

func IsValidMaintenanceTransition(from, to model.MaintenanceStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case model.MaintenanceStatusOpen:
		return to == model.MaintenanceStatusInProgress || to == model.MaintenanceStatusCancelled
	case model.MaintenanceStatusInProgress:
		return to == model.MaintenanceStatusCompleted || to == model.MaintenanceStatusCancelled
	case model.MaintenanceStatusCompleted, model.MaintenanceStatusCancelled:
		// Terminal states
		return false
	default:
		return false
	}
}

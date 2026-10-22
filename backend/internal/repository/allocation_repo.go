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
	ErrAllocationNotFound = errors.New("alokasi aset tidak ditemukan")
	ErrAssetNotAvailable  = errors.New("aset tidak berstatus AVAILABLE")
	ErrDoubleAllocation   = errors.New("aset sedang dialokasikan dan belum dikembalikan")
	ErrAlreadyReturned    = errors.New("alokasi aset sudah dikembalikan sebelumnya")
)

type AllocationRepository interface {
	CreateWithTransaction(ctx context.Context, alloc *model.Allocation) error
	ReturnWithTransaction(ctx context.Context, allocationID string, returnedAt time.Time, conditionAfter string, notes *string) (*model.Allocation, error)
	GetByID(ctx context.Context, id string) (*model.Allocation, error)
	List(ctx context.Context, query model.AllocationListQuery) ([]model.Allocation, int, error)
	GetActiveAllocationByAssetID(ctx context.Context, assetID string) (*model.Allocation, error)
	GetAssetTimeline(ctx context.Context, assetID string) ([]model.AssetTimelineEvent, error)
}

type PostgresAllocationRepository struct {
	db *DB
}

func NewPostgresAllocationRepository(db *DB) *PostgresAllocationRepository {
	return &PostgresAllocationRepository{db: db}
}

func (r *PostgresAllocationRepository) CreateWithTransaction(ctx context.Context, a *model.Allocation) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	tx, err := r.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi alokasi: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Kunci baris aset menggunakan SELECT ... FOR UPDATE untuk mencegah race condition
	var currentStatus string
	checkQuery := `SELECT status FROM assets WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, checkQuery, a.AssetID).Scan(&currentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAssetNotFound
	}
	if err != nil {
		return err
	}

	if currentStatus != string(model.AssetStatusAvailable) {
		return ErrAssetNotAvailable
	}

	// 2. Simpan alokasi aset
	insertQuery := `
		INSERT INTO asset_allocations (id, asset_id, user_id, session_id, allocated_at, returned_at, condition_before, notes, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, NULL, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, insertQuery,
		a.AssetID, a.UserID, a.SessionID, a.AllocatedAt.UTC(), a.ConditionBefore, a.Notes,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "idx_asset_allocations_active") || strings.Contains(err.Error(), "23505") {
			return ErrDoubleAllocation
		}
		return err
	}

	// 3. Update status aset menjadi ALLOCATED dalam transaksi yang sama
	updateQuery := `UPDATE assets SET status = 'ALLOCATED', updated_at = NOW() WHERE id = $1`
	if _, err := tx.ExecContext(ctx, updateQuery, a.AssetID); err != nil {
		return fmt.Errorf("gagal mengubah status aset menjadi ALLOCATED: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi alokasi: %w", err)
	}

	a.AllocatedAt = a.AllocatedAt.UTC()
	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return nil
}

func (r *PostgresAllocationRepository) ReturnWithTransaction(
	ctx context.Context,
	allocationID string,
	returnedAt time.Time,
	conditionAfter string,
	notes *string,
) (*model.Allocation, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	tx, err := r.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi pengembalian: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Kunci baris alokasi
	var assetID string
	var existingReturnedAt sql.NullTime
	selectQuery := `SELECT asset_id, returned_at FROM asset_allocations WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, selectQuery, allocationID).Scan(&assetID, &existingReturnedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAllocationNotFound
	}
	if err != nil {
		return nil, err
	}

	if existingReturnedAt.Valid {
		return nil, ErrAlreadyReturned
	}

	// 2. Kunci baris aset
	lockAssetQuery := `SELECT id FROM assets WHERE id = $1 FOR UPDATE`
	if _, err := tx.ExecContext(ctx, lockAssetQuery, assetID); err != nil {
		return nil, err
	}

	// 3. Update pengembalian alokasi
	updateAllocQuery := `
		UPDATE asset_allocations
		SET returned_at = $1, condition_after = $2, notes = COALESCE($3, notes), updated_at = NOW()
		WHERE id = $4`

	if _, err := tx.ExecContext(ctx, updateAllocQuery, returnedAt.UTC(), conditionAfter, notes, allocationID); err != nil {
		return nil, fmt.Errorf("gagal mencatat data pengembalian alokasi: %w", err)
	}

	// 4. Tentukan status aset setelah return
	targetStatus := model.AssetStatusAvailable
	if strings.EqualFold(conditionAfter, "DAMAGED") {
		targetStatus = model.AssetStatusMaintenance
	}

	updateAssetQuery := `UPDATE assets SET status = $1, updated_at = NOW() WHERE id = $2`
	if _, err := tx.ExecContext(ctx, updateAssetQuery, string(targetStatus), assetID); err != nil {
		return nil, fmt.Errorf("gagal memperbarui status aset: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("gagal commit transaksi pengembalian aset: %w", err)
	}

	return r.GetByID(ctx, allocationID)
}

func (r *PostgresAllocationRepository) GetByID(ctx context.Context, id string) (*model.Allocation, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT al.id, al.asset_id, al.user_id, al.session_id, al.allocated_at, al.returned_at,
		       al.condition_before, al.condition_after, al.notes, al.created_at, al.updated_at,
		       ast.asset_code, ast.asset_name, u.identifier, u.full_name, s.code, s.title
		FROM asset_allocations al
		JOIN assets ast ON al.asset_id = ast.id
		JOIN users u ON al.user_id = u.id
		JOIN lab_sessions s ON al.session_id = s.id
		WHERE al.id = $1
		LIMIT 1`

	var a model.Allocation
	var returnedAt sql.NullTime
	var conditionAfter, notes sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.AssetID, &a.UserID, &a.SessionID, &a.AllocatedAt, &returnedAt,
		&a.ConditionBefore, &conditionAfter, &notes, &a.CreatedAt, &a.UpdatedAt,
		&a.AssetCode, &a.AssetName, &a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAllocationNotFound
	}
	if err != nil {
		return nil, err
	}

	if returnedAt.Valid {
		t := returnedAt.Time.UTC()
		a.ReturnedAt = &t
	}
	if conditionAfter.Valid {
		a.ConditionAfter = &conditionAfter.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	a.AllocatedAt = a.AllocatedAt.UTC()
	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAllocationRepository) GetActiveAllocationByAssetID(ctx context.Context, assetID string) (*model.Allocation, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT al.id, al.asset_id, al.user_id, al.session_id, al.allocated_at, al.returned_at,
		       al.condition_before, al.condition_after, al.notes, al.created_at, al.updated_at,
		       ast.asset_code, ast.asset_name, u.identifier, u.full_name, s.code, s.title
		FROM asset_allocations al
		JOIN assets ast ON al.asset_id = ast.id
		JOIN users u ON al.user_id = u.id
		JOIN lab_sessions s ON al.session_id = s.id
		WHERE al.asset_id = $1 AND al.returned_at IS NULL
		LIMIT 1`

	var a model.Allocation
	var returnedAt sql.NullTime
	var conditionAfter, notes sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, assetID).Scan(
		&a.ID, &a.AssetID, &a.UserID, &a.SessionID, &a.AllocatedAt, &returnedAt,
		&a.ConditionBefore, &conditionAfter, &notes, &a.CreatedAt, &a.UpdatedAt,
		&a.AssetCode, &a.AssetName, &a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAllocationNotFound
	}
	if err != nil {
		return nil, err
	}

	if returnedAt.Valid {
		t := returnedAt.Time.UTC()
		a.ReturnedAt = &t
	}
	if conditionAfter.Valid {
		a.ConditionAfter = &conditionAfter.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	return &a, nil
}

func (r *PostgresAllocationRepository) List(ctx context.Context, q model.AllocationListQuery) ([]model.Allocation, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, errors.New("database tidak tersedia")
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.AssetID != "" {
		conditions = append(conditions, fmt.Sprintf("al.asset_id = $%d", argIdx))
		args = append(args, q.AssetID)
		argIdx++
	}
	if q.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("al.user_id = $%d", argIdx))
		args = append(args, q.UserID)
		argIdx++
	}
	if q.SessionID != "" {
		conditions = append(conditions, fmt.Sprintf("al.session_id = $%d", argIdx))
		args = append(args, q.SessionID)
		argIdx++
	}
	if strings.EqualFold(q.Status, "ACTIVE") {
		conditions = append(conditions, "al.returned_at IS NULL")
	} else if strings.EqualFold(q.Status, "RETURNED") {
		conditions = append(conditions, "al.returned_at IS NOT NULL")
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM asset_allocations al WHERE %s`, whereClause)
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
		SELECT al.id, al.asset_id, al.user_id, al.session_id, al.allocated_at, al.returned_at,
		       al.condition_before, al.condition_after, al.notes, al.created_at, al.updated_at,
		       ast.asset_code, ast.asset_name, u.identifier, u.full_name, s.code, s.title
		FROM asset_allocations al
		JOIN assets ast ON al.asset_id = ast.id
		JOIN users u ON al.user_id = u.id
		JOIN lab_sessions s ON al.session_id = s.id
		WHERE %s
		ORDER BY al.allocated_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.Allocation{}
	for rows.Next() {
		var a model.Allocation
		var returnedAt sql.NullTime
		var conditionAfter, notes sql.NullString
		if err := rows.Scan(
			&a.ID, &a.AssetID, &a.UserID, &a.SessionID, &a.AllocatedAt, &returnedAt,
			&a.ConditionBefore, &conditionAfter, &notes, &a.CreatedAt, &a.UpdatedAt,
			&a.AssetCode, &a.AssetName, &a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle,
		); err != nil {
			return nil, 0, err
		}
		if returnedAt.Valid {
			t := returnedAt.Time.UTC()
			a.ReturnedAt = &t
		}
		if conditionAfter.Valid {
			a.ConditionAfter = &conditionAfter.String
		}
		if notes.Valid {
			a.Notes = &notes.String
		}
		a.AllocatedAt = a.AllocatedAt.UTC()
		a.CreatedAt = a.CreatedAt.UTC()
		a.UpdatedAt = a.UpdatedAt.UTC()

		list = append(list, a)
	}

	return list, total, nil
}

func (r *PostgresAllocationRepository) GetAssetTimeline(ctx context.Context, assetID string) ([]model.AssetTimelineEvent, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	// 1. Ambil data alokasi dan pengembalian
	allocQuery := `
		SELECT al.id, al.allocated_at, al.returned_at, al.condition_before, al.condition_after, al.notes,
		       u.id, u.full_name, u.identifier, s.code, s.title
		FROM asset_allocations al
		JOIN users u ON al.user_id = u.id
		JOIN lab_sessions s ON al.session_id = s.id
		WHERE al.asset_id = $1
		ORDER BY al.allocated_at DESC`

	rows, err := r.db.SQL.QueryContext(ctx, allocQuery, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeline []model.AssetTimelineEvent
	for rows.Next() {
		var id, condBefore string
		var returnedAt sql.NullTime
		var condAfter, notes sql.NullString
		var userID, fullName, identifier, sessCode, sessTitle string
		var allocatedAt time.Time

		if err := rows.Scan(
			&id, &allocatedAt, &returnedAt, &condBefore, &condAfter, &notes,
			&userID, &fullName, &identifier, &sessCode, &sessTitle,
		); err != nil {
			return nil, err
		}

		// Event Peminjaman
		timeline = append(timeline, model.AssetTimelineEvent{
			ID:          fmt.Sprintf("alloc-%s", id),
			EventType:   "ALLOCATION",
			Title:       "Aset Dipinjam",
			Description: fmt.Sprintf("Aset dialokasikan kepada %s (%s) pada sesi %s", fullName, identifier, sessCode),
			Timestamp:   allocatedAt.UTC().Format(time.RFC3339),
			ActorID:     &userID,
			ActorName:   &fullName,
			Details: map[string]interface{}{
				"condition_before": condBefore,
				"session_code":     sessCode,
				"session_title":    sessTitle,
			},
		})

		// Event Pengembalian jika sudah dikembalikan
		if returnedAt.Valid {
			condAfterVal := "GOOD"
			if condAfter.Valid {
				condAfterVal = condAfter.String
			}
			timeline = append(timeline, model.AssetTimelineEvent{
				ID:          fmt.Sprintf("return-%s", id),
				EventType:   "RETURN",
				Title:       "Aset Dikembalikan",
				Description: fmt.Sprintf("Aset dikembalikan oleh %s (%s) dengan kondisi %s", fullName, identifier, condAfterVal),
				Timestamp:   returnedAt.Time.UTC().Format(time.RFC3339),
				ActorID:     &userID,
				ActorName:   &fullName,
				Details: map[string]interface{}{
					"condition_after": condAfterVal,
					"session_code":    sessCode,
				},
			})
		}
	}

	// 2. Ambil data pemeliharaan untuk aset ini
	maintQuery := `
		SELECT mn.id, mn.status, mn.priority, mn.issue_type, mn.description, mn.created_at, mn.resolution,
		       u.id, u.full_name, u.identifier
		FROM maintenance mn
		JOIN users u ON mn.reported_by = u.id
		WHERE mn.asset_id = $1
		ORDER BY mn.created_at DESC`

	maintRows, err := r.db.SQL.QueryContext(ctx, maintQuery, assetID)
	if err == nil {
		defer maintRows.Close()
		for maintRows.Next() {
			var id, status, priority, issueType, description string
			var createdAt time.Time
			var resolution sql.NullString
			var reporterID, reporterName, reporterIdent string

			if err := maintRows.Scan(
				&id, &status, &priority, &issueType, &description, &createdAt,
				&resolution, &reporterID, &reporterName, &reporterIdent,
			); err == nil {
				timeline = append(timeline, model.AssetTimelineEvent{
					ID:          fmt.Sprintf("maint-%s", id),
					EventType:   "MAINTENANCE",
					Title:       fmt.Sprintf("Pemeliharaan (%s - %s)", issueType, status),
					Description: description,
					Timestamp:   createdAt.UTC().Format(time.RFC3339),
					ActorID:     &reporterID,
					ActorName:   &reporterName,
					Details: map[string]interface{}{
						"priority":   priority,
						"status":     status,
						"issue_type": issueType,
					},
				})
			}
		}
	}

	return timeline, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"labsync-backend/internal/model"
)

var (
	ErrAssetNotFound             = errors.New("aset tidak ditemukan")
	ErrDuplicateAssetCode        = errors.New("kode aset sudah terdaftar")
	ErrDuplicateSerialNumber     = errors.New("nomor seri sudah terdaftar pada aset lain")
	ErrAssetHasActiveAllocation  = errors.New("aset sedang dialokasikan dan tidak dapat diubah statusnya atau dihapus secara langsung")
	ErrAssetHasActiveMaintenance = errors.New("aset sedang dalam proses pemeliharaan aktif dan tidak dapat diubah statusnya atau dihapus")
)

type AssetRepository interface {
	GetByID(ctx context.Context, id string) (*model.Asset, error)
	GetByCode(ctx context.Context, code string) (*model.Asset, error)
	GetBySerialNumber(ctx context.Context, sn string) (*model.Asset, error)
	Create(ctx context.Context, asset *model.Asset) error
	Update(ctx context.Context, asset *model.Asset) error
	UpdateStatus(ctx context.Context, id string, status model.AssetStatus) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, query model.AssetListQuery) ([]model.Asset, int, error)
	HasActiveAllocation(ctx context.Context, assetID string) (bool, error)
	HasActiveMaintenance(ctx context.Context, assetID string) (bool, error)
}

type PostgresAssetRepository struct {
	db *DB
}

func NewPostgresAssetRepository(db *DB) *PostgresAssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (r *PostgresAssetRepository) GetByID(ctx context.Context, id string) (*model.Asset, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, asset_code, asset_name, category, brand, model, serial_number, location, status, notes, created_at, updated_at
		FROM assets
		WHERE id = $1
		LIMIT 1`

	var a model.Asset
	var brand, modelName, sn, notes sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.AssetCode, &a.AssetName, &a.Category, &brand, &modelName, &sn, &a.Location, &a.Status, &notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetNotFound
	}
	if err != nil {
		return nil, err
	}

	if brand.Valid {
		a.Brand = &brand.String
	}
	if modelName.Valid {
		a.Model = &modelName.String
	}
	if sn.Valid {
		a.SerialNumber = &sn.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAssetRepository) GetByCode(ctx context.Context, code string) (*model.Asset, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, asset_code, asset_name, category, brand, model, serial_number, location, status, notes, created_at, updated_at
		FROM assets
		WHERE asset_code = $1
		LIMIT 1`

	var a model.Asset
	var brand, modelName, sn, notes sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, code).Scan(
		&a.ID, &a.AssetCode, &a.AssetName, &a.Category, &brand, &modelName, &sn, &a.Location, &a.Status, &notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetNotFound
	}
	if err != nil {
		return nil, err
	}

	if brand.Valid {
		a.Brand = &brand.String
	}
	if modelName.Valid {
		a.Model = &modelName.String
	}
	if sn.Valid {
		a.SerialNumber = &sn.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAssetRepository) GetBySerialNumber(ctx context.Context, sn string) (*model.Asset, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, asset_code, asset_name, category, brand, model, serial_number, location, status, notes, created_at, updated_at
		FROM assets
		WHERE serial_number = $1
		LIMIT 1`

	var a model.Asset
	var brand, modelName, serialNum, notes sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, sn).Scan(
		&a.ID, &a.AssetCode, &a.AssetName, &a.Category, &brand, &modelName, &serialNum, &a.Location, &a.Status, &notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetNotFound
	}
	if err != nil {
		return nil, err
	}

	if brand.Valid {
		a.Brand = &brand.String
	}
	if modelName.Valid {
		a.Model = &modelName.String
	}
	if serialNum.Valid {
		a.SerialNumber = &serialNum.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}

	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAssetRepository) Create(ctx context.Context, a *model.Asset) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		INSERT INTO assets (id, asset_code, asset_name, category, brand, model, serial_number, location, status, notes, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.SQL.QueryRowContext(ctx, query,
		a.AssetCode, a.AssetName, a.Category, a.Brand, a.Model, a.SerialNumber, a.Location, string(a.Status), a.Notes,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "assets_asset_code_key") || strings.Contains(err.Error(), "asset_code") {
			return ErrDuplicateAssetCode
		}
		if strings.Contains(err.Error(), "assets_serial_number_key") || strings.Contains(err.Error(), "serial_number") {
			return ErrDuplicateSerialNumber
		}
		return err
	}

	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return nil
}

func (r *PostgresAssetRepository) Update(ctx context.Context, a *model.Asset) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		UPDATE assets
		SET asset_code = $1, asset_name = $2, category = $3, brand = $4, model = $5,
		    serial_number = $6, location = $7, status = $8, notes = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at`

	err := r.db.SQL.QueryRowContext(ctx, query,
		a.AssetCode, a.AssetName, a.Category, a.Brand, a.Model, a.SerialNumber,
		a.Location, string(a.Status), a.Notes, a.ID,
	).Scan(&a.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAssetNotFound
		}
		if strings.Contains(err.Error(), "assets_asset_code_key") || strings.Contains(err.Error(), "asset_code") {
			return ErrDuplicateAssetCode
		}
		if strings.Contains(err.Error(), "assets_serial_number_key") || strings.Contains(err.Error(), "serial_number") {
			return ErrDuplicateSerialNumber
		}
		return err
	}

	a.UpdatedAt = a.UpdatedAt.UTC()
	return nil
}

func (r *PostgresAssetRepository) UpdateStatus(ctx context.Context, id string, status model.AssetStatus) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `UPDATE assets SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.SQL.ExecContext(ctx, query, string(status), id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAssetNotFound
	}
	return nil
}

func (r *PostgresAssetRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `DELETE FROM assets WHERE id = $1`
	res, err := r.db.SQL.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAssetNotFound
	}
	return nil
}

func (r *PostgresAssetRepository) List(ctx context.Context, q model.AssetListQuery) ([]model.Asset, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, errors.New("database tidak tersedia")
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, strings.ToUpper(q.Status))
		argIdx++
	}

	if q.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category ILIKE $%d", argIdx))
		args = append(args, q.Category)
		argIdx++
	}

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(asset_code ILIKE $%d OR asset_name ILIKE $%d OR COALESCE(brand, '') ILIKE $%d OR COALESCE(serial_number, '') ILIKE $%d OR location ILIKE $%d)", argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM assets WHERE %s`, whereClause)
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
		SELECT id, asset_code, asset_name, category, brand, model, serial_number, location, status, notes, created_at, updated_at
		FROM assets
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.Asset{}
	for rows.Next() {
		var a model.Asset
		var brand, modelName, sn, notes sql.NullString
		if err := rows.Scan(
			&a.ID, &a.AssetCode, &a.AssetName, &a.Category, &brand, &modelName, &sn, &a.Location, &a.Status, &notes, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if brand.Valid {
			a.Brand = &brand.String
		}
		if modelName.Valid {
			a.Model = &modelName.String
		}
		if sn.Valid {
			a.SerialNumber = &sn.String
		}
		if notes.Valid {
			a.Notes = &notes.String
		}
		a.CreatedAt = a.CreatedAt.UTC()
		a.UpdatedAt = a.UpdatedAt.UTC()

		list = append(list, a)
	}

	return list, total, nil
}

func (r *PostgresAssetRepository) HasActiveAllocation(ctx context.Context, assetID string) (bool, error) {
	if r.db == nil || r.db.SQL == nil {
		return false, nil
	}

	query := `SELECT EXISTS(SELECT 1 FROM asset_allocations WHERE asset_id = $1 AND returned_at IS NULL)`
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx, query, assetID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresAssetRepository) HasActiveMaintenance(ctx context.Context, assetID string) (bool, error) {
	if r.db == nil || r.db.SQL == nil {
		return false, nil
	}

	query := `SELECT EXISTS(SELECT 1 FROM maintenance WHERE asset_id = $1 AND status IN ('OPEN', 'IN_PROGRESS'))`
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx, query, assetID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

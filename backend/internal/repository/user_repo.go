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
	ErrUserNotFound = errors.New("pengguna tidak ditemukan")
	ErrRoleNotFound = errors.New("peran tidak ditemukan")
	ErrDuplicateKey = errors.New("identifier atau email sudah terdaftar")
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, query model.UserListQuery) ([]model.User, int, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetRoleByCode(ctx context.Context, code string) (*model.Role, error)
	HasFaceTemplate(ctx context.Context, userID string) (bool, error)
}

type PostgresUserRepository struct {
	db *DB
}

func NewPostgresUserRepository(db *DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT u.id, u.role_id, r.code, u.identifier, u.email, u.full_name, u.password_hash, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1`

	var u model.User
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.RoleID, &u.RoleCode, &u.Identifier, &u.Email, &u.FullName, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT u.id, u.role_id, r.code, u.identifier, u.email, u.full_name, u.password_hash, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.identifier = $1`

	var u model.User
	err := r.db.SQL.QueryRowContext(ctx, query, identifier).Scan(
		&u.ID, &u.RoleID, &u.RoleCode, &u.Identifier, &u.Email, &u.FullName, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT u.id, u.role_id, r.code, u.identifier, u.email, u.full_name, u.password_hash, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1`

	var u model.User
	err := r.db.SQL.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.RoleID, &u.RoleCode, &u.Identifier, &u.Email, &u.FullName, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, q model.UserListQuery) ([]model.User, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, errors.New("database tidak tersedia")
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.identifier ILIKE $%d OR u.full_name ILIKE $%d OR u.email ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	if q.Role != "" {
		conditions = append(conditions, fmt.Sprintf("r.code = $%d", argIdx))
		args = append(args, strings.ToUpper(q.Role))
		argIdx++
	}

	if q.Status != "" {
		conditions = append(conditions, fmt.Sprintf("u.status = $%d", argIdx))
		args = append(args, strings.ToUpper(q.Status))
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Hitung total
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u JOIN roles r ON u.role_id = r.id WHERE %s`, whereClause)
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
		SELECT u.id, u.role_id, r.code, u.identifier, u.email, u.full_name, u.password_hash, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.RoleID, &u.RoleCode, &u.Identifier, &u.Email, &u.FullName, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *model.User) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		INSERT INTO users (id, role_id, identifier, email, full_name, password_hash, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.SQL.QueryRowContext(ctx, query,
		u.RoleID, u.Identifier, u.Email, u.FullName, u.PasswordHash, u.Status,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, u *model.User) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		UPDATE users
		SET identifier = $1, email = $2, full_name = $3, password_hash = $4, role_id = $5, status = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`

	err := r.db.SQL.QueryRowContext(ctx, query,
		u.Identifier, u.Email, u.FullName, u.PasswordHash, u.RoleID, u.Status, u.ID,
	).Scan(&u.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.SQL.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) GetRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `SELECT id, code, name FROM roles WHERE code = $1`
	var role model.Role
	err := r.db.SQL.QueryRowContext(ctx, query, strings.ToUpper(code)).Scan(&role.ID, &role.Code, &role.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PostgresUserRepository) HasFaceTemplate(ctx context.Context, userID string) (bool, error) {
	if r.db == nil || r.db.SQL == nil {
		return false, nil
	}

	query := `SELECT EXISTS(SELECT 1 FROM face_templates WHERE user_id = $1 AND status = 'ACTIVE')`
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

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
	ErrAttendanceNotFound  = errors.New("data presensi tidak ditemukan")
	ErrDuplicateAttendance = errors.New("mahasiswa sudah melakukan presensi pada sesi ini")
	ErrAlreadyCheckedOut   = errors.New("mahasiswa sudah melakukan check-out sebelumnya")
)

type AttendanceRepository interface {
	GetByUserAndSession(ctx context.Context, userID, sessionID string) (*model.Attendance, error)
	GetByID(ctx context.Context, id string) (*model.Attendance, error)
	Create(ctx context.Context, att *model.Attendance) error
	CheckOut(ctx context.Context, id string, checkedOutAt time.Time) error
	List(ctx context.Context, query model.AttendanceListQuery) ([]model.Attendance, int, error)
}

type PostgresAttendanceRepository struct {
	db *DB
}

func NewPostgresAttendanceRepository(db *DB) *PostgresAttendanceRepository {
	return &PostgresAttendanceRepository{db: db}
}

func (r *PostgresAttendanceRepository) GetByUserAndSession(ctx context.Context, userID, sessionID string) (*model.Attendance, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT a.id, a.user_id, a.session_id, a.checked_in_at, a.checked_out_at,
		       a.verification_status, a.created_at, a.updated_at,
		       u.identifier, u.full_name, s.code, s.title, TO_CHAR(s.session_date, 'YYYY-MM-DD')
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		JOIN lab_sessions s ON a.session_id = s.id
		WHERE a.user_id = $1 AND a.session_id = $2
		LIMIT 1`

	var a model.Attendance
	var checkedOutAt sql.NullTime
	err := r.db.SQL.QueryRowContext(ctx, query, userID, sessionID).Scan(
		&a.ID, &a.UserID, &a.SessionID, &a.CheckedInAt, &checkedOutAt,
		&a.VerificationStatus, &a.CreatedAt, &a.UpdatedAt,
		&a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle, &a.SessionDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAttendanceNotFound
	}
	if err != nil {
		return nil, err
	}

	if checkedOutAt.Valid {
		t := checkedOutAt.Time.UTC()
		a.CheckedOutAt = &t
	}
	a.CheckedInAt = a.CheckedInAt.UTC()
	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAttendanceRepository) GetByID(ctx context.Context, id string) (*model.Attendance, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT a.id, a.user_id, a.session_id, a.checked_in_at, a.checked_out_at,
		       a.verification_status, a.created_at, a.updated_at,
		       u.identifier, u.full_name, s.code, s.title, TO_CHAR(s.session_date, 'YYYY-MM-DD')
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		JOIN lab_sessions s ON a.session_id = s.id
		WHERE a.id = $1
		LIMIT 1`

	var a model.Attendance
	var checkedOutAt sql.NullTime
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.UserID, &a.SessionID, &a.CheckedInAt, &checkedOutAt,
		&a.VerificationStatus, &a.CreatedAt, &a.UpdatedAt,
		&a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle, &a.SessionDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAttendanceNotFound
	}
	if err != nil {
		return nil, err
	}

	if checkedOutAt.Valid {
		t := checkedOutAt.Time.UTC()
		a.CheckedOutAt = &t
	}
	a.CheckedInAt = a.CheckedInAt.UTC()
	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return &a, nil
}

func (r *PostgresAttendanceRepository) Create(ctx context.Context, a *model.Attendance) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		INSERT INTO attendance (id, user_id, session_id, checked_in_at, checked_out_at, verification_status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NULL, $4, NOW(), NOW())
		RETURNING id, checked_in_at, created_at, updated_at`

	err := r.db.SQL.QueryRowContext(ctx, query,
		a.UserID, a.SessionID, a.CheckedInAt.UTC(), a.VerificationStatus,
	).Scan(&a.ID, &a.CheckedInAt, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "idx_attendance_user_session") || strings.Contains(err.Error(), "23505") {
			return ErrDuplicateAttendance
		}
		return err
	}

	a.CheckedInAt = a.CheckedInAt.UTC()
	a.CreatedAt = a.CreatedAt.UTC()
	a.UpdatedAt = a.UpdatedAt.UTC()

	return nil
}

func (r *PostgresAttendanceRepository) CheckOut(ctx context.Context, id string, checkedOutAt time.Time) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		UPDATE attendance
		SET checked_out_at = $1, updated_at = NOW()
		WHERE id = $2 AND checked_out_at IS NULL`

	res, err := r.db.SQL.ExecContext(ctx, query, checkedOutAt.UTC(), id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// Cek apakah memang tidak ada atau sudah checkout sebelumnya
		existing, err := r.GetByID(ctx, id)
		if err != nil {
			return ErrAttendanceNotFound
		}
		if existing.CheckedOutAt != nil {
			return ErrAlreadyCheckedOut
		}
		return ErrAttendanceNotFound
	}

	return nil
}

func (r *PostgresAttendanceRepository) List(ctx context.Context, q model.AttendanceListQuery) ([]model.Attendance, int, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, 0, errors.New("database tidak tersedia")
	}

	conditions := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if q.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("a.user_id = $%d", argIdx))
		args = append(args, q.UserID)
		argIdx++
	}

	if q.SessionID != "" {
		conditions = append(conditions, fmt.Sprintf("a.session_id = $%d", argIdx))
		args = append(args, q.SessionID)
		argIdx++
	}

	if q.Date != "" {
		conditions = append(conditions, fmt.Sprintf("s.session_date = $%d::date", argIdx))
		args = append(args, q.Date)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM attendance a
		JOIN lab_sessions s ON a.session_id = s.id
		WHERE %s`, whereClause)

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
		SELECT a.id, a.user_id, a.session_id, a.checked_in_at, a.checked_out_at,
		       a.verification_status, a.created_at, a.updated_at,
		       u.identifier, u.full_name, s.code, s.title, TO_CHAR(s.session_date, 'YYYY-MM-DD')
		FROM attendance a
		JOIN users u ON a.user_id = u.id
		JOIN lab_sessions s ON a.session_id = s.id
		WHERE %s
		ORDER BY a.checked_in_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []model.Attendance{}
	for rows.Next() {
		var a model.Attendance
		var checkedOutAt sql.NullTime
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.SessionID, &a.CheckedInAt, &checkedOutAt,
			&a.VerificationStatus, &a.CreatedAt, &a.UpdatedAt,
			&a.UserIdentifier, &a.UserFullName, &a.SessionCode, &a.SessionTitle, &a.SessionDate,
		); err != nil {
			return nil, 0, err
		}
		if checkedOutAt.Valid {
			t := checkedOutAt.Time.UTC()
			a.CheckedOutAt = &t
		}
		a.CheckedInAt = a.CheckedInAt.UTC()
		a.CreatedAt = a.CreatedAt.UTC()
		a.UpdatedAt = a.UpdatedAt.UTC()

		list = append(list, a)
	}

	return list, total, nil
}

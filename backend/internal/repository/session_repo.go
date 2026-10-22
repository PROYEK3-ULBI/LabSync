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
	ErrSessionNotFound      = errors.New("sesi praktikum tidak ditemukan")
	ErrDuplicateSessionCode = errors.New("kode sesi sudah digunakan")
	ErrSessionHasAttendance = errors.New("sesi praktikum telah memiliki data presensi dan tidak dapat dihapus secara fisik")
)

type SessionRepository interface {
	GetByID(ctx context.Context, id string) (*model.LabSession, error)
	GetByCode(ctx context.Context, code string) (*model.LabSession, error)
	List(ctx context.Context, query model.SessionListQuery) ([]model.LabSession, int, error)
	GetActive(ctx context.Context) ([]model.LabSession, error)
	Create(ctx context.Context, session *model.LabSession) error
	Update(ctx context.Context, session *model.LabSession) error
	UpdateStatus(ctx context.Context, id string, status model.SessionStatus) error
	HasAttendance(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, id string) error
}

type PostgresSessionRepository struct {
	db *DB
}

func NewPostgresSessionRepository(db *DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) GetByID(ctx context.Context, id string) (*model.LabSession, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, code, title, course_code, location, TO_CHAR(session_date, 'YYYY-MM-DD'),
		       start_at, end_at, status, created_by, created_at, updated_at
		FROM lab_sessions
		WHERE id = $1`

	var s model.LabSession
	var courseCode, createdBy sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Code, &s.Title, &courseCode, &s.Location, &s.SessionDate,
		&s.StartAt, &s.EndAt, &s.Status, &createdBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	if courseCode.Valid {
		s.CourseCode = &courseCode.String
	}
	if createdBy.Valid {
		s.CreatedBy = &createdBy.String
	}

	s.StartAt = s.StartAt.UTC()
	s.EndAt = s.EndAt.UTC()
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()

	return &s, nil
}

func (r *PostgresSessionRepository) GetByCode(ctx context.Context, code string) (*model.LabSession, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, code, title, course_code, location, TO_CHAR(session_date, 'YYYY-MM-DD'),
		       start_at, end_at, status, created_by, created_at, updated_at
		FROM lab_sessions
		WHERE code = $1`

	var s model.LabSession
	var courseCode, createdBy sql.NullString
	err := r.db.SQL.QueryRowContext(ctx, query, code).Scan(
		&s.ID, &s.Code, &s.Title, &courseCode, &s.Location, &s.SessionDate,
		&s.StartAt, &s.EndAt, &s.Status, &createdBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	if courseCode.Valid {
		s.CourseCode = &courseCode.String
	}
	if createdBy.Valid {
		s.CreatedBy = &createdBy.String
	}

	s.StartAt = s.StartAt.UTC()
	s.EndAt = s.EndAt.UTC()
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()

	return &s, nil
}

func (r *PostgresSessionRepository) List(ctx context.Context, q model.SessionListQuery) ([]model.LabSession, int, error) {
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

	if q.SessionDate != "" {
		conditions = append(conditions, fmt.Sprintf("session_date = $%d::date", argIdx))
		args = append(args, q.SessionDate)
		argIdx++
	}

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(code ILIKE $%d OR title ILIKE $%d OR location ILIKE $%d OR COALESCE(course_code, '') ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM lab_sessions WHERE %s`, whereClause)
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
		SELECT id, code, title, course_code, location, TO_CHAR(session_date, 'YYYY-MM-DD'),
		       start_at, end_at, status, created_by, created_at, updated_at
		FROM lab_sessions
		WHERE %s
		ORDER BY session_date DESC, start_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, q.Limit, offset)

	rows, err := r.db.SQL.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	sessions := []model.LabSession{}
	for rows.Next() {
		var s model.LabSession
		var courseCode, createdBy sql.NullString
		if err := rows.Scan(
			&s.ID, &s.Code, &s.Title, &courseCode, &s.Location, &s.SessionDate,
			&s.StartAt, &s.EndAt, &s.Status, &createdBy, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if courseCode.Valid {
			s.CourseCode = &courseCode.String
		}
		if createdBy.Valid {
			s.CreatedBy = &createdBy.String
		}
		s.StartAt = s.StartAt.UTC()
		s.EndAt = s.EndAt.UTC()
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()

		sessions = append(sessions, s)
	}

	return sessions, total, nil
}

func (r *PostgresSessionRepository) GetActive(ctx context.Context) ([]model.LabSession, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, code, title, course_code, location, TO_CHAR(session_date, 'YYYY-MM-DD'),
		       start_at, end_at, status, created_by, created_at, updated_at
		FROM lab_sessions
		WHERE status = 'ACTIVE'
		ORDER BY start_at ASC`

	rows, err := r.db.SQL.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []model.LabSession{}
	for rows.Next() {
		var s model.LabSession
		var courseCode, createdBy sql.NullString
		if err := rows.Scan(
			&s.ID, &s.Code, &s.Title, &courseCode, &s.Location, &s.SessionDate,
			&s.StartAt, &s.EndAt, &s.Status, &createdBy, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if courseCode.Valid {
			s.CourseCode = &courseCode.String
		}
		if createdBy.Valid {
			s.CreatedBy = &createdBy.String
		}
		s.StartAt = s.StartAt.UTC()
		s.EndAt = s.EndAt.UTC()
		s.CreatedAt = s.CreatedAt.UTC()
		s.UpdatedAt = s.UpdatedAt.UTC()

		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (r *PostgresSessionRepository) Create(ctx context.Context, s *model.LabSession) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		INSERT INTO lab_sessions (id, code, title, course_code, location, session_date, start_at, end_at, status, created_by, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5::date, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	var courseCode, createdBy sql.NullString
	if s.CourseCode != nil {
		courseCode = sql.NullString{String: *s.CourseCode, Valid: true}
	}
	if s.CreatedBy != nil {
		createdBy = sql.NullString{String: *s.CreatedBy, Valid: true}
	}

	err := r.db.SQL.QueryRowContext(ctx, query,
		s.Code, s.Title, courseCode, s.Location, s.SessionDate,
		s.StartAt.UTC(), s.EndAt.UTC(), string(s.Status), createdBy,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			return ErrDuplicateSessionCode
		}
		return err
	}
	return nil
}

func (r *PostgresSessionRepository) Update(ctx context.Context, s *model.LabSession) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		UPDATE lab_sessions
		SET code = $1, title = $2, course_code = $3, location = $4,
		    session_date = $5::date, start_at = $6, end_at = $7, status = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at`

	var courseCode sql.NullString
	if s.CourseCode != nil {
		courseCode = sql.NullString{String: *s.CourseCode, Valid: true}
	}

	err := r.db.SQL.QueryRowContext(ctx, query,
		s.Code, s.Title, courseCode, s.Location,
		s.SessionDate, s.StartAt.UTC(), s.EndAt.UTC(), string(s.Status), s.ID,
	).Scan(&s.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSessionNotFound
		}
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			return ErrDuplicateSessionCode
		}
		return err
	}
	return nil
}

func (r *PostgresSessionRepository) UpdateStatus(ctx context.Context, id string, status model.SessionStatus) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `UPDATE lab_sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.SQL.ExecContext(ctx, query, string(status), id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (r *PostgresSessionRepository) HasAttendance(ctx context.Context, sessionID string) (bool, error) {
	if r.db == nil || r.db.SQL == nil {
		return false, nil
	}

	query := `SELECT EXISTS(SELECT 1 FROM attendance WHERE session_id = $1)`
	var exists bool
	err := r.db.SQL.QueryRowContext(ctx, query, sessionID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresSessionRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	hasAtt, err := r.HasAttendance(ctx, id)
	if err != nil {
		return err
	}
	if hasAtt {
		return ErrSessionHasAttendance
	}

	query := `DELETE FROM lab_sessions WHERE id = $1`
	res, err := r.db.SQL.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

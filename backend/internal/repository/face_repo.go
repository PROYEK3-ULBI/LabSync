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
	ErrFaceTemplateNotFound = errors.New("template biometrik wajah tidak ditemukan")
	ErrFaceAlreadyEnrolled  = errors.New("pengguna sudah memiliki pendaftaran wajah aktif")
)

type FaceRepository interface {
	GetActiveTemplateByUserID(ctx context.Context, userID string) (*model.FaceTemplate, []model.FaceTemplateSample, error)
	SaveEnrollment(ctx context.Context, template *model.FaceTemplate, samples []model.FaceTemplateSample, revokeOld bool) error
	RevokeTemplate(ctx context.Context, userID string) error
}

type PostgresFaceRepository struct {
	db *DB
}

func NewPostgresFaceRepository(db *DB) *PostgresFaceRepository {
	return &PostgresFaceRepository{db: db}
}

func (r *PostgresFaceRepository) GetActiveTemplateByUserID(ctx context.Context, userID string) (*model.FaceTemplate, []model.FaceTemplateSample, error) {
	if r.db == nil || r.db.SQL == nil {
		return nil, nil, errors.New("database tidak tersedia")
	}

	query := `
		SELECT id, user_id, model_name, model_version, consent_version, consented_at, status, enrolled_at, updated_at
		FROM face_templates
		WHERE user_id = $1 AND status = 'ACTIVE'
		LIMIT 1`

	var t model.FaceTemplate
	err := r.db.SQL.QueryRowContext(ctx, query, userID).Scan(
		&t.ID, &t.UserID, &t.ModelName, &t.ModelVersion, &t.ConsentVersion,
		&t.ConsentedAt, &t.Status, &t.EnrolledAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrFaceTemplateNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	// Ambil sampel-sampel yang terasosiasi
	sampleQuery := `
		SELECT id, face_template_id, pose, embedding_ciphertext, quality_status, created_at
		FROM face_template_samples
		WHERE face_template_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.SQL.QueryContext(ctx, sampleQuery, t.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var samples []model.FaceTemplateSample
	for rows.Next() {
		var s model.FaceTemplateSample
		if err := rows.Scan(
			&s.ID, &s.FaceTemplateID, &s.Pose, &s.EmbeddingCiphertext, &s.QualityStatus, &s.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		samples = append(samples, s)
	}

	return &t, samples, nil
}

func (r *PostgresFaceRepository) SaveEnrollment(ctx context.Context, template *model.FaceTemplate, samples []model.FaceTemplateSample, revokeOld bool) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	tx, err := r.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if revokeOld {
		revokeQuery := `
			UPDATE face_templates
			SET status = 'REVOKED', updated_at = NOW()
			WHERE user_id = $1 AND status = 'ACTIVE'`
		if _, err := tx.ExecContext(ctx, revokeQuery, template.UserID); err != nil {
			return fmt.Errorf("gagal mencabut template lama: %w", err)
		}
	}

	insertTemplateQuery := `
		INSERT INTO face_templates (id, user_id, model_name, model_version, consent_version, consented_at, status, enrolled_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'ACTIVE', NOW(), NOW())
		RETURNING id, enrolled_at, updated_at`

	err = tx.QueryRowContext(ctx, insertTemplateQuery,
		template.UserID, template.ModelName, template.ModelVersion, template.ConsentVersion, template.ConsentedAt,
	).Scan(&template.ID, &template.EnrolledAt, &template.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "idx_face_templates_active_user") || strings.Contains(err.Error(), "23505") {
			return ErrFaceAlreadyEnrolled
		}
		return fmt.Errorf("gagal menyimpan metadata template wajah: %w", err)
	}

	insertSampleQuery := `
		INSERT INTO face_template_samples (id, face_template_id, pose, embedding_ciphertext, quality_status, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())
		RETURNING id, created_at`

	for i := range samples {
		samples[i].FaceTemplateID = template.ID
		err = tx.QueryRowContext(ctx, insertSampleQuery,
			template.ID, samples[i].Pose, samples[i].EmbeddingCiphertext, samples[i].QualityStatus,
		).Scan(&samples[i].ID, &samples[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("gagal menyimpan sampel pose %s: %w", samples[i].Pose, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi pendaftaran wajah: %w", err)
	}

	return nil
}

func (r *PostgresFaceRepository) RevokeTemplate(ctx context.Context, userID string) error {
	if r.db == nil || r.db.SQL == nil {
		return errors.New("database tidak tersedia")
	}

	query := `
		UPDATE face_templates
		SET status = 'REVOKED', updated_at = NOW()
		WHERE user_id = $1 AND status = 'ACTIVE'`

	res, err := r.db.SQL.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrFaceTemplateNotFound
	}
	return nil
}

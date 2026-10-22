package repository

import "context"

type HealthRepository struct {
	db *DB
}

func NewHealthRepository(db *DB) *HealthRepository {
	return &HealthRepository{db: db}
}

// CheckDB mengecek koneksi database
func (r *HealthRepository) CheckDB(ctx context.Context) error {
	if r.db == nil {
		return nil
	}
	return r.db.Ping(ctx)
}

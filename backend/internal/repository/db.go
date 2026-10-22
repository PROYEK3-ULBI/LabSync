package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

// NewDB menginisialisasi connection pool ke database PostgreSQL
func NewDB(databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return &DB{SQL: nil}, nil
	}

	connConfig, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("konfigurasi database tidak valid: %w", err)
	}
	// Supabase transaction pooler tidak aman dengan prepared statement cache pgx default.
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*connConfig)

	// Konfigurasi connection pooling standar
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		// Mengembalikan instance agar aplikasi tetap dapat boot (status degraded), disertai error
		return &DB{SQL: db}, fmt.Errorf("ping database gagal: %w", err)
	}

	return &DB{SQL: db}, nil
}

// Close menutup seluruh koneksi pada pool
func (d *DB) Close() error {
	if d.SQL != nil {
		return d.SQL.Close()
	}
	return nil
}

// Ping memeriksa konektivitas aktif ke database
func (d *DB) Ping(ctx context.Context) error {
	if d.SQL == nil {
		return fmt.Errorf("koneksi database belum terkonfigurasi")
	}
	return d.SQL.PingContext(ctx)
}

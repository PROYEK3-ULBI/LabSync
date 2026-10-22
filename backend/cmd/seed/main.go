package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"labsync-backend/internal/config"
	"labsync-backend/internal/service"
)

const (
	adminIdentifier    = "ADM-DEMO-001"
	studentIdentifier  = "MHS-DEMO-001"
	defaultAdminPass   = "AdminDemo123!"
	defaultStudentPass = "StudentDemo123!"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi tidak valid: %v", err)
	}
	if cfg.AppEnv == "production" {
		log.Fatal("seed development ditolak pada APP_ENV=production")
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}

	adminPassword := os.Getenv("DEV_SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = defaultAdminPass
	}
	studentPassword := os.Getenv("DEV_SEED_STUDENT_PASSWORD")
	if studentPassword == "" {
		studentPassword = defaultStudentPass
	}

	adminHash, err := service.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("gagal hash akun admin: %v", err)
	}
	studentHash, err := service.HashPassword(studentPassword)
	if err != nil {
		log.Fatalf("gagal hash akun mahasiswa: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal membuat pool database: %v", err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatalf("gagal mulai transaksi seed: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := seed(ctx, tx, adminHash, studentHash); err != nil {
		log.Fatalf("seed gagal: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit seed gagal: %v", err)
	}

	fmt.Println("Seed development selesai.")
	fmt.Printf("Admin login: identifier=%s password=%s\n", adminIdentifier, adminPassword)
	fmt.Printf("Mahasiswa login: identifier=%s password=%s\n", studentIdentifier, studentPassword)
	fmt.Println("Simpan kredensial ini lokal saja. Jangan commit ke repo.")
}

func seed(ctx context.Context, tx pgx.Tx, adminHash, studentHash string) error {
	queries := []struct {
		name string
		sql  string
	}{
		{
			name: "users",
			sql: `
INSERT INTO users (role_id, identifier, email, full_name, password_hash, status)
VALUES
  ((SELECT id FROM roles WHERE code = 'ADMIN'), 'ADM-DEMO-001', 'admin.demo@labsync.local', 'Admin Demo LabSync', $1, 'ACTIVE'),
  ((SELECT id FROM roles WHERE code = 'STUDENT'), 'MHS-DEMO-001', 'mahasiswa.demo@labsync.local', 'Mahasiswa Demo LabSync', $2, 'ACTIVE')
ON CONFLICT (identifier) DO UPDATE SET
  email = EXCLUDED.email,
  full_name = EXCLUDED.full_name,
  password_hash = EXCLUDED.password_hash,
  status = EXCLUDED.status,
  role_id = EXCLUDED.role_id;`,
		},
		{
			name: "sessions",
			sql: `
INSERT INTO lab_sessions (code, title, course_code, location, session_date, start_at, end_at, status, created_by)
VALUES
  ('PWB-DEMO-01', 'Praktikum Pemrograman Web Demo', 'IF2101', 'Lab Komputer 1', CURRENT_DATE, NOW() - INTERVAL '30 minutes', NOW() + INTERVAL '2 hours', 'ACTIVE', (SELECT id FROM users WHERE identifier = 'ADM-DEMO-001')),
  ('BDT-DEMO-01', 'Praktikum Basis Data Demo', 'IF2202', 'Lab Komputer 2', CURRENT_DATE + INTERVAL '1 day', NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day 2 hours', 'SCHEDULED', (SELECT id FROM users WHERE identifier = 'ADM-DEMO-001'))
ON CONFLICT (code) DO UPDATE SET
  title = EXCLUDED.title,
  course_code = EXCLUDED.course_code,
  location = EXCLUDED.location,
  session_date = EXCLUDED.session_date,
  start_at = EXCLUDED.start_at,
  end_at = EXCLUDED.end_at,
  status = EXCLUDED.status,
  created_by = EXCLUDED.created_by;`,
		},
		{
			name: "assets",
			sql: `
INSERT INTO assets (asset_code, asset_name, category, brand, model, serial_number, location, status, notes)
VALUES
  ('LAB-DEMO-LT-001', 'Laptop Demo Dell Latitude', 'Laptop', 'Dell', 'Latitude 5490', 'DEMO-LT-001', 'Lab Komputer 1', 'AVAILABLE', 'Data seed development'),
  ('LAB-DEMO-PC-001', 'PC Demo HP ProDesk', 'Desktop', 'HP', 'ProDesk 400', 'DEMO-PC-001', 'Lab Komputer 2', 'AVAILABLE', 'Data seed development'),
  ('LAB-DEMO-SW-001', 'Switch Demo Cisco', 'Perangkat Jaringan', 'Cisco', 'Catalyst 2960', 'DEMO-SW-001', 'Lab Komputer 1', 'MAINTENANCE', 'Data seed development')
ON CONFLICT (asset_code) DO UPDATE SET
  asset_name = EXCLUDED.asset_name,
  category = EXCLUDED.category,
  brand = EXCLUDED.brand,
  model = EXCLUDED.model,
  serial_number = EXCLUDED.serial_number,
  location = EXCLUDED.location,
  status = EXCLUDED.status,
  notes = EXCLUDED.notes;`,
		},
		{
			name: "maintenance",
			sql: `
INSERT INTO maintenance (asset_id, reported_by, issue_type, description, priority, status)
SELECT a.id, u.id, 'HARDWARE', 'Perangkat demo membutuhkan pemeriksaan', 'MEDIUM', 'OPEN'
FROM assets a
CROSS JOIN users u
WHERE a.asset_code = 'LAB-DEMO-SW-001' AND u.identifier = 'ADM-DEMO-001'
AND NOT EXISTS (
  SELECT 1 FROM maintenance m WHERE m.asset_id = a.id AND m.status IN ('OPEN', 'IN_PROGRESS')
);`,
		},
	}

	for _, q := range queries {
		args := []any{}
		if q.name == "users" {
			args = []any{adminHash, studentHash}
		}
		if _, err := tx.Exec(ctx, q.sql, args...); err != nil {
			return fmt.Errorf("seed %s: %w", q.name, err)
		}
	}
	return nil
}

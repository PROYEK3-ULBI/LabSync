-- Migration: Init Schema (UP)

-- ==========================================
-- 1. ROLES
-- ==========================================
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);

-- ==========================================
-- 2. USERS
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    identifier VARCHAR(50) NOT NULL UNIQUE, -- NIM/NIP
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(150) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Seed initial roles
INSERT INTO roles (code, name) VALUES
    ('ADMIN', 'Administrator'),
    ('STUDENT', 'Mahasiswa')
ON CONFLICT (code) DO NOTHING;

-- ==========================================
-- 3. FACE_TEMPLATES
-- ==========================================
CREATE TABLE face_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    model_name VARCHAR(100) NOT NULL,
    model_version VARCHAR(50) NOT NULL,
    consent_version VARCHAR(50) NOT NULL,
    consented_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'REVOKED')),
    enrolled_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Constraint: Maksimal 1 enrollment aktif per user
CREATE UNIQUE INDEX idx_face_templates_active_user ON face_templates(user_id) WHERE status = 'ACTIVE';

-- ==========================================
-- 4. FACE_TEMPLATE_SAMPLES
-- ==========================================
CREATE TABLE face_template_samples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    face_template_id UUID NOT NULL REFERENCES face_templates(id) ON DELETE CASCADE,
    pose VARCHAR(20) NOT NULL CHECK (pose IN ('FRONT', 'LEFT', 'RIGHT')),
    embedding_ciphertext BYTEA NOT NULL,
    quality_status VARCHAR(20) NOT NULL DEFAULT 'VALID',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Constraint: Dalam 1 enrollment set, tidak boleh ada pose yang ganda
CREATE UNIQUE INDEX idx_face_samples_pose ON face_template_samples(face_template_id, pose);

-- ==========================================
-- 5. LAB_SESSIONS
-- ==========================================
CREATE TABLE lab_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    course_code VARCHAR(50),
    location VARCHAR(150) NOT NULL,
    session_date DATE NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED' CHECK (status IN ('SCHEDULED', 'ACTIVE', 'CLOSED', 'CANCELLED')),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_session_time CHECK (start_at < end_at)
);

-- ==========================================
-- 6. ATTENDANCE
-- ==========================================
CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    session_id UUID NOT NULL REFERENCES lab_sessions(id) ON DELETE RESTRICT,
    checked_in_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    checked_out_at TIMESTAMP WITH TIME ZONE,
    verification_status VARCHAR(50) NOT NULL DEFAULT 'VERIFIED_FACE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_checkout_time CHECK (checked_out_at IS NULL OR checked_out_at >= checked_in_at)
);

-- Constraint: 1 user hanya 1 presensi per sesi
CREATE UNIQUE INDEX idx_attendance_user_session ON attendance(user_id, session_id);

-- ==========================================
-- 7. ASSETS
-- ==========================================
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_code VARCHAR(100) NOT NULL UNIQUE,
    asset_name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'Umum',
    brand VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(100) UNIQUE,
    location VARCHAR(150) NOT NULL DEFAULT 'Lab Komputer',
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE', 'ALLOCATED', 'MAINTENANCE', 'RETIRED')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 8. ASSET_ALLOCATIONS
-- ==========================================
CREATE TABLE asset_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    session_id UUID NOT NULL REFERENCES lab_sessions(id) ON DELETE RESTRICT,
    allocated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    returned_at TIMESTAMP WITH TIME ZONE,
    condition_before VARCHAR(30) DEFAULT 'GOOD',
    condition_after VARCHAR(30),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_return_time CHECK (returned_at IS NULL OR returned_at >= allocated_at)
);

-- Constraint: Mencegah alokasi ganda untuk aset yang belum dikembalikan
CREATE UNIQUE INDEX idx_asset_allocations_active ON asset_allocations(asset_id) WHERE returned_at IS NULL;

-- ==========================================
-- 9. MAINTENANCE
-- ==========================================
CREATE TABLE maintenance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    reported_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    issue_type VARCHAR(50) NOT NULL DEFAULT 'HARDWARE',
    description TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM' CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH')),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED')),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    resolution TEXT,
    cost NUMERIC(14,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 10. AUDIT_LOGS
-- ==========================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Null jika user dihapus (walau seharusnya disoft-delete)
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'SUCCESS',
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ==========================================
-- TRIGGER UNTUK UPDATED_AT (Opsional untuk testing lokal)
-- ==========================================
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_face_templates
BEFORE UPDATE ON face_templates
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_lab_sessions
BEFORE UPDATE ON lab_sessions
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_assets
BEFORE UPDATE ON assets
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_asset_allocations
BEFORE UPDATE ON asset_allocations
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_attendance
BEFORE UPDATE ON attendance
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_maintenance
BEFORE UPDATE ON maintenance
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

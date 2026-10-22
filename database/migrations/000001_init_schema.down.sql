-- Migration: Init Schema (DOWN)

-- Drop triggers
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS set_timestamp_face_templates ON face_templates;
DROP TRIGGER IF EXISTS set_timestamp_lab_sessions ON lab_sessions;
DROP TRIGGER IF EXISTS set_timestamp_assets ON assets;
DROP TRIGGER IF EXISTS set_timestamp_asset_allocations ON asset_allocations;
DROP TRIGGER IF EXISTS set_timestamp_attendance ON attendance;
DROP TRIGGER IF EXISTS set_timestamp_maintenance ON maintenance;

-- Drop function
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS maintenance CASCADE;
DROP TABLE IF EXISTS asset_allocations CASCADE;
DROP TABLE IF EXISTS assets CASCADE;
DROP TABLE IF EXISTS attendance CASCADE;
DROP TABLE IF EXISTS lab_sessions CASCADE;
DROP TABLE IF EXISTS face_template_samples CASCADE;
DROP TABLE IF EXISTS face_templates CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

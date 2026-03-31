-- Patch Migration: Fix uhids table timestamp columns
-- Run this on existing databases where 001 was already applied
-- The UHID model uses int64 (BIGINT milliseconds) for issued_date, created_at, updated_at

ALTER TABLE uhids
  ALTER COLUMN issued_date TYPE BIGINT USING EXTRACT(EPOCH FROM issued_date)::BIGINT * 1000;

ALTER TABLE uhids
  ALTER COLUMN created_at TYPE BIGINT USING EXTRACT(EPOCH FROM created_at)::BIGINT * 1000;

ALTER TABLE uhids
  ALTER COLUMN updated_at TYPE BIGINT USING EXTRACT(EPOCH FROM updated_at)::BIGINT * 1000;

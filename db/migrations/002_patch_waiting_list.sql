-- Patch Migration: Fix waiting_list_entries table
-- Run this on existing databases where 001 was already applied

-- Add missing notes column
ALTER TABLE waiting_list_entries ADD COLUMN IF NOT EXISTS notes TEXT;

-- Fix created_at / updated_at to BIGINT (model uses int64 milliseconds)
ALTER TABLE waiting_list_entries
  ALTER COLUMN created_at TYPE BIGINT USING EXTRACT(EPOCH FROM created_at)::BIGINT * 1000;

ALTER TABLE waiting_list_entries
  ALTER COLUMN updated_at TYPE BIGINT USING EXTRACT(EPOCH FROM updated_at)::BIGINT * 1000;

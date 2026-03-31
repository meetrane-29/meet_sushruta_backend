-- ============================================================
-- Patch: Make admission_id Nullable in progress_notes
-- Purpose: Support SOAP notes for OPD consultations (no admission)
-- Since: SOAP notes should work for both IPD (with admission) and OPD (without admission)
-- ============================================================

-- Drop the foreign key constraint first
ALTER TABLE progress_notes 
DROP CONSTRAINT IF EXISTS fk_progress_notes_admission;

-- Make admission_id nullable
ALTER TABLE progress_notes 
ALTER COLUMN admission_id DROP NOT NULL;

-- Recreate the foreign key constraint (now with ON DELETE SET NULL)
ALTER TABLE progress_notes
ADD CONSTRAINT fk_progress_notes_admission 
FOREIGN KEY(admission_id) REFERENCES admission_records(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- Add missing columns if they don't exist
ALTER TABLE progress_notes 
ADD COLUMN IF NOT EXISTS recorded_date VARCHAR(100),
ADD COLUMN IF NOT EXISTS next_review_date VARCHAR(100),
ADD COLUMN IF NOT EXISTS medications TEXT;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_progress_notes_recorded_date ON progress_notes(recorded_date);
CREATE INDEX IF NOT EXISTS idx_progress_notes_patient_id ON progress_notes(patient_id);

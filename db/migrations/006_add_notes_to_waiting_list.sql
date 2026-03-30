-- Add missing notes column to waiting_list_entries table
-- Created: March 30, 2026

ALTER TABLE IF EXISTS waiting_list_entries
ADD COLUMN IF NOT EXISTS notes TEXT;

-- Migration completed successfully
-- This adds the notes column that the model expects

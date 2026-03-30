-- ============================================================
-- meet_sushruta - Patient Rating System Tables
-- Created for Doctor Module - Patient Ratings Feature
-- Run this migration to enable patient ratings features
-- ============================================================

-- ============================================================
-- 1. RATINGS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS ratings (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at          TIMESTAMPTZ,

    doctor_id           UUID NOT NULL,
    patient_id          UUID NOT NULL,
    
    rating              INTEGER NOT NULL DEFAULT 5,           -- Overall rating (1-5)
    professionalism     INTEGER NOT NULL DEFAULT 5,           -- Doctor's professionalism (1-5)
    communication       INTEGER NOT NULL DEFAULT 5,           -- Communication clarity (1-5)
    punctuality         INTEGER NOT NULL DEFAULT 5,           -- Doctor's punctuality (1-5)
    cleanliness         INTEGER NOT NULL DEFAULT 5,           -- Clinic/facility cleanliness (1-5)
    comment             TEXT,                                  -- Patient's feedback/comments

    CONSTRAINT fk_ratings_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_ratings_patient FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_rating_range CHECK (rating >= 1 AND rating <= 5),
    CONSTRAINT chk_professionalism_range CHECK (professionalism >= 1 AND professionalism <= 5),
    CONSTRAINT chk_communication_range CHECK (communication >= 1 AND communication <= 5),
    CONSTRAINT chk_punctuality_range CHECK (punctuality >= 1 AND punctuality <= 5),
    CONSTRAINT chk_cleanliness_range CHECK (cleanliness >= 1 AND cleanliness <= 5)
);

-- Indexes for optimal query performance
CREATE INDEX IF NOT EXISTS idx_ratings_doctor_id ON ratings(doctor_id);
CREATE INDEX IF NOT EXISTS idx_ratings_patient_id ON ratings(patient_id);
CREATE INDEX IF NOT EXISTS idx_ratings_created_at ON ratings(created_at);
CREATE INDEX IF NOT EXISTS idx_ratings_deleted_at ON ratings(deleted_at);
CREATE INDEX IF NOT EXISTS idx_ratings_doctor_created ON ratings(doctor_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ratings_patient_created ON ratings(patient_id, created_at DESC);

-- ============================================================
-- END OF MIGRATION
-- ============================================================

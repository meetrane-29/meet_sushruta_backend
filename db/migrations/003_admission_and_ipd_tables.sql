-- ============================================================
-- meet_sushruta - Admission & IPD Management Tables
-- Created for Doctor Module Phase 1
-- Run this migration to enable IPD/admission features
-- ============================================================

-- ============================================================
-- 1. ADMISSION RECORDS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS admission_records (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at          TIMESTAMPTZ,

    patient_id          UUID NOT NULL,
    doctor_id           UUID NOT NULL,
    bed_id              UUID,
    
    admission_date      VARCHAR(100) NOT NULL,  -- YYYY-MM-DD HH:MM
    discharge_date      VARCHAR(100),           -- YYYY-MM-DD HH:MM (null if active)
    reason              TEXT,                   -- Reason for admission
    diagnosis           TEXT,                   -- Initial diagnosis
    status              VARCHAR(50) NOT NULL DEFAULT 'active',  -- active | discharged | cancelled
    ward                VARCHAR(100),           -- ICU | General | etc
    room_number         VARCHAR(50),
    is_emergency        BOOLEAN DEFAULT FALSE,
    notes               TEXT,

    CONSTRAINT fk_admission_patient FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_admission_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_admission_bed FOREIGN KEY(bed_id) REFERENCES beds(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_admission_records_patient_id ON admission_records(patient_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_doctor_id ON admission_records(doctor_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_bed_id ON admission_records(bed_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_status ON admission_records(status);
CREATE INDEX IF NOT EXISTS idx_admission_records_admission_date ON admission_records(admission_date);
CREATE INDEX IF NOT EXISTS idx_admission_records_deleted_at ON admission_records(deleted_at);

-- ============================================================
-- 2. PROGRESS NOTES TABLE (SOAP Format)
-- ============================================================
CREATE TABLE IF NOT EXISTS progress_notes (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at          TIMESTAMPTZ,

    admission_id        UUID NOT NULL,
    patient_id          UUID NOT NULL,
    doctor_id           UUID NOT NULL,
    
    -- SOAP Format
    subjective          TEXT,           -- Patient's complaints and history
    objective           TEXT,           -- Physical examination findings
    assessment          TEXT,           -- Doctor's analysis and diagnosis
    plan                TEXT,           -- Treatment plan and instructions
    
    vitals              TEXT,           -- JSON: HR, BP, Temp, RR, SpO2
    notes               TEXT,           -- Additional clinical notes
    signature_text      VARCHAR(255),   -- Doctor's signature/approval

    CONSTRAINT fk_progress_notes_admission FOREIGN KEY(admission_id) REFERENCES admission_records(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_progress_notes_patient FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_progress_notes_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_progress_notes_admission_id ON progress_notes(admission_id);
CREATE INDEX IF NOT EXISTS idx_progress_notes_doctor_id ON progress_notes(doctor_id);
CREATE INDEX IF NOT EXISTS idx_progress_notes_created_at ON progress_notes(created_at);
CREATE INDEX IF NOT EXISTS idx_progress_notes_deleted_at ON progress_notes(deleted_at);

-- ============================================================
-- 3. NURSE INSTRUCTIONS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS nurse_instructions (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at              TIMESTAMPTZ,

    admission_id            UUID NOT NULL,
    doctor_id               UUID NOT NULL,
    
    vitals_frequency        VARCHAR(100),           -- Every 4 hours | Every 6 hours | etc
    dietary_restrictions    TEXT,                   -- Diet instructions
    medication_instructions TEXT,                   -- How to administer medications
    activity_level          VARCHAR(100),           -- Bed rest | Restricted mobility | etc
    special_monitoring      TEXT,                   -- Special observation points
    additional_notes        TEXT,                   -- General instructions
    
    created_by_doctor_id    UUID NOT NULL,

    CONSTRAINT fk_nurse_instr_admission FOREIGN KEY(admission_id) REFERENCES admission_records(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_nurse_instr_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_nurse_instr_creator FOREIGN KEY(created_by_doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_nurse_instructions_admission_id ON nurse_instructions(admission_id);
CREATE INDEX IF NOT EXISTS idx_nurse_instructions_doctor_id ON nurse_instructions(doctor_id);
CREATE INDEX IF NOT EXISTS idx_nurse_instructions_deleted_at ON nurse_instructions(deleted_at);

-- ============================================================
-- 4. DISCHARGE SUMMARIES TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS discharge_summaries (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at              TIMESTAMPTZ,

    admission_id            UUID NOT NULL UNIQUE,
    patient_id              UUID NOT NULL,
    doctor_id               UUID NOT NULL,
    
    final_diagnosis         TEXT,                   -- Final confirmed diagnosis
    procedures_done         TEXT,                   -- Procedures/surgeries performed
    discharge_medications   TEXT,                   -- Medications to continue
    follow_up_instructions  TEXT,                   -- Follow-up care instructions
    diet_recommendations    TEXT,                   -- Dietary advice
    activity_restrictions   TEXT,                   -- Activity limitations
    warning_symptoms        TEXT,                   -- Symptoms requiring emergency care
    patient_outcome         VARCHAR(100),           -- Cured | Improved | Stable | etc
    discharge_date          VARCHAR(100),           -- YYYY-MM-DD HH:MM
    discharge_notes         TEXT,                   -- Additional discharge notes

    CONSTRAINT fk_discharge_admission FOREIGN KEY(admission_id) REFERENCES admission_records(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_discharge_patient FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_discharge_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discharge_summaries_admission_id ON discharge_summaries(admission_id);
CREATE INDEX IF NOT EXISTS idx_discharge_summaries_patient_id ON discharge_summaries(patient_id);
CREATE INDEX IF NOT EXISTS idx_discharge_summaries_doctor_id ON discharge_summaries(doctor_id);
CREATE INDEX IF NOT EXISTS idx_discharge_summaries_deleted_at ON discharge_summaries(deleted_at);

-- ============================================================
-- 5. DOCTOR SCHEDULE TABLE (Operation Theatre / OT)
-- ============================================================
CREATE TABLE IF NOT EXISTS doctor_schedules (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    updated_at              BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
    deleted_at              TIMESTAMPTZ,

    doctor_id               UUID NOT NULL,
    operation_theatre_id    UUID,           -- Which OT is assigned
    
    scheduled_date          VARCHAR(100) NOT NULL,  -- YYYY-MM-DD
    start_time              VARCHAR(20),            -- HH:MM
    end_time                VARCHAR(20),            -- HH:MM
    
    procedure_name          VARCHAR(255),           -- Name of procedure/surgery
    complexity_level        VARCHAR(50),            -- Simple | Moderate | Complex
    status                  VARCHAR(50) NOT NULL DEFAULT 'scheduled',  -- scheduled | in-progress | completed | cancelled
    notes                   TEXT,

    CONSTRAINT fk_doc_schedule_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_doctor_schedules_doctor_id ON doctor_schedules(doctor_id);
CREATE INDEX IF NOT EXISTS idx_doctor_schedules_scheduled_date ON doctor_schedules(scheduled_date);
CREATE INDEX IF NOT EXISTS idx_doctor_schedules_status ON doctor_schedules(status);
CREATE INDEX IF NOT EXISTS idx_doctor_schedules_deleted_at ON doctor_schedules(deleted_at);

-- ============================================================
-- Migration Status
-- ============================================================
-- All tables created successfully
-- Ready for Phase 1: IPD Patient Management
-- All foreign keys and indexes configured
-- Timestamps configured for GORM compatibility

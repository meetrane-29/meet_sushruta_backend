-- Performance Optimization - Database Indexes
-- Created: 2026-03-29
-- Purpose: Add indexes for frequently queried columns to improve query performance

-- ============================================================
-- APPOINTMENTS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_date_range ON appointments(appointment_date);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_status ON appointments(doctor_id, status);
CREATE INDEX IF NOT EXISTS idx_appointments_date_status ON appointments(appointment_date, status);

-- ============================================================
-- PRESCRIPTIONS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_prescriptions_doctor_id ON prescriptions(doctor_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_patient_id ON prescriptions(patient_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_status ON prescriptions(status);
CREATE INDEX IF NOT EXISTS idx_prescriptions_created_at ON prescriptions(created_at);
CREATE INDEX IF NOT EXISTS idx_prescriptions_doctor_status ON prescriptions(doctor_id, status);

-- ============================================================
-- LAB_REQUESTS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_lab_requests_doctor_id ON lab_requests(doctor_id);
CREATE INDEX IF NOT EXISTS idx_lab_requests_patient_id ON lab_requests(patient_id);
CREATE INDEX IF NOT EXISTS idx_lab_requests_status ON lab_requests(status);
CREATE INDEX IF NOT EXISTS idx_lab_requests_test_type ON lab_requests(test_type);
CREATE INDEX IF NOT EXISTS idx_lab_requests_priority ON lab_requests(priority);
CREATE INDEX IF NOT EXISTS idx_lab_requests_created_at ON lab_requests(created_at);
CREATE INDEX IF NOT EXISTS idx_lab_requests_status_created ON lab_requests(status, created_at);

-- ============================================================
-- PATIENTS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_patients_uhid ON patients(uhid);
CREATE INDEX IF NOT EXISTS idx_patients_email ON patients(email);
CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(first_name, last_name);

-- ============================================================
-- DOCTORS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_doctors_specialization ON doctors(specialization);
CREATE INDEX IF NOT EXISTS idx_doctors_email ON doctors(email);
CREATE INDEX IF NOT EXISTS idx_doctors_status ON doctors(status);

-- ============================================================
-- ADMISSION_RECORDS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_admission_records_doctor_id ON admission_records(doctor_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_patient_id ON admission_records(patient_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_bed_id ON admission_records(bed_id);
CREATE INDEX IF NOT EXISTS idx_admission_records_status ON admission_records(status);
CREATE INDEX IF NOT EXISTS idx_admission_records_admission_date ON admission_records(admission_date);
CREATE INDEX IF NOT EXISTS idx_admission_records_doctor_status ON admission_records(doctor_id, status);

-- ============================================================
-- PROGRESS_NOTES TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_progress_notes_admission_id ON progress_notes(admission_id);
CREATE INDEX IF NOT EXISTS idx_progress_notes_doctor_id ON progress_notes(doctor_id);
CREATE INDEX IF NOT EXISTS idx_progress_notes_created_at ON progress_notes(created_at);

-- ============================================================
-- BILLS TABLE INDEXES
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_bills_patient_id ON bills(patient_id);
CREATE INDEX IF NOT EXISTS idx_bills_doctor_id ON bills(doctor_id);
CREATE INDEX IF NOT EXISTS idx_bills_created_at ON bills(created_at);
CREATE INDEX IF NOT EXISTS idx_bills_status ON bills(status);

-- ============================================================
-- AUDIT_LOGS TABLE INDEXES (if exists)
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- Analyze updated statistics
ANALYZE;

-- Reception/OPD Module - Database Migration
-- Created: March 29, 2026
-- Purpose: Create 4 new tables for Reception/OPD functionality

-- ============================================================
-- 1. UHID (Unique Hospital ID) TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS uhids (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL UNIQUE,
  uhid VARCHAR(15) UNIQUE NOT NULL,
  hospital_code VARCHAR(5) NOT NULL DEFAULT 'MS',
  sequence_number INTEGER NOT NULL,
  issued_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Create indexes on uhids table
CREATE INDEX IF NOT EXISTS idx_uhids_patient_id ON uhids(patient_id);
CREATE INDEX IF NOT EXISTS idx_uhids_uhid ON uhids(uhid);
CREATE INDEX IF NOT EXISTS idx_uhids_hospital_code ON uhids(hospital_code);

-- ============================================================
-- 2. INSURANCE POLICIES TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS insurance_policies (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL,
  provider_name VARCHAR(100) NOT NULL,
  policy_number VARCHAR(50) UNIQUE NOT NULL,
  member_id VARCHAR(50),
  coverage_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
  copay_percentage DECIMAL(5, 2) DEFAULT 10,
  deductible_amount DECIMAL(10, 2) DEFAULT 0,
  valid_from DATE NOT NULL,
  valid_upto DATE NOT NULL,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Create indexes on insurance_policies table
CREATE INDEX IF NOT EXISTS idx_insurance_patient_id ON insurance_policies(patient_id);
CREATE INDEX IF NOT EXISTS idx_insurance_policy_number ON insurance_policies(policy_number);
CREATE INDEX IF NOT EXISTS idx_insurance_provider ON insurance_policies(provider_name);
CREATE INDEX IF NOT EXISTS idx_insurance_valid_dates ON insurance_policies(valid_from, valid_upto);

-- ============================================================
-- 3. WAITING LIST ENTRIES TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS waiting_list_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  appointment_id UUID,
  patient_id UUID NOT NULL,
  doctor_id UUID NOT NULL,
  token_number INTEGER NOT NULL,
  status VARCHAR(20) DEFAULT 'waiting' CHECK (status IN ('waiting', 'called', 'seen', 'completed', 'cancelled')),
  arrival_time TIMESTAMP,
  called_time TIMESTAMP,
  seen_time TIMESTAMP,
  completion_time TIMESTAMP,
  estimated_wait_time INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Create indexes on waiting_list_entries table
CREATE INDEX IF NOT EXISTS idx_waiting_list_appointment_id ON waiting_list_entries(appointment_id);
CREATE INDEX IF NOT EXISTS idx_waiting_list_patient_id ON waiting_list_entries(patient_id);
CREATE INDEX IF NOT EXISTS idx_waiting_list_doctor_id ON waiting_list_entries(doctor_id);
CREATE INDEX IF NOT EXISTS idx_waiting_list_status ON waiting_list_entries(status);
CREATE INDEX IF NOT EXISTS idx_waiting_list_token ON waiting_list_entries(token_number, doctor_id);
CREATE INDEX IF NOT EXISTS idx_waiting_list_dates ON waiting_list_entries(created_at);

-- ============================================================
-- 4. OPD RECEIPTS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS opd_receipts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bill_id UUID,
  appointment_id UUID,
  patient_id UUID NOT NULL,
  receipt_number VARCHAR(20) UNIQUE NOT NULL,
  receipt_date DATE DEFAULT CURRENT_DATE,
  consultation_fee DECIMAL(10, 2) NOT NULL DEFAULT 500,
  insurance_covered DECIMAL(10, 2) DEFAULT 0,
  patient_payable DECIMAL(10, 2) NOT NULL,
  paid_amount DECIMAL(10, 2) DEFAULT 0,
  payment_method VARCHAR(50),
  transaction_ref VARCHAR(100),
  receipt_status VARCHAR(20) DEFAULT 'generated' CHECK (receipt_status IN ('generated', 'partially_paid', 'paid', 'cancelled')),
  printed_count INTEGER DEFAULT 0,
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Create indexes on opd_receipts table
CREATE INDEX IF NOT EXISTS idx_opd_receipt_patient_id ON opd_receipts(patient_id);
CREATE INDEX IF NOT EXISTS idx_opd_receipt_appointment_id ON opd_receipts(appointment_id);
CREATE INDEX IF NOT EXISTS idx_opd_receipt_number ON opd_receipts(receipt_number);
CREATE INDEX IF NOT EXISTS idx_opd_receipt_status ON opd_receipts(receipt_status);
CREATE INDEX IF NOT EXISTS idx_opd_receipt_date ON opd_receipts(receipt_date);

-- ============================================================
-- Migration completed successfully
-- ============================================================
-- Run this script once to set up all tables:
-- psql -U postgres -d meet_sushruta -f 001_reception_opd_tables.sql
--
-- Verify with:
-- \dt+ 
-- (should show 4 new tables)
-- ============================================================

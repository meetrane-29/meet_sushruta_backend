-- Reception/OPD Module - Sample Data (Next 5 Days)
-- Created: March 29, 2026
-- Purpose: Seed database with realistic test data

BEGIN;

-- ============================================================
-- IMPORTANT NOTES:
-- 1. This script assumes users table exists with patients and doctors
-- 2. Run migration first: 001_reception_opd_tables.sql
-- 3. The Go backend should have created sample data in users, appointments, bills
-- ============================================================

-- ============================================================
-- STEP 1: Create Sample UHIDs for First 15 Patients
-- ============================================================
INSERT INTO uhids (patient_id, uhid, hospital_code, sequence_number, issued_date, is_active)
SELECT 
  u.id,
  CONCAT('MS-2026-', LPAD(ROW_NUMBER() OVER (ORDER BY u.id)::TEXT, 5, '0')),
  'MS',
  ROW_NUMBER() OVER (ORDER BY u.id),
  CURRENT_TIMESTAMP,
  true
FROM users u
WHERE u.role = 'patient' 
  AND u.is_active = true
  AND NOT EXISTS (SELECT 1 FROM uhids WHERE uhids.patient_id = u.id)
LIMIT 15
ON CONFLICT DO NOTHING;

-- ============================================================
-- STEP 1: Get or Create Doctor IDs (assuming doctors exist)
-- We'll reference them by role and assign appointments
-- ============================================================

-- ============================================================
-- STEP 2: Create Sample UHIDs for Patients
-- ============================================================
INSERT INTO uhids (patient_id, uhid, hospital_code, sequence_number, issued_date, is_active)
SELECT 
  u.id,
  CONCAT('MS-2026-', LPAD(ROW_NUMBER() OVER (ORDER BY u.id)::TEXT, 5, '0')),
  'MS',
  ROW_NUMBER() OVER (ORDER BY u.id),
  CURRENT_TIMESTAMP,
  true
FROM users u
WHERE u.role = 'patient' 
  AND u.is_active = true
  AND NOT EXISTS (SELECT 1 FROM uhids WHERE uhids.patient_id = u.id)
LIMIT 15;

-- ============================================================
-- STEP 3: Create Sample Insurance Policies
-- ============================================================
INSERT INTO insurance_policies (
  patient_id, 
  provider_name, 
  policy_number, 
  member_id, 
  coverage_amount, 
  copay_percentage, 
  deductible_amount, 
  valid_from, 
  valid_upto, 
  is_active
)
SELECT 
  u.id,
  CASE (RANDOM() * 3)::INT
    WHEN 0 THEN 'Aditya Birla Health Insurance'
    WHEN 1 THEN 'Star Health Insurance'
    WHEN 2 THEN 'HDFC ERGO Insurance'
    ELSE 'United India Insurance'
  END,
  CONCAT('POL-', TO_CHAR(CURRENT_DATE, 'YYYY'), '-', LPAD(ROW_NUMBER() OVER (ORDER BY u.id)::TEXT, 6, '0')),
  CONCAT('MEM-', LPAD(ROW_NUMBER() OVER (ORDER BY u.id)::TEXT, 8, '0')),
  (RANDOM() * (100000 - 25000) + 25000)::DECIMAL(10,2),
  (RANDOM() * (20 - 5) + 5)::DECIMAL(5,2),
  (RANDOM() * (10000 - 2000) + 2000)::DECIMAL(10,2),
  CURRENT_DATE,
  CURRENT_DATE + INTERVAL '1 year',
  true
FROM users u
WHERE u.role = 'patient' 
  AND u.is_active = true
  AND EXISTS (SELECT 1 FROM uhids WHERE uhids.patient_id = u.id)
LIMIT 10;

-- ============================================================
-- STEP 4: Create Sample Appointments for Next 5 Days
-- ============================================================
INSERT INTO appointments (
  patient_id,
  doctor_id,
  appointment_date,
  appointment_time,
  status,
  appointment_type,
  notes,
  created_at
)
SELECT
  p.id,
  d.id,
  CURRENT_DATE + (ROW_NUMBER() OVER (ORDER BY p.id) % 5)::INT,
  CASE (ROW_NUMBER() OVER (ORDER BY p.id) % 6)
    WHEN 0 THEN '09:00:00'
    WHEN 1 THEN '09:30:00'
    WHEN 2 THEN '10:00:00'
    WHEN 3 THEN '10:30:00'
    WHEN 4 THEN '14:00:00'
    ELSE '14:30:00'
  END,
  'scheduled',
  'OPD Consultation',
  CASE (RANDOM() * 2)::INT
    WHEN 0 THEN 'General checkup required'
    WHEN 1 THEN 'Follow-up consultation'
    ELSE 'New patient consultation'
  END,
  CURRENT_TIMESTAMP
FROM (
  SELECT u.id FROM users u 
  WHERE u.role = 'patient' 
    AND u.is_active = true
    AND EXISTS (SELECT 1 FROM uhids WHERE uhids.patient_id = u.id)
  LIMIT 12
) p
CROSS JOIN (
  SELECT u.id FROM users u 
  WHERE u.role = 'doctor' 
    AND u.is_active = true
  LIMIT 3
) d
WHERE NOT EXISTS (
  SELECT 1 FROM appointments 
  WHERE appointments.patient_id = p.id 
    AND appointments.doctor_id = d.id
    AND appointments.appointment_date >= CURRENT_DATE
)
LIMIT 15;

-- ============================================================
-- STEP 5: Create Sample Waiting List Entries
-- ============================================================
INSERT INTO waiting_list_entries (
  appointment_id,
  patient_id,
  doctor_id,
  token_number,
  status,
  arrival_time,
  estimated_wait_time
)
SELECT
  a.id,
  a.patient_id,
  a.doctor_id,
  ROW_NUMBER() OVER (PARTITION BY a.doctor_id, a.appointment_date ORDER BY a.appointment_time),
  CASE 
    WHEN CURRENT_DATE = a.appointment_date 
      AND CURRENT_TIME >= a.appointment_time::TIME - INTERVAL '30 minutes'
    THEN 'called'
    ELSE 'waiting'
  END,
  CASE 
    WHEN CURRENT_DATE = a.appointment_date 
      AND CURRENT_TIME >= a.appointment_time::TIME - INTERVAL '30 minutes'
    THEN CURRENT_TIMESTAMP
    ELSE NULL
  END,
  15 + (ROW_NUMBER() OVER (PARTITION BY a.doctor_id, a.appointment_date ORDER BY a.appointment_time) - 1) * 15
FROM appointments a
WHERE a.status = 'scheduled'
  AND a.appointment_date >= CURRENT_DATE
  AND a.appointment_date < CURRENT_DATE + INTERVAL '5 days'
  AND NOT EXISTS (
    SELECT 1 FROM waiting_list_entries 
    WHERE waiting_list_entries.appointment_id = a.id
  );

-- ============================================================
-- STEP 6: Create Sample OPD Receipts
-- ============================================================
INSERT INTO opd_receipts (
  appointment_id,
  patient_id,
  receipt_number,
  receipt_date,
  consultation_fee,
  insurance_covered,
  patient_payable,
  paid_amount,
  payment_method,
  transaction_ref,
  receipt_status,
  printed_count,
  notes
)
SELECT
  a.id,
  a.patient_id,
  CONCAT('RCP-2026-', LPAD(ROW_NUMBER() OVER (ORDER BY a.id)::TEXT, 5, '0')),
  CURRENT_DATE,
  500.00,
  COALESCE(ip.coverage_amount * ip.copay_percentage / 100, 0),
  500.00 - COALESCE((ip.coverage_amount * ip.copay_percentage / 100), 0),
  CASE WHEN RANDOM() > 0.3 THEN 500.00 - COALESCE((ip.coverage_amount * ip.copay_percentage / 100), 0) ELSE 0 END,
  CASE (RANDOM() * 3)::INT
    WHEN 0 THEN 'cash'
    WHEN 1 THEN 'card'
    ELSE 'upi'
  END,
  CONCAT('TXN-', TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDDHH24MISS'), '-', LPAD((RANDOM() * 10000)::INT::TEXT, 5, '0')),
  CASE WHEN RANDOM() > 0.3 THEN 'paid' ELSE 'generated' END,
  0,
  'OPD Consultation Receipt'
FROM appointments a
LEFT JOIN insurance_policies ip ON a.patient_id = ip.patient_id AND ip.is_active = true
WHERE a.appointment_date = CURRENT_DATE
  AND EXISTS (SELECT 1 FROM uhids WHERE uhids.patient_id = a.patient_id)
  AND NOT EXISTS (
    SELECT 1 FROM opd_receipts 
    WHERE opd_receipts.appointment_id = a.id
  )
LIMIT 10;

-- ============================================================
-- STEP 7: Display Summary of Created Data
-- ============================================================
-- You can run these queries to verify the data was created:
-- SELECT COUNT(*) as total_uhids FROM uhids;
-- SELECT COUNT(*) as total_insurance_policies FROM insurance_policies;
-- SELECT COUNT(*) as total_appointments FROM appointments WHERE appointment_date >= CURRENT_DATE;
-- SELECT COUNT(*) as total_waiting_list_entries FROM waiting_list_entries;
-- SELECT COUNT(*) as total_receipts FROM opd_receipts;

COMMIT;

-- ============================================================
-- Sample Data Import Complete
-- ============================================================
-- 
-- Summary of what was created:
-- 1. UHIDs for patients (MS-2026-00001, MS-2026-00002, etc.)
-- 2. Insurance policies (10 policies with realistic coverage amounts)
-- 3. Appointments (15 appointments for next 5 days across 3 doctors)
-- 4. Waiting list entries (auto-generated from appointments)
-- 5. OPD receipts (for today's appointments, 60-70% marked as paid)
--
-- To view the data created:
-- psql -d meet_sushruta -c "
--   SELECT 'UHIDs: ' || COUNT(*) as table_count FROM uhids
--   UNION ALL
--   SELECT 'Insurance: ' || COUNT(*) FROM insurance_policies
--   UNION ALL
--   SELECT 'Appointments: ' || COUNT(*) FROM appointments WHERE appointment_date >= CURRENT_DATE
--   UNION ALL
--   SELECT 'Waiting List: ' || COUNT(*) FROM waiting_list_entries
--   UNION ALL
--   SELECT 'Receipts: ' || COUNT(*) FROM opd_receipts;
-- "
-- ============================================================

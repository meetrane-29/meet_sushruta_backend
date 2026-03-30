-- ============================================================
-- COMPLETE DATABASE SETUP & APPOINTMENT DATA
-- ============================================================
-- Run this entire script in pgAdmin Query Tool or psql
-- This will handle ALL appointments setup for March 29 - April 3

-- Step 1: Clear existing appointment data (optional - comment out if you want to keep old data)
DELETE FROM appointments WHERE appointment_date >= '2026-03-29';

-- Step 2: Verify we have patients and doctors
-- VERIFY: Run these to check data exists
-- SELECT COUNT(*) as total_patients FROM patients;
-- SELECT COUNT(*) as total_doctors FROM doctors;

-- ============================================================
-- MARCH 29 (TODAY) - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 0), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 0), '2026-03-29', '09:00', 'pending', 'General checkup', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-03-29', '10:30', 'confirmed', 'Follow-up consultation', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-03-29', '14:00', 'in_progress', 'Neurological assessment', 'Patient status: in_progress. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-03-29', '11:00', 'pending', 'Pediatric consultation', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-03-29', '12:30', 'confirmed', 'Skin treatment', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- MARCH 30 (TOMORROW) - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-03-30', '09:00', 'pending', 'Chronic disease management', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 0), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-03-30', '10:30', 'confirmed', 'Orthopedic follow-up', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-03-30', '14:00', 'pending', 'Neurological consultation', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-03-30', '11:00', 'confirmed', 'Pediatric checkup', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-03-30', '12:30', 'pending', 'Dermatology appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- MARCH 31 - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 0), '2026-03-31', '09:30', 'pending', 'Eye surgery consultation', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-03-31', '10:00', 'confirmed', 'ENT examination', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 0), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-03-31', '13:00', 'pending', 'General medicine checkup', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-03-31', '15:00', 'confirmed', 'Psychiatric evaluation', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-03-31', '14:00', 'pending', 'Womens health appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- APRIL 1 - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-04-01', '09:00', 'pending', 'Cardiac assessment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-04-01', '10:30', 'confirmed', 'Bone density scan appointment', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-04-01', '14:00', 'pending', 'Brain imaging consultation', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 0), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-04-01', '11:00', 'confirmed', 'Vaccination appointment', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-04-01', '12:30', 'pending', 'Skin biopsy discussion', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- APRIL 2 - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 0), '2026-04-02', '09:30', 'pending', 'Post-operative checkup', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-04-02', '10:00', 'confirmed', 'Hearing test appointment', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-04-02', '13:00', 'pending', 'Blood work appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-04-02', '15:00', 'confirmed', 'Mental health session', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 0), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-04-02', '14:00', 'pending', 'Pregnancy monitoring appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- APRIL 3 - 5 Appointments
-- ============================================================
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-04-03', '09:00', 'pending', 'Treadmill test appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-04-03', '10:30', 'confirmed', 'Physical therapy session', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-04-03', '14:00', 'pending', 'MRI scan appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-04-03', '11:00', 'confirmed', 'Child immunization', 'Patient status: confirmed. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-04-03', '12:30', 'pending', 'Allergy testing appointment', 'Patient status: pending. Appointment recorded.', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- ============================================================
-- VERIFICATION QUERIES
-- ============================================================

-- Total count
SELECT '=== TOTAL APPOINTMENTS ===' as info;
SELECT COUNT(*) as total_appointments FROM appointments;

-- Day-wise breakdown
SELECT '=== APPOINTMENTS BY DATE ===' as info;
SELECT 
  appointment_date,
  COUNT(*) as daily_count,
  string_agg(DISTINCT status, ', ') as statuses
FROM appointments 
WHERE appointment_date BETWEEN '2026-03-29' AND '2026-04-03'
GROUP BY appointment_date 
ORDER BY appointment_date;

-- Status breakdown
SELECT '=== APPOINTMENTS BY STATUS ===' as info;
SELECT 
  status,
  COUNT(*) as count
FROM appointments
WHERE appointment_date BETWEEN '2026-03-29' AND '2026-04-03'
GROUP BY status;

-- Sample appointment data
SELECT '=== SAMPLE DATA (First 5) ===' as info;
SELECT 
  a.id,
  p.first_name || ' ' || p.last_name as patient_name,
  d.first_name || ' ' || d.last_name as doctor_name,
  a.appointment_date,
  a.appointment_time,
  a.status,
  a.reason
FROM appointments a
LEFT JOIN patients p ON a.patient_id = p.id
LEFT JOIN doctors d ON a.doctor_id = d.id
WHERE a.appointment_date BETWEEN '2026-03-29' AND '2026-04-03'
LIMIT 5;

-- ============================================================
-- SUCCESS MESSAGE
-- ============================================================
SELECT '✅ DATABASE SETUP COMPLETE!' as status,
       COUNT(*) as appointments_inserted
FROM appointments
WHERE appointment_date BETWEEN '2026-03-29' AND '2026-04-03';

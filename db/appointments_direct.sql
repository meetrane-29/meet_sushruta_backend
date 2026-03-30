-- Simple direct SQL insert for appointments
-- PostgreSQL script to insert appointment data directly
-- Run this in pgAdmin or psql command line

-- First, get the IDs we need
-- These are for 6 patients and 10 doctors (from the seed)
-- Patient IDs: the first 6 patient records
-- Doctor IDs: the first 10 doctor records  

-- Check if we have patients and doctors
-- SELECT COUNT(*) FROM patients;
-- SELECT COUNT(*) FROM doctors;

-- INSERT 30 APPOINTMENTS for MARCH 29-30-31 and APRIL 01-02-03

-- MARCH 29 (Today) - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) 
SELECT 
  gen_random_uuid(),
  p.id,
  d.id,
  '2026-03-29',
  times.time,
  status,
  reason,
  'Patient status: ' || status || '. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM (
  SELECT * FROM patients ORDER BY created_at LIMIT 5
) p
CROSS JOIN (
  SELECT * FROM doctors ORDER BY created_at LIMIT 1
) d
CROSS JOIN (
  SELECT '09:00' as time, 'pending' as status, 'General checkup' as reason
  UNION ALL SELECT '10:30', 'confirmed', 'Follow-up consultation'
  UNION ALL SELECT '14:00', 'in_progress', 'Neurological assessment'
  UNION ALL SELECT '11:00', 'pending', 'Pediatric consultation'
  UNION ALL SELECT '12:30', 'confirmed', 'Skin treatment'
) times;

-- MARCH 30 - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  (SELECT id FROM patients ORDER BY created_at OFFSET i LIMIT 1),
  (SELECT id FROM doctors ORDER BY created_at OFFSET i LIMIT 1),
  '2026-03-30',
  times.time,
  'pending',
  times.reason,
  'Patient status: pending. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM generate_series(0, 4) AS i
CROSS JOIN (
  SELECT '09:00' as time, 'Chronic disease management' as reason
  UNION ALL SELECT '10:30', 'Orthopedic follow-up'
  UNION ALL SELECT '14:00', 'Neurological consultation'
  UNION ALL SELECT '11:00', 'Pediatric checkup'
  UNION ALL SELECT '12:30', 'Dermatology appointment'
) times;

-- MARCH 31 - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET (ROW_NUMBER() OVER () - 1) % 6),
  (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET (ROW_NUMBER() OVER () - 1) % 10),
  '2026-03-31',
  times.time,
  'confirmed',
  times.reason,
  'Patient status: confirmed. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM generate_series(1, 5) AS i
CROSS JOIN (
  SELECT '09:30' as time, 'Eye surgery consultation' as reason
  UNION ALL SELECT '10:00', 'ENT examination'
  UNION ALL SELECT '13:00', 'General medicine checkup'
  UNION ALL SELECT '15:00', 'Psychiatric evaluation'
  UNION ALL SELECT '14:00', 'Women''s health appointment'
) times;

-- APRIL 1 - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  (SELECT id FROM patients ORDER BY random() LIMIT 1),
  (SELECT id FROM doctors ORDER BY random() LIMIT 1),
  '2026-04-01',
  times.time,
  'pending',
  times.reason,
  'Patient status: pending. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM (
  SELECT '09:00' as time, 'Cardiac assessment' as reason
  UNION ALL SELECT '10:30', 'Bone density scan appointment'
  UNION ALL SELECT '14:00', 'Brain imaging consultation'
  UNION ALL SELECT '11:00', 'Vaccination appointment'
  UNION ALL SELECT '12:30', 'Skin biopsy discussion'
) times;

-- APRIL 2 - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  (SELECT id FROM patients ORDER BY random() LIMIT 1),
  (SELECT id FROM doctors ORDER BY random() LIMIT 1),
  '2026-04-02',
  times.time,
  'confirmed',
  times.reason,
  'Patient status: confirmed. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM (
  SELECT '09:30' as time, 'Post-operative checkup' as reason
  UNION ALL SELECT '10:00', 'Hearing test appointment'
  UNION ALL SELECT '13:00', 'Blood work appointment'
  UNION ALL SELECT '15:00', 'Mental health session'
  UNION ALL SELECT '14:00', 'Pregnancy monitoring appointment'
) times;

-- APRIL 3 - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  (SELECT id FROM patients ORDER BY random() LIMIT 1),
  (SELECT id FROM doctors ORDER BY random() LIMIT 1),
  '2026-04-03',
  times.time,
  'pending',
  times.reason,
  'Patient status: pending. Appointment recorded.',
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000,
  (EXTRACT(EPOCH FROM NOW())::bigint) * 1000
FROM (
  SELECT '09:00' as time, 'Treadmill test appointment' as reason
  UNION ALL SELECT '10:30', 'Physical therapy session'
  UNION ALL SELECT '14:00', 'MRI scan appointment'
  UNION ALL SELECT '11:00', 'Child immunization'
  UNION ALL SELECT '12:30', 'Allergy testing appointment'
) times;

-- Verify the data
SELECT COUNT(*) as total_appointments FROM appointments;
SELECT appointment_date, COUNT(*) as count_per_day FROM appointments 
WHERE appointment_date BETWEEN '2026-03-29' AND '2026-04-03'
GROUP BY appointment_date ORDER BY appointment_date;

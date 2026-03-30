-- Direct appointment data insertion
-- Run this file directly in PostgreSQL to populate appointment data
-- Date: 29 March 2026
-- Format: INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)

-- Get doctor IDs (from seeded doctors - doctor@test.com patient@test.com etc)
-- Get patient IDs (from seeded patients - patient1@test.com, patient2@test.com etc)

-- TODAY (29 March) - 18 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at) 
SELECT 
  gen_random_uuid(), patients.id, doctors.id, '2026-03-29', '09:00', 'pending', 'General checkup', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL
FROM patients LIMIT 1, doctors LIMIT 1
ON CONFLICT DO NOTHING;

-- TOMORROW (30 March) - 5 appointments  
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)
VALUES
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 1), (SELECT id FROM doctors LIMIT 1 OFFSET 1), '2026-03-30', '09:00', 'pending', 'Chronic disease management', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 2), (SELECT id FROM doctors LIMIT 1 OFFSET 2), '2026-03-30', '10:30', 'confirmed', 'Orthopedic follow-up', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 3), (SELECT id FROM doctors LIMIT 1 OFFSET 3), '2026-03-30', '14:00', 'pending', 'Neurological consultation', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 4), (SELECT id FROM doctors LIMIT 1 OFFSET 4), '2026-03-30', '11:00', 'confirmed', 'Pediatric checkup', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 5), (SELECT id FROM doctors LIMIT 1 OFFSET 5), '2026-03-30', '12:30', 'pending', 'Dermatology appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL);

-- 31 MARCH - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)
VALUES
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1), (SELECT id FROM doctors LIMIT 1 OFFSET 6), '2026-03-31', '09:30', 'pending', 'Eye surgery consultation', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 1), (SELECT id FROM doctors LIMIT 1 OFFSET 7), '2026-03-31', '10:00', 'confirmed', 'ENT examination', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 2), (SELECT id FROM doctors LIMIT 1 OFFSET 8), '2026-03-31', '13:00', 'pending', 'General medicine checkup', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 3), (SELECT id FROM doctors LIMIT 1 OFFSET 9), '2026-03-31', '15:00', 'confirmed', 'Psychiatric evaluation', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 4), (SELECT id FROM doctors LIMIT 1), '2026-03-31', '14:00', 'pending', 'Women''s health appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL);

-- 1 APRIL - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)
VALUES
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 5), (SELECT id FROM doctors LIMIT 1 OFFSET 1), '2026-04-01', '09:00', 'pending', 'Cardiac assessment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1), (SELECT id FROM doctors LIMIT 1 OFFSET 2), '2026-04-01', '10:30', 'confirmed', 'Bone density scan appointment', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 1), (SELECT id FROM doctors LIMIT 1 OFFSET 3), '2026-04-01', '14:00', 'pending', 'Brain imaging consultation', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 2), (SELECT id FROM doctors LIMIT 1 OFFSET 4), '2026-04-01', '11:00', 'confirmed', 'Vaccination appointment', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 3), (SELECT id FROM doctors LIMIT 1 OFFSET 5), '2026-04-01', '12:30', 'pending', 'Skin biopsy discussion', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL);

-- 2 APRIL - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)
VALUES
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 4), (SELECT id FROM doctors LIMIT 1 OFFSET 6), '2026-04-02', '09:30', 'pending', 'Post-operative checkup', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 5), (SELECT id FROM doctors LIMIT 1 OFFSET 7), '2026-04-02', '10:00', 'confirmed', 'Hearing test appointment', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1), (SELECT id FROM doctors LIMIT 1 OFFSET 8), '2026-04-02', '13:00', 'pending', 'Blood work appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 1), (SELECT id FROM doctors LIMIT 1 OFFSET 9), '2026-04-02', '15:00', 'confirmed', 'Mental health session', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 2), (SELECT id FROM doctors LIMIT 1), '2026-04-02', '14:00', 'pending', 'Pregnancy monitoring appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL);

-- 3 APRIL - 5 appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at, deleted_at)
VALUES
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 3), (SELECT id FROM doctors LIMIT 1 OFFSET 1), '2026-04-03', '09:00', 'pending', 'Treadmill test appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 4), (SELECT id FROM doctors LIMIT 1 OFFSET 2), '2026-04-03', '10:30', 'confirmed', 'Physical therapy session', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 5), (SELECT id FROM doctors LIMIT 1 OFFSET 3), '2026-04-03', '14:00', 'pending', 'MRI scan appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1), (SELECT id FROM doctors LIMIT 1 OFFSET 4), '2026-04-03', '11:00', 'confirmed', 'Child immunization', 'Patient status: confirmed. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL),
  (gen_random_uuid(), (SELECT id FROM patients LIMIT 1 OFFSET 1), (SELECT id FROM doctors LIMIT 1 OFFSET 5), '2026-04-03', '12:30', 'pending', 'Allergy testing appointment', 'Patient status: pending. Appointment recorded.', EXTRACT(EPOCH FROM NOW())::bigint * 1000, EXTRACT(EPOCH FROM NOW())::bigint * 1000, NULL);

-- Verify insertions
SELECT COUNT(*) as total_appointments FROM appointments;

-- ============================================================
-- SIMPLE DIRECT INSERT - Copy aur paste karo pgAdmin mein
-- ============================================================
-- Get patient and doctor IDs first by running:
-- SELECT id, email FROM patients LIMIT 6;
-- SELECT id, email FROM doctors ORDER BY created_at LIMIT 10;

-- MARCH 29 (TODAY) - 5 Appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1), '2026-03-29', '09:00', 'pending', 'General checkup', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-03-29', '10:30', 'confirmed', 'Follow-up consultation', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-03-29', '14:00', 'in_progress', 'Neurological assessment', 'Patient status: in_progress', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-03-29', '11:00', 'pending', 'Pediatric consultation', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-03-29', '12:30', 'confirmed', 'Skin treatment', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- MARCH 30 (Tomorrow) - 5 Appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-03-30', '09:00', 'pending', 'Chronic disease management', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-03-30', '10:30', 'confirmed', 'Orthopedic follow-up', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-03-30', '14:00', 'pending', 'Neurological consultation', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-03-30', '11:00', 'confirmed', 'Pediatric checkup', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-03-30', '12:30', 'pending', 'Dermatology appointment', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- MARCH 31 - 5 Appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1), '2026-03-31', '09:30', 'pending', 'Eye surgery consultation', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-03-31', '10:00', 'confirmed', 'ENT examination', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-03-31', '13:00', 'pending', 'General medicine checkup', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-03-31', '15:00', 'confirmed', 'Psychiatric evaluation', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-03-31', '14:00', 'pending', 'Women ''s health appointment', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- APRIL 1 - 5 Appointments  
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-04-01', '09:00', 'pending', 'Cardiac assessment', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-04-01', '10:30', 'confirmed', 'Bone density scan', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-04-01', '14:00', 'pending', 'Brain imaging consultation', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-04-01', '11:00', 'confirmed', 'Vaccination appointment', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-04-01', '12:30', 'pending', 'Skin biopsy discussion', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- APRIL 2 - 5 Appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1), '2026-04-02', '09:30', 'pending', 'Post-operative checkup', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 1), '2026-04-02', '10:00', 'confirmed', 'Hearing test appointment', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 2), '2026-04-02', '13:00', 'pending', 'Blood work appointment', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 3), '2026-04-02', '15:00', 'confirmed', 'Mental health session', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 4), '2026-04-02', '14:00', 'pending', 'Pregnancy monitoring', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- APRIL 3 - 5 Appointments
INSERT INTO appointments (id, patient_id, doctor_id, appointment_date, appointment_time, status, reason, notes, created_at, updated_at) VALUES
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 1), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 5), '2026-04-03', '09:00', 'pending', 'Treadmill test', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 2), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 6), '2026-04-03', '10:30', 'confirmed', 'Physical therapy', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 3), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 7), '2026-04-03', '14:00', 'pending', 'MRI scan appointment', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 4), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 8), '2026-04-03', '11:00', 'confirmed', 'Child immunization', 'Patient status: confirmed', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint),
(gen_random_uuid(), (SELECT id FROM patients ORDER BY created_at LIMIT 1 OFFSET 5), (SELECT id FROM doctors ORDER BY created_at LIMIT 1 OFFSET 9), '2026-04-03', '12:30', 'pending', 'Allergy testing', 'Patient status: pending', (extract(epoch from now())*1000)::bigint, (extract(epoch from now())*1000)::bigint);

-- VERIFY - Check if data inserted
SELECT COUNT(*) as total_appointments FROM appointments;
SELECT appointment_date, COUNT(*) as daily_count FROM appointments 
WHERE appointment_date >= '2026-03-29' 
GROUP BY appointment_date ORDER BY appointment_date;

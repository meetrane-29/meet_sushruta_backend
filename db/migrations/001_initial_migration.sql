-- ============================================================
-- meet_sushruta - Initial Database Migration
-- Generated from seed.go model analysis
-- Run this ONCE on a fresh database before seeding
-- ============================================================

-- Enable UUID extension (PostgreSQL)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- 1. USERS
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    email       VARCHAR(255) NOT NULL UNIQUE,
    phone       VARCHAR(20)  NOT NULL UNIQUE,
    first_name  VARCHAR(100) NOT NULL,
    last_name   VARCHAR(100) NOT NULL,
    password    TEXT         NOT NULL,
    role        VARCHAR(50)  NOT NULL,   -- admin | doctor | nurse | pharmacy | lab | patient
    active      BOOLEAN      NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email      ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role       ON users(role);

-- ============================================================
-- 2. SPECIALIZATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS specializations (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    name       VARCHAR(100) NOT NULL UNIQUE,
    category   VARCHAR(100) NOT NULL,   -- Medicine | Surgery | Diagnostics
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_specializations_deleted_at ON specializations(deleted_at);

-- ============================================================
-- 3. HOSPITALS
-- ============================================================
CREATE TABLE IF NOT EXISTS hospitals (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,

    name            VARCHAR(200) NOT NULL UNIQUE,
    address         TEXT,
    city            VARCHAR(100),
    state           VARCHAR(100),
    phone           VARCHAR(20),
    email           VARCHAR(255),
    total_beds      INTEGER      NOT NULL DEFAULT 0,
    available_beds  INTEGER      NOT NULL DEFAULT 0,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    is_verified     BOOLEAN      NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_hospitals_deleted_at ON hospitals(deleted_at);

-- ============================================================
-- 4. BEDS
-- ============================================================
CREATE TABLE IF NOT EXISTS beds (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    bed_number  VARCHAR(50)  NOT NULL UNIQUE,
    ward        VARCHAR(100) NOT NULL,   -- ICU | General | Private | Semi-Private
    room        VARCHAR(100),
    floor       INTEGER,
    bed_type    VARCHAR(100) NOT NULL,   -- ICU | General | Private | Semi-Private
    status      VARCHAR(50)  NOT NULL DEFAULT 'available',  -- available | occupied | maintenance
    features    TEXT,
    daily_rate  NUMERIC(10,2) NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_beds_deleted_at ON beds(deleted_at);
CREATE INDEX IF NOT EXISTS idx_beds_status     ON beds(status);
CREATE INDEX IF NOT EXISTS idx_beds_ward       ON beds(ward);

-- ============================================================
-- 5. MEDICAL EQUIPMENT
-- ============================================================
CREATE TABLE IF NOT EXISTS medical_equipments (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at         TIMESTAMPTZ,
    updated_at         TIMESTAMPTZ,
    deleted_at         TIMESTAMPTZ,

    equipment_name     VARCHAR(200) NOT NULL UNIQUE,
    equipment_type     VARCHAR(100) NOT NULL,
    model              VARCHAR(200),
    manufacturer       VARCHAR(200),
    location           VARCHAR(200),
    status             VARCHAR(50)  NOT NULL DEFAULT 'working',  -- working | repair | under_maintenance
    critical_equipment BOOLEAN      NOT NULL DEFAULT FALSE,
    daily_rental_rate  NUMERIC(10,2) NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_medical_equipments_deleted_at ON medical_equipments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_medical_equipments_status     ON medical_equipments(status);

-- ============================================================
-- 6. OPERATION THEATRES
-- ============================================================
CREATE TABLE IF NOT EXISTS operation_theatres (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,

    theatre_name   VARCHAR(200) NOT NULL UNIQUE,
    floor          INTEGER,
    capacity       INTEGER      NOT NULL DEFAULT 1,
    status         VARCHAR(50)  NOT NULL DEFAULT 'available',  -- available | occupied | maintenance
    features       TEXT,
    equipment_list TEXT
);
CREATE INDEX IF NOT EXISTS idx_operation_theatres_deleted_at ON operation_theatres(deleted_at);

-- ============================================================
-- 7. DOCTORS
-- ============================================================
CREATE TABLE IF NOT EXISTS doctors (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,

    user_id          UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    specialization   VARCHAR(100) NOT NULL,
    license_number   VARCHAR(100) NOT NULL UNIQUE,
    consultation_fee NUMERIC(10,2) NOT NULL DEFAULT 0,
    department       VARCHAR(100),
    bio              TEXT
);
CREATE INDEX IF NOT EXISTS idx_doctors_deleted_at ON doctors(deleted_at);
CREATE INDEX IF NOT EXISTS idx_doctors_user_id    ON doctors(user_id);

-- ============================================================
-- 8. NURSES
-- ============================================================
CREATE TABLE IF NOT EXISTS nurses (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,

    user_id        UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    license_number VARCHAR(100) NOT NULL UNIQUE,
    role           VARCHAR(100) NOT NULL,   -- Staff Nurse | ICU Nurse | Charge Nurse | OT Nurse | Pediatric Nurse
    department     VARCHAR(100),
    shift          VARCHAR(50)              -- morning | evening | night
);
CREATE INDEX IF NOT EXISTS idx_nurses_deleted_at ON nurses(deleted_at);
CREATE INDEX IF NOT EXISTS idx_nurses_user_id    ON nurses(user_id);

-- ============================================================
-- 9. PATIENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS patients (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,

    user_id           UUID         NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    date_of_birth     VARCHAR(20),
    gender            VARCHAR(20),
    blood_group       VARCHAR(10),
    address           TEXT,
    emergency_contact VARCHAR(20),
    medical_history   TEXT,
    allergies         TEXT
);
CREATE INDEX IF NOT EXISTS idx_patients_deleted_at ON patients(deleted_at);
CREATE INDEX IF NOT EXISTS idx_patients_user_id    ON patients(user_id);

-- ============================================================
-- 10. APPOINTMENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS appointments (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,

    patient_id       UUID         NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id        UUID         NOT NULL REFERENCES doctors(id)  ON DELETE CASCADE,
    appointment_date VARCHAR(20)  NOT NULL,   -- YYYY-MM-DD
    appointment_time VARCHAR(10)  NOT NULL,   -- HH:MM
    status           VARCHAR(50)  NOT NULL DEFAULT 'pending',
                                              -- pending | confirmed | in_progress | completed | cancelled
    reason           TEXT,
    notes            TEXT
);
CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at       ON appointments(deleted_at);
CREATE INDEX IF NOT EXISTS idx_appointments_patient_id       ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id        ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_status           ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_date             ON appointments(appointment_date);

-- ============================================================
-- 11. VITALS
-- ============================================================
CREATE TABLE IF NOT EXISTS vitals (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,

    patient_id       UUID           NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    temperature      NUMERIC(5,2),                  -- °C
    blood_pressure   VARCHAR(20),                   -- e.g. 120/80
    heart_rate       INTEGER,                       -- bpm
    respiratory_rate INTEGER,                       -- breaths/min
    weight           NUMERIC(6,2),                  -- kg
    height           NUMERIC(6,2),                  -- cm
    blood_sugar      NUMERIC(6,2),                  -- mg/dL
    oxygen           INTEGER,                       -- SpO2 %
    recorded_at      VARCHAR(30)    NOT NULL        -- YYYY-MM-DD HH:MM
);
CREATE INDEX IF NOT EXISTS idx_vitals_deleted_at  ON vitals(deleted_at);
CREATE INDEX IF NOT EXISTS idx_vitals_patient_id  ON vitals(patient_id);

-- ============================================================
-- 12. MEDICINES
-- ============================================================
CREATE TABLE IF NOT EXISTS medicines (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,

    name           VARCHAR(200) NOT NULL UNIQUE,
    generic_name   VARCHAR(200),
    dosage         VARCHAR(100),
    stock_quantity INTEGER      NOT NULL DEFAULT 0,
    reorder_level  INTEGER      NOT NULL DEFAULT 0,
    price          NUMERIC(10,2) NOT NULL DEFAULT 0,
    description    TEXT,
    manufacturer   VARCHAR(200),
    expiry_date    VARCHAR(20),
    active         BOOLEAN      NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_medicines_deleted_at ON medicines(deleted_at);
CREATE INDEX IF NOT EXISTS idx_medicines_name       ON medicines(name);

-- ============================================================
-- 13. PRESCRIPTIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS prescriptions (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at        TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,

    patient_id        UUID         NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id         UUID         NOT NULL REFERENCES doctors(id)  ON DELETE CASCADE,
    appointment_id    UUID         REFERENCES appointments(id)      ON DELETE SET NULL,
    prescription_date VARCHAR(20)  NOT NULL,   -- YYYY-MM-DD
    issued_at         VARCHAR(30),             -- YYYY-MM-DD HH:MM
    status            VARCHAR(50)  NOT NULL DEFAULT 'active',  -- active | dispensed | expired
    notes             TEXT
);
CREATE INDEX IF NOT EXISTS idx_prescriptions_deleted_at    ON prescriptions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_prescriptions_patient_id    ON prescriptions(patient_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_doctor_id     ON prescriptions(doctor_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_appointment_id ON prescriptions(appointment_id);

-- ============================================================
-- 14. PRESCRIPTION ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS prescription_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at      TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,

    prescription_id UUID         NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    medicine_id     UUID         NOT NULL REFERENCES medicines(id)      ON DELETE RESTRICT,
    dosage          VARCHAR(100),
    frequency       VARCHAR(100),
    duration        VARCHAR(100),
    time_of_day     VARCHAR(100),
    quantity        INTEGER      NOT NULL DEFAULT 0,
    instructions    TEXT
);
CREATE INDEX IF NOT EXISTS idx_prescription_items_deleted_at     ON prescription_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_prescription_items_prescription_id ON prescription_items(prescription_id);

-- ============================================================
-- 15. LAB REQUESTS
-- ============================================================
CREATE TABLE IF NOT EXISTS lab_requests (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,

    patient_id     UUID         NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id      UUID         NOT NULL REFERENCES doctors(id)  ON DELETE CASCADE,
    test_type      VARCHAR(100) NOT NULL,   -- BloodWork | Metabolic | Urinalysis ...
    test_name      VARCHAR(200) NOT NULL,
    requested_at   VARCHAR(30)  NOT NULL,   -- YYYY-MM-DD HH:MM
    status         VARCHAR(50)  NOT NULL DEFAULT 'pending',  -- pending | in_progress | completed
    priority       VARCHAR(50)  NOT NULL DEFAULT 'normal',   -- normal | urgent | stat
    result_url     TEXT,
    completed_at   VARCHAR(30),
    result_summary TEXT
);
CREATE INDEX IF NOT EXISTS idx_lab_requests_deleted_at ON lab_requests(deleted_at);
CREATE INDEX IF NOT EXISTS idx_lab_requests_patient_id ON lab_requests(patient_id);
CREATE INDEX IF NOT EXISTS idx_lab_requests_doctor_id  ON lab_requests(doctor_id);
CREATE INDEX IF NOT EXISTS idx_lab_requests_status     ON lab_requests(status);

-- ============================================================
-- 16. BILLS
-- ============================================================
CREATE TABLE IF NOT EXISTS bills (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at       TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,

    patient_id       UUID           NOT NULL REFERENCES patients(id)     ON DELETE CASCADE,
    appointment_id   UUID           REFERENCES appointments(id)          ON DELETE SET NULL,
    bill_number      VARCHAR(50)    NOT NULL UNIQUE,
    bill_date        VARCHAR(20)    NOT NULL,   -- YYYY-MM-DD
    consultation_fee NUMERIC(10,2)  NOT NULL DEFAULT 0,
    lab_tests        NUMERIC(10,2)  NOT NULL DEFAULT 0,
    medicines        NUMERIC(10,2)  NOT NULL DEFAULT 0,
    other_charges    NUMERIC(10,2)  NOT NULL DEFAULT 0,
    discount         NUMERIC(10,2)  NOT NULL DEFAULT 0,
    tax_percentage   NUMERIC(5,2)   NOT NULL DEFAULT 0,
    tax_amount       NUMERIC(10,2)  NOT NULL DEFAULT 0,
    total_amount     NUMERIC(10,2)  NOT NULL DEFAULT 0,
    paid_amount      NUMERIC(10,2)  NOT NULL DEFAULT 0,
    due_amount       NUMERIC(10,2)  NOT NULL DEFAULT 0,
    payment_status   VARCHAR(50)    NOT NULL DEFAULT 'pending',  -- pending | paid | partial | cancelled
    payment_method   VARCHAR(50),                                -- upi | cash | card | online
    payment_date     VARCHAR(20),
    notes            TEXT
);
CREATE INDEX IF NOT EXISTS idx_bills_deleted_at    ON bills(deleted_at);
CREATE INDEX IF NOT EXISTS idx_bills_patient_id    ON bills(patient_id);
CREATE INDEX IF NOT EXISTS idx_bills_appointment_id ON bills(appointment_id);
CREATE INDEX IF NOT EXISTS idx_bills_payment_status ON bills(payment_status);

-- ============================================================
-- 17. BILL ITEMS
-- ============================================================
CREATE TABLE IF NOT EXISTS bill_items (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,

    bill_id     UUID           NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    item_type   VARCHAR(50)    NOT NULL,   -- consultation | medicine | lab | other
    description VARCHAR(500)   NOT NULL,
    quantity    INTEGER        NOT NULL DEFAULT 1,
    unit_price  NUMERIC(10,2)  NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_bill_items_deleted_at ON bill_items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_bill_items_bill_id    ON bill_items(bill_id);

-- ============================================================
-- 18. OPERATION SCHEDULES
-- ============================================================
CREATE TABLE IF NOT EXISTS operation_schedules (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    deleted_at           TIMESTAMPTZ,

    theatre_id           UUID         NOT NULL REFERENCES operation_theatres(id) ON DELETE RESTRICT,
    patient_id           UUID         NOT NULL REFERENCES patients(id)           ON DELETE CASCADE,
    surgie_doctor_id     UUID         NOT NULL REFERENCES users(id)              ON DELETE RESTRICT,
    operation_date       VARCHAR(20)  NOT NULL,   -- YYYY-MM-DD
    operation_time       VARCHAR(10)  NOT NULL,   -- HH:MM
    estimated_duration   INTEGER      NOT NULL DEFAULT 60,  -- minutes
    actual_duration      INTEGER,
    operation_type       VARCHAR(200) NOT NULL,
    status               VARCHAR(50)  NOT NULL DEFAULT 'scheduled',
                                                  -- scheduled | in_progress | completed | cancelled
    diagnosis            TEXT,
    pre_operative_notes  TEXT,
    post_operative_notes TEXT
);
CREATE INDEX IF NOT EXISTS idx_operation_schedules_deleted_at  ON operation_schedules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_operation_schedules_theatre_id  ON operation_schedules(theatre_id);
CREATE INDEX IF NOT EXISTS idx_operation_schedules_patient_id  ON operation_schedules(patient_id);
CREATE INDEX IF NOT EXISTS idx_operation_schedules_status      ON operation_schedules(status);

-- ============================================================
-- END OF MIGRATION
-- ============================================================



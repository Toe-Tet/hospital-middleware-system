CREATE TABLE patients (
    id SERIAL PRIMARY KEY,

    hospital_id INTEGER NOT NULL
        REFERENCES hospitals(id)
        ON DELETE CASCADE,

    first_name_th VARCHAR(100),
    middle_name_th VARCHAR(100),
    last_name_th VARCHAR(100),

    first_name_en VARCHAR(100) NOT NULL,
    middle_name_en VARCHAR(100),
    last_name_en VARCHAR(100) NOT NULL,

    patient_hn VARCHAR(50) NOT NULL,
    passport_id VARCHAR(50),
    national_id VARCHAR(50),

    date_of_birth DATE NOT NULL,

    gender VARCHAR(1) NOT NULL
        CHECK (gender IN ('M', 'F')),

    email VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    UNIQUE (hospital_id, patient_hn),
    UNIQUE (email)
);

CREATE INDEX idx_patients_hospital_id
    ON patients(hospital_id);

CREATE INDEX idx_patients_national_id
    ON patients(national_id);

CREATE INDEX idx_patients_passport_id
    ON patients(passport_id);

CREATE INDEX idx_patients_email
    ON patients(email);

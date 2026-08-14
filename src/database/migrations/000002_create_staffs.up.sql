CREATE TABLE staffs (
    id SERIAL PRIMARY KEY,
    hospital_id INTEGER NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_staffs_hospital_id ON staffs(hospital_id);
CREATE UNIQUE INDEX idx_staffs_username_unique ON staffs(username);
CREATE INDEX idx_staffs_status ON staffs(status);

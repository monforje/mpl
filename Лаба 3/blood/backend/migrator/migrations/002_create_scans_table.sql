-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS scans (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    sex VARCHAR(10) NOT NULL,
    hemoglobin DECIMAL(10, 2),
    erythrocytes DECIMAL(10, 2),
    hematocrit DECIMAL(10, 2),
    mcv DECIMAL(10, 2),
    leukocytes DECIMAL(10, 2),
    neutrophils DECIMAL(10, 2),
    lymphocytes DECIMAL(10, 2),
    monocytes DECIMAL(10, 2),
    eosinophils DECIMAL(10, 2),
    basophils DECIMAL(10, 2),
    platelets DECIMAL(10, 2),
    mpv DECIMAL(10, 2),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_scans_user_id ON scans(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_scans_user_id;
DROP TABLE IF EXISTS scans;

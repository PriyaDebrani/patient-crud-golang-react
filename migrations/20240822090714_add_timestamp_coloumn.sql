-- +goose Up
ALTER TABLE patients
ADD COLUMN created_at timestamptz ,
ADD COLUMN updated_at timestamptz ;

-- +goose Down
ALTER TABLE patients
DROP COLUMN created_at,
DROP COLUMN updated_at;

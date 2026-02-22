-- +goose Up
CREATE TABLE department (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    parent_id   INT REFERENCES department(id) ON DELETE SET NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE department;
-- +goose Up
CREATE TABLE employee(
    id              SERIAL PRIMARY KEY,
    department_id   INT NOT NULL REFERENCES department(id),
    full_name       VARCHAR(200) NOT NULL,
    position        VARCHAR(200) NOT NULL,
    hired_at        DATE,
    created_at      TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE employee;
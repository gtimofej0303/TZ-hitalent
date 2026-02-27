-- +goose Up
CREATE TABLE employees(
    id              SERIAL PRIMARY KEY,
    department_id   INT NOT NULL REFERENCES departments(id),
    fullname       VARCHAR(200) NOT NULL,
    position        VARCHAR(200) NOT NULL,
    hired_at        DATE,
    created_at      TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE employees;
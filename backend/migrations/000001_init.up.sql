-- Enable UUID generation.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users -----------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(120) NOT NULL,
    email         VARCHAR(160) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'Employee',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- employees -------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employees (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_code VARCHAR(40)  NOT NULL,
    first_name    VARCHAR(80)  NOT NULL,
    last_name     VARCHAR(80)  NOT NULL,
    email         VARCHAR(160) NOT NULL,
    phone         VARCHAR(20),
    department    VARCHAR(80),
    designation   VARCHAR(80),
    salary        NUMERIC(12,2) NOT NULL DEFAULT 0,
    joining_date  DATE         NOT NULL,
    status        VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_code  ON employees (employee_code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_email ON employees (email);
CREATE INDEX IF NOT EXISTS idx_employees_department   ON employees (department);
CREATE INDEX IF NOT EXISTS idx_employees_status       ON employees (status);
CREATE INDEX IF NOT EXISTS idx_employees_joining_date ON employees (joining_date);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at   ON employees (deleted_at);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(64)  NOT NULL,
    image TEXT,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15),
    refresh TEXT,
    password TEXT,
    role VARCHAR(9) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE resumes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    url TEXT NOT NULL,
    salary NUMERIC NOT NULL DEFAULT 0,
    job_title VARCHAR(100) NOT NULL,
    region VARCHAR(150) NOT NULL,
    job_location VARCHAR(100) NOT NULL DEFAULT 'offline',
    job_type VARCHAR(25) NOT NULL,
    experience INT NOT NULL,
    template VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

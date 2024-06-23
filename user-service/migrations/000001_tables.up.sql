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
    filename VARCHAR(200) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    job_title VARCHAR(100) NOT NULL,
    summary TEXT,
    website TEXT,
    profile_image TEXT,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE locations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    city VARCHAR(100),
    country_code VARCHAR(5),
    region VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    network VARCHAR(50) NOT NULL,
    username VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE works (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    position VARCHAR(100) NOT NULL,
    company VARCHAR(200) NOT NULL,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    location VARCHAR(100) NOT NULL,
    summary TEXT,
    skills TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE projects (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    url TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE educations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    institution VARCHAR(300) NOT NULL,
    area VARCHAR(200) NOT NULL,
    location VARCHAR(200) NOT NULL,
    study_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE,
    score VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE courses (
    id UUID PRIMARY KEY,
    education_id UUID NOT NULL,
    course_name VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (education_id) REFERENCES educations(id)
);

CREATE TABLE certificates (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    title VARCHAR(100) NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    issuer VARCHAR(64) NOT NULL,
    score VARCHAR(50),
    url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE hard_skills (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    level VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE soft_skills (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE languages (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    language VARCHAR(100) NOT NULL,
    fluency VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

CREATE TABLE interests (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    resume_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id)
);

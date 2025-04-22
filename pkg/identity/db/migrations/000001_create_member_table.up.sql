CREATE TABLE IF NOT EXISTS member (
    id BIGSERIAL PRIMARY KEY,
    login_name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    currency VARCHAR(50) NOT NULL,
    full_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(255) NOT NULL,
    email_verify_status SMALLINT NOT NULL DEFAULT 1,
    phone_verify_status SMALLINT NOT NULL DEFAULT 1,
    address VARCHAR(255),
    created_at  TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at  TIMESTAMP DEFAULT timezone('UTC', now())
);

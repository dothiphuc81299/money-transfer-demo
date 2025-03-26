CREATE TABLE IF NOT EXISTS member_account (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGSERIAL NOT NULL,
    login_name VARCHAR(255) UNIQUE NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    currency VARCHAR(50) NOT NULL,
    balance DECIMAL NOT NULL,
    outstanding_balance DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE  IF NOT EXISTS deposit (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(255) NOT NULL,
    bank_account_id BIGINT  NULL,
    member_id BIGINT NOT NULL,
    login_name VARCHAR(255) NOT NULL,
    payment_method_code VARCHAR(255) NOT NULL,
    ref_code VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    detail JSONB  NULL,
    amount NUMERIC(20, 2) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', now())
);

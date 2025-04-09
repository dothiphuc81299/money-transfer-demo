CREATE TABLE IF NOT EXISTS member_account_hold (
    id BIGSERIAL PRIMARY KEY,
    member_account_id BIGINT NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    is_hold BOOLEAN NOT NULL DEFAULT false,
    note VARCHAR(255),
    created_at TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', now())
)
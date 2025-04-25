CREATE TABLE IF NOT EXISTS member_payment_account(
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL,
    payment_method_code VARCHAR(200) NOT NULL,
    detail JSONB NULL,
    verify_status SMALLINT NOT NULL DEFAULT 1,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', now())
)


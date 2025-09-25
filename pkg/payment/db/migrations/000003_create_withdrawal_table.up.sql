CREATE TABLE  IF NOT EXISTS withdrawal (
    id SERIAL PRIMARY KEY,
    member_payment_account_id BIGINT NOT NULL,
    payment_method_code VARCHAR(255) NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    status SMALLINT NOT NULL,
    member_id BIGINT NOT NULL,
    login_name VARCHAR(255) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    gross_amount NUMERIC(20, 2) NOT NULL,
    charge_amount NUMERIC(20, 2)  NULL,
    net_amount NUMERIC(20, 2)  NULL,
    detail JSONB NOT NULL,
    bank_account_id BIGINT  NULL,
    created_at TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', now()) 
)
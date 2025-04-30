CREATE TABLE IF NOT EXISTS transfer_timeline (
    id BIGSERIAL PRIMARY KEY,
    transfer_id BIGINT NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    message VARCHAR(255) NOT NULL,
    note VARCHAR(255), 
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT timezone('UTC', now())
);
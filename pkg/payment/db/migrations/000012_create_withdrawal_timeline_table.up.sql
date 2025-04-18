CREATE TABLE withdrawal_timeline (
    id SERIAL PRIMARY KEY,
    withdrawal_id BIGINT NOT NULL,
    message VARCHAR(255) NOT NULL,
    additional_content JSONB ,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT timezone('UTC', now())
);

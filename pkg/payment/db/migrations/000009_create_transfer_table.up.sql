CREATE TABLE transfer (
    id BIGSERIAL PRIMARY KEY,
    from_member_id INT NOT NULL,
    from_login_name VARCHAR(255) NOT NULL,
    to_member_id INT NOT NULL,
    to_login_name VARCHAR(255) NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    amount DECIMAL NOT NULL,
    status SMALLINT, 
    created_at TIMESTAMP DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP DEFAULT timezone('UTC', now())
);

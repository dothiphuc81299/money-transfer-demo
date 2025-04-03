CREATE TABLE bank_account (
    id BIGSERIAL PRIMARY KEY,             
    bank_code VARCHAR(150) NOT NULL,      
    account_no VARCHAR(200) NOT NULL,     
    balance DECIMAL(18, 2) DEFAULT 0,     -
    outstanding_balance DECIMAL(18, 2) DEFAULT 0, 
    status SMALLINT NOT NULL DEFAULT 1,   
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()), 
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now())  
);

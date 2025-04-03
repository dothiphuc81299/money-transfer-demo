CREATE TABLE "user" (
    id BIGSERIAL PRIMARY KEY, 
    login_name  VARCHAR(150) UNIQUE NOT NULL,  
    password VARCHAR(255) NOT NULL,         
    status TINYINT NOT NULL DEFAULT 1,      
    created_at TIMESTAMPZ DEFAULT timezone('UTC', now()), 
    updated_at TIMESTAMPZ DEFAULT timezone('UTC', now())
);

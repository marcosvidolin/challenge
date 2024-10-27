CREATE SCHEMA IF NOT EXISTS challenge;

CREATE TABLE IF NOT EXISTS challenge.users (
    id BIGINT PRIMARY KEY,              
    first_name VARCHAR(100),           
    last_name VARCHAR(100),
    email_address VARCHAR(255) UNIQUE, 
    created_at TIMESTAMPTZ,           
    deleted_at TIMESTAMPTZ,          
    merged_at TIMESTAMPTZ,          
    parent_user_id BIGINT,         
    CONSTRAINT fk_parent_user
        FOREIGN KEY (parent_user_id)
        REFERENCES challenge.users(id)
        ON DELETE SET NULL
);

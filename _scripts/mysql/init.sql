CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY UNIQUE,
    document_number VARCHAR(11) NOT NULL UNIQUE,
    created_at TIMESTAMP
);
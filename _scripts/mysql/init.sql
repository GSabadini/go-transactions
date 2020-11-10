CREATE TABLE accounts (
    id VARCHAR(36) PRIMARY KEY UNIQUE,
    document_number VARCHAR(50) NOT NULL UNIQUE,
    available_credit_limit INTEGER NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE operations (
    id VARCHAR(36) PRIMARY KEY UNIQUE,
    description VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL
);

CREATE TABLE transactions (
    id VARCHAR(36) PRIMARY KEY UNIQUE,
    account_id VARCHAR(36) NOT NULL,
    operation_id VARCHAR(36) NOT NULL,
    amount INTEGER NOT NULL,
    balance INTEGER NOT NULL,
    created_at TIMESTAMP,

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (operation_id) REFERENCES operations(id)
);

INSERT
    INTO
        `operations` (`id`, `description`, `type`)
    VALUES
        ('fd426041-0648-40f6-9d04-5284295c5095', 'COMPRA A VISTA', 'DEBIT'),
        ('b03dcb59-006f-472f-a8f1-58651990dea6', 'COMPRA PARCELADA', 'DEBIT'),
        ('3f973e5b-cb9f-475c-b27d-8f855a0b90b0', 'SAQUE', 'DEBIT'),
        ('976f88ea-eb2f-4325-a106-26f9cb35810d', 'PAGAMENTO', 'CREDIT');
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
        ('1', 'COMPRA A VISTA', 'DEBIT'),
        ('2', 'COMPRA PARCELADA', 'DEBIT'),
        ('3', 'SAQUE', 'DEBIT'),
        ('4', 'PAGAMENTO', 'CREDIT');
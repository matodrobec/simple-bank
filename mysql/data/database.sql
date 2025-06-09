-- Active: 1742649766239@@127.0.0.1@3306@bank
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL,
    currency VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE entries (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    account_id BIGINT,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    from_account_id BIGINT ,
    to_account_id BIGINT ,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE entries
ADD CONSTRAINT fk_entries_account FOREIGN KEY (account_id)
REFERENCES accounts (id)
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE transfers
ADD CONSTRAINT fk_transfers_from FOREIGN KEY (from_account_id)
REFERENCES accounts (id)
ON DELETE SET NULL
ON UPDATE CASCADE;

ALTER TABLE transfers
ADD CONSTRAINT fk_transfers_to FOREIGN KEY (to_account_id)
REFERENCES accounts (id)
ON DELETE SET NULL
ON UPDATE CASCADE;

-- MySQL does not support COMMENT ON COLUMN, so use inline comments or a COMMENT clause in table creation:
ALTER TABLE transfers MODIFY COLUMN amount BIGINT NOT NULL COMMENT 'must be positive';
ALTER TABLE entries MODIFY COLUMN amount BIGINT NOT NULL COMMENT 'can be negative or positive';

CREATE INDEX idx_accounts_owner ON accounts(owner);
CREATE INDEX idx_entries_account_id ON entries(account_id);
CREATE INDEX idx_transfers_from ON transfers(from_account_id);
CREATE INDEX idx_transfers_to ON transfers(to_account_id);
CREATE INDEX idx_transfers_from_to ON transfers(from_account_id, to_account_id);

INSERT INTO accounts VALUES (null, "test", 100, "usd", now());
INSERT INTO accounts VALUES (null, "test2", 100, "usd", now());
INSERT INTO accounts VALUES (null, "test3", 100, "usd", now());


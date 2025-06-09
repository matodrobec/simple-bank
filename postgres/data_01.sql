
CREATE TABLE accounts (
	id bigserial PRIMARY KEY,
	owner varchar NOT NULL,
	balance bigint NOT NULL,
	currency varchar NOT NULL,
	created_at timestamptz DEFAULT now()
);

CREATE TABLE entries (
	id bigserial PRIMARY KEY,
	accounts_id bigint NOT NULL,
	amount bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE transfers (
	id bigserial PRIMARY KEY,
	from_accounts_id bigint NOT NULL,
	to_accounts_id bigint NOT NULL,
	amount bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE entries ADD FOREIGN KEY (accounts_id) REFERENCES accounts (id)
ON DELETE SET NULL ON UPDATE CASCADE;
COMMENT ON COLUMN transfers.amount IS E'must be positive';
COMMENT ON COLUMN entries.amount IS E'can be negative or positive';


ALTER TABLE transfers ADD FOREIGN KEY (from_accounts_id) REFERENCES accounts (id) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE transfers ADD FOREIGN KEY (to_accounts_id) REFERENCES accounts (id) ON DELETE SET NULL ON UPDATE CASCADE;

CREATE INDEX ON accounts(owner);
CREATE INDEX ON entries(accounts_id);
CREATE INDEX ON transfers(from_accounts_id);
CREATE INDEX ON transfers(to_accounts_id);
CREATE INDEX ON transfers(from_accounts_id, to_accounts_id);



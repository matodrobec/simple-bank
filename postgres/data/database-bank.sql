-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler version: 1.1.5
-- PostgreSQL version: 17.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: bank | type: DATABASE --
-- DROP DATABASE IF EXISTS bank;
CREATE DATABASE bank;
-- ddl-end --


-- object: public.accounts | type: TABLE --
-- DROP TABLE IF EXISTS public.accounts CASCADE;
CREATE TABLE public.accounts (
	id bigserial NOT NULL,
	owner varchar NOT NULL,
	balance bigint NOT NULL,
	currency varchar NOT NULL,
	created_at timestamptz DEFAULT now(),
	CONSTRAINT accounts_pk PRIMARY KEY (id)
);
-- ddl-end --
ALTER TABLE public.accounts OWNER TO postgres;
-- ddl-end --

-- object: public.entries | type: TABLE --
-- DROP TABLE IF EXISTS public.entries CASCADE;
CREATE TABLE public.entries (
	id bigserial NOT NULL,
	accounts_id bigint,
	amount bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT entries_pk PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN public.entries.amount IS E'can be negative or positive';
-- ddl-end --
ALTER TABLE public.entries OWNER TO postgres;
-- ddl-end --

-- object: accounts_fk | type: CONSTRAINT --
-- ALTER TABLE public.entries DROP CONSTRAINT IF EXISTS accounts_fk CASCADE;
ALTER TABLE public.entries ADD CONSTRAINT accounts_fk FOREIGN KEY (accounts_id)
REFERENCES public.accounts (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: public.transfers | type: TABLE --
-- DROP TABLE IF EXISTS public.transfers CASCADE;
CREATE TABLE public.transfers (
	id bigserial NOT NULL,
	amount bigint NOT NULL,
	to_id_accounts bigint,
	created_at timestamptz NOT NULL DEFAULT now(),
	from_accounts_id bigint,
	CONSTRAINT transfers_pk PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN public.transfers.amount IS E'must be positive';
-- ddl-end --
ALTER TABLE public.transfers OWNER TO postgres;
-- ddl-end --

-- object: form_accounts_fk | type: CONSTRAINT --
-- ALTER TABLE public.transfers DROP CONSTRAINT IF EXISTS form_accounts_fk CASCADE;
ALTER TABLE public.transfers ADD CONSTRAINT form_accounts_fk FOREIGN KEY (from_accounts_id)
REFERENCES public.accounts (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: to_accounts_fk | type: CONSTRAINT --
-- ALTER TABLE public.transfers DROP CONSTRAINT IF EXISTS to_accounts_fk CASCADE;
ALTER TABLE public.transfers ADD CONSTRAINT to_accounts_fk FOREIGN KEY (to_id_accounts)
REFERENCES public.accounts (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: owner_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.owner_idx CASCADE;
CREATE INDEX owner_idx ON public.accounts
USING btree
(
	owner
);
-- ddl-end --

-- object: account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.account_id_idx CASCADE;
CREATE INDEX account_id_idx ON public.entries
USING btree
(
	accounts_id
);
-- ddl-end --

-- object: from_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.from_account_id_idx CASCADE;
CREATE INDEX from_account_id_idx ON public.transfers
USING btree
(
	from_accounts_id
);
-- ddl-end --

-- object: to_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.to_account_id_idx CASCADE;
CREATE INDEX to_account_id_idx ON public.transfers
USING btree
(
	to_id_accounts
);
-- ddl-end --

-- object: from_to_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.from_to_account_id_idx CASCADE;
CREATE INDEX from_to_account_id_idx ON public.transfers
USING btree
(
	from_accounts_id,
	to_id_accounts
);
-- ddl-end --



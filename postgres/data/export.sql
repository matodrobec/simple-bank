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
CREATE DATABASE bank
	ENCODING = 'UTF8'
	LC_COLLATE = 'en_US.utf8'
	LC_CTYPE = 'en_US.utf8'
	TABLESPACE = pg_default
	OWNER = postgres;
-- ddl-end --


-- object: public.schema_migrations | type: TABLE --
-- DROP TABLE IF EXISTS public.schema_migrations CASCADE;
CREATE TABLE public.schema_migrations (
	version bigint NOT NULL,
	dirty boolean NOT NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);
-- ddl-end --
ALTER TABLE public.schema_migrations OWNER TO postgres;
-- ddl-end --

-- object: public.accounts_id_seq | type: SEQUENCE --
-- DROP SEQUENCE IF EXISTS public.accounts_id_seq CASCADE;
CREATE SEQUENCE public.accounts_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

-- ddl-end --
ALTER SEQUENCE public.accounts_id_seq OWNER TO postgres;
-- ddl-end --

-- object: public.accounts | type: TABLE --
-- DROP TABLE IF EXISTS public.accounts CASCADE;
CREATE TABLE public.accounts (
	id bigint NOT NULL DEFAULT nextval('public.accounts_id_seq'::regclass),
	owner character varying NOT NULL,
	balance bigint NOT NULL,
	currency character varying NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	CONSTRAINT accounts_pkey PRIMARY KEY (id)
);
-- ddl-end --
ALTER TABLE public.accounts OWNER TO postgres;
-- ddl-end --

-- object: public.entries_id_seq | type: SEQUENCE --
-- DROP SEQUENCE IF EXISTS public.entries_id_seq CASCADE;
CREATE SEQUENCE public.entries_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

-- ddl-end --
ALTER SEQUENCE public.entries_id_seq OWNER TO postgres;
-- ddl-end --

-- object: public.entries | type: TABLE --
-- DROP TABLE IF EXISTS public.entries CASCADE;
CREATE TABLE public.entries (
	id bigint NOT NULL DEFAULT nextval('public.entries_id_seq'::regclass),
	account_id bigint NOT NULL,
	amount bigint NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	CONSTRAINT entries_pkey PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN public.entries.amount IS E'can be negative or positive';
-- ddl-end --
ALTER TABLE public.entries OWNER TO postgres;
-- ddl-end --

-- object: public.transfers_id_seq | type: SEQUENCE --
-- DROP SEQUENCE IF EXISTS public.transfers_id_seq CASCADE;
CREATE SEQUENCE public.transfers_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

-- ddl-end --
ALTER SEQUENCE public.transfers_id_seq OWNER TO postgres;
-- ddl-end --

-- object: public.transfers | type: TABLE --
-- DROP TABLE IF EXISTS public.transfers CASCADE;
CREATE TABLE public.transfers (
	id bigint NOT NULL DEFAULT nextval('public.transfers_id_seq'::regclass),
	from_account_id bigint NOT NULL,
	to_account_id bigint NOT NULL,
	amount bigint NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	CONSTRAINT transfers_pkey PRIMARY KEY (id)
);
-- ddl-end --
COMMENT ON COLUMN public.transfers.amount IS E'must be positive';
-- ddl-end --
ALTER TABLE public.transfers OWNER TO postgres;
-- ddl-end --

-- object: accounts_owner_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.accounts_owner_idx CASCADE;
CREATE INDEX accounts_owner_idx ON public.accounts
USING btree
(
	owner
)
WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: entries_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.entries_account_id_idx CASCADE;
CREATE INDEX entries_account_id_idx ON public.entries
USING btree
(
	account_id
)
WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: transfers_from_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.transfers_from_account_id_idx CASCADE;
CREATE INDEX transfers_from_account_id_idx ON public.transfers
USING btree
(
	from_account_id
)
WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: transfers_to_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.transfers_to_account_id_idx CASCADE;
CREATE INDEX transfers_to_account_id_idx ON public.transfers
USING btree
(
	to_account_id
)
WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: transfers_from_account_id_to_account_id_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.transfers_from_account_id_to_account_id_idx CASCADE;
CREATE INDEX transfers_from_account_id_to_account_id_idx ON public.transfers
USING btree
(
	from_account_id,
	to_account_id
)
WITH (FILLFACTOR = 90);
-- ddl-end --

-- object: public.test | type: TABLE --
-- DROP TABLE IF EXISTS public.test CASCADE;
CREATE TABLE public.test (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	note text NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT now(),
	CONSTRAINT test_pkey PRIMARY KEY (id)
);
-- ddl-end --
ALTER TABLE public.test OWNER TO postgres;
-- ddl-end --

-- object: public.users | type: TABLE --
-- DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users (
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar,
	password_changed_at timestamptz DEFAULT 0001-01-01 00:00:00Z,
	created_ad timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT users_pk PRIMARY KEY (username),
	CONSTRAINT email UNIQUE (email)
);
-- ddl-end --
ALTER TABLE public.users OWNER TO postgres;
-- ddl-end --

-- object: entries_account_id_fkey | type: CONSTRAINT --
-- ALTER TABLE public.entries DROP CONSTRAINT IF EXISTS entries_account_id_fkey CASCADE;
ALTER TABLE public.entries ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY (account_id)
REFERENCES public.accounts (id) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: transfers_from_account_id_fkey | type: CONSTRAINT --
-- ALTER TABLE public.transfers DROP CONSTRAINT IF EXISTS transfers_from_account_id_fkey CASCADE;
ALTER TABLE public.transfers ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY (from_account_id)
REFERENCES public.accounts (id) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: transfers_to_account_id_fkey | type: CONSTRAINT --
-- ALTER TABLE public.transfers DROP CONSTRAINT IF EXISTS transfers_to_account_id_fkey CASCADE;
ALTER TABLE public.transfers ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY (to_account_id)
REFERENCES public.accounts (id) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --



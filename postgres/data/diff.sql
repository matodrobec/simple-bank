-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 1.1.5
-- Diff date: 2025-05-19 11:49:00
-- Source model: bank
-- Database: bank
-- PostgreSQL version: 17.0

-- [ Diff summary ]
-- Dropped objects: 0
-- Created objects: 3
-- Changed objects: 0

SET search_path=public,pg_catalog;
-- ddl-end --


-- [ Created objects ] --
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



-- [ Created constraints ] --
-- object: account_owner_currency | type: CONSTRAINT --
-- ALTER TABLE public.accounts DROP CONSTRAINT IF EXISTS account_owner_currency CASCADE;
ALTER TABLE public.accounts ADD CONSTRAINT account_owner_currency UNIQUE (owner,currency);
-- ddl-end --



-- [ Created foreign keys ] --
-- object: accounts_owner_fkey | type: CONSTRAINT --
-- ALTER TABLE public.accounts DROP CONSTRAINT IF EXISTS accounts_owner_fkey CASCADE;
ALTER TABLE public.accounts ADD CONSTRAINT accounts_owner_fkey FOREIGN KEY (owner)
REFERENCES public.users (username) MATCH SIMPLE
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


-- Diff code generated with pgModeler (PostgreSQL Database Modeler)
-- pgModeler version: 1.1.5
-- Diff date: 2025-06-01 19:56:27
-- Source model: bank
-- Database: bank
-- PostgreSQL version: 17.0

-- [ Diff summary ]
-- Dropped objects: 2
-- Created objects: 3
-- Changed objects: 0

SET search_path=public,pg_catalog;
-- ddl-end --


-- [ Dropped objects ] --
ALTER TABLE public.sessions DROP CONSTRAINT IF EXISTS sessions_username_fkey CASCADE;
-- ddl-end --
DROP TABLE IF EXISTS public.sessions CASCADE;
-- ddl-end --


-- [ Created objects ] --
-- object: is_email_verified | type: COLUMN --
-- ALTER TABLE public.users DROP COLUMN IF EXISTS is_email_verified CASCADE;
ALTER TABLE public.users ADD COLUMN is_email_verified bool NOT NULL DEFAULT false;
-- ddl-end --


-- object: public.verify_emails | type: TABLE --
-- DROP TABLE IF EXISTS public.verify_emails CASCADE;
CREATE TABLE public.verify_emails (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	email varchar NOT NULL,
	secret_code varchar NOT NULL,
	is_used boolean NOT NULL DEFAULT false,
	CONSTRAINT verify_emails_pk PRIMARY KEY (id)
);
-- ddl-end --
ALTER TABLE public.verify_emails OWNER TO postgres;
-- ddl-end --



-- [ Created foreign keys ] --
-- object: verify_emails_username_fkey | type: CONSTRAINT --
-- ALTER TABLE public.verify_emails DROP CONSTRAINT IF EXISTS verify_emails_username_fkey CASCADE;
ALTER TABLE public.verify_emails ADD CONSTRAINT verify_emails_username_fkey FOREIGN KEY (username)
REFERENCES public.users (username) MATCH SIMPLE
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --


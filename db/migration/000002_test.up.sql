CREATE TABLE test (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  note TEXT NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

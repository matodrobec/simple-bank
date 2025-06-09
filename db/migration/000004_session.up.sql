CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
	"username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT false,
	"expires_at" timestamptz NOT NULL,
	"created_at" timestamptz NOT NULL DEFAULT (now())
);
ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE SET NULL ON UPDATE CASCADE;
CREATE INDEX idx_expires_at ON sessions(expires_at);

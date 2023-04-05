CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar  NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_blocked" boolean NOT NULL DEFAULT false,
  expires_it timestamptz NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");


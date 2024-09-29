CREATE TABLE IF NOT EXISTS "channels" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "description" text,
  "created_by" integer NOT NULL,
  "last_modified_by" integer NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz
);

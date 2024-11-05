CREATE TABLE IF NOT EXISTS "plans" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "description" text,
  "created_by" integer NOT NULL,
  "last_modified_by" integer NOT NULL,
  "is_published" boolean NOT NULL DEFAULT false,
  "public" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz
);

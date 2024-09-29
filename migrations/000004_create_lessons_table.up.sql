CREATE TABLE IF NOT EXISTS "lessons" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "created_by" integer NOT NULL,
  "last_modified_by" integer NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz
);

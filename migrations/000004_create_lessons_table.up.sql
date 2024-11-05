CREATE TABLE IF NOT EXISTS "lessons" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "created_by" varchar(24) NOT NULL,
  "last_modified_by" varchar(24) NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz
);

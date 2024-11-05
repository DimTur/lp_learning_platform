CREATE TABLE IF NOT EXISTS "pages_abstractpages" (
  "id" SERIAL PRIMARY KEY,
  "lesson_id" integer,
  "created_by" varchar(24) NOT NULL,
  "last_modified_by" varchar(24) NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  "content_type" text NOT NULL CHECK (content_type IN ('pdf', 'video', 'image', 'question'))
);

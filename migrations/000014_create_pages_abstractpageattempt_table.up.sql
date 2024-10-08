CREATE TABLE IF NOT EXISTS "pages_abstractpageattempt" (
  "id" SERIAL PRIMARY KEY,
  "lesson_attempt_id" integer,
  "content_type" text NOT NULL CHECK (content_type IN ('pdf', 'video', 'image', 'question')),
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  CONSTRAINT fk_lesson_attempt FOREIGN KEY ("lesson_attempt_id") REFERENCES "attempt_lessonattempt" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "question_abstractquestionattempt" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  "question_type" text NOT NULL CHECK (question_type IN ('multichoice', 'short_answer')),
  "page_attempt_id" integer UNIQUE,
  "is_successful" boolean NOT NULL DEFAULT false,
  CONSTRAINT fk_page_attempt FOREIGN KEY ("page_attempt_id") REFERENCES "pages_abstractpageattempt" ("id") ON DELETE CASCADE
);

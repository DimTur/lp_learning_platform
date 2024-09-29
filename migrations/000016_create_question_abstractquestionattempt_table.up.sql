CREATE TABLE IF NOT EXISTS "question_abstractquestionattempt" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  "content_type" text NOT NULL CHECK (content_type IN ('multichoice', 'short_answer')),
  "page_attempt_id" integer UNIQUE,
  "is_successfull" boolean NOT NULL DEFAULT false,
  CONSTRAINT fk_page_attempt FOREIGN KEY ("page_attempt_id") REFERENCES "pages_abstractpageattempt" ("id") ON DELETE CASCADE
);

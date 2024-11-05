CREATE TABLE IF NOT EXISTS "question_questionpageattempt" (
  "id" SERIAL PRIMARY KEY,
  "page_id" integer,
  "page_attempt_id" integer UNIQUE,
  "user_answer" varchar(8),
  CONSTRAINT fk_page FOREIGN KEY ("page_id") REFERENCES "question_questionpage" ("id") ON DELETE CASCADE,
  CONSTRAINT fk_page_attempt FOREIGN KEY ("page_attempt_id") REFERENCES "question_abstractquestionattempt" ("id") ON DELETE CASCADE
);

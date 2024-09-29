CREATE TABLE IF NOT EXISTS "question_abstractquestion" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  "content_type" text NOT NULL CHECK (content_type IN ('multichoice', 'short_answer'))
);

CREATE TABLE IF NOT EXISTS "question_questionpage" (
  "id" SERIAL PRIMARY KEY,
  "abstractpage_id" integer UNIQUE,
  "question_id" integer UNIQUE,
  CONSTRAINT fk_abstractpage FOREIGN KEY ("abstractpage_id") REFERENCES "pages_abstractpages" ("id") ON DELETE CASCADE,
  CONSTRAINT fk_question FOREIGN KEY ("question_id") REFERENCES "question_abstractquestion" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "question_multichoicequestion" (
  "id" SERIAL PRIMARY KEY,
  "question_abstractquestion_id" integer UNIQUE,
  "question" text,
  "option_a" varchar(512) NOT NULL,
  "option_b" varchar(512) NOT NULL,
  "option_c" varchar(512) NOT NULL,
  "option_d" varchar(512) NOT NULL,
  "option_e" varchar(512) NOT NULL,
  "answer" varchar(1) NOT NULL,
  CONSTRAINT fk_question_abstractquestion FOREIGN KEY ("question_abstractquestion_id") REFERENCES "question_abstractquestion" ("id") ON DELETE CASCADE
);

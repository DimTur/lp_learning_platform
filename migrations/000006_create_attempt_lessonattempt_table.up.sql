CREATE TABLE IF NOT EXISTS "attempt_lessonattempt" (
  "id" SERIAL PRIMARY KEY,
  "lesson_id" integer,
  "plan_id" integer,
  "channel_id" integer,
  "start_time" timestamptz DEFAULT (now()),
  "end_time" timestamptz,
  "user_id" integer NOT NULL,
  "last_modified_by" integer NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "modified" timestamptz,
  "is_complete" boolean NOT NULL DEFAULT false,
  "is_successfull" boolean NOT NULL DEFAULT false,
  "percentage_score" integer DEFAULT 0 CHECK (percentage_score BETWEEN 0 AND 100)
);

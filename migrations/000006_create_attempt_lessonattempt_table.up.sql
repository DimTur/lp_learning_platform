CREATE TABLE IF NOT EXISTS "attempt_lessonattempt" (
  "id" SERIAL PRIMARY KEY,
  "lesson_id" integer NOT NULL,
  "plan_id" integer NOT NULL,
  "channel_id" integer NOT NULL,
  "start_time" timestamptz DEFAULT (now()),
  "end_time" timestamptz,
  "user_id" integer NOT NULL,
  "is_complete" boolean NOT NULL DEFAULT false,
  "is_successful" boolean NOT NULL DEFAULT false,
  "percentage_score" integer DEFAULT 0 CHECK (percentage_score BETWEEN 0 AND 100)
);

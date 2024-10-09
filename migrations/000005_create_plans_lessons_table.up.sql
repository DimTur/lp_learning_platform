CREATE TABLE IF NOT EXISTS "plans_lessons" (
  "id" SERIAL PRIMARY KEY,
  "plan_id" integer NOT NULL,
  "lesson_id" integer NOT NULL,
  CONSTRAINT fk_plan FOREIGN KEY ("plan_id") REFERENCES "plans" ("id") ON DELETE CASCADE,
  CONSTRAINT fk_lesson FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("id") ON DELETE CASCADE
);

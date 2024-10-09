CREATE TABLE IF NOT EXISTS "channels_plans" (
  "id" SERIAL PRIMARY KEY,
  "channel_id" integer NOT NULL,
  "plan_id" integer NOT NULL,
  CONSTRAINT fk_channel FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE,
  CONSTRAINT fk_plan FOREIGN KEY ("plan_id") REFERENCES "plans" ("id") ON DELETE CASCADE
);

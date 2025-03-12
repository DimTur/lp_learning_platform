CREATE TABLE IF NOT EXISTS "shared_channels_learninggroups" (
  "id" SERIAL PRIMARY KEY,
  "channel_id" integer NOT NULL,
  "learning_group_id" varchar(24) NOT NULL,
  "created_by" varchar(24) NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  CONSTRAINT fk_channel FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE
);

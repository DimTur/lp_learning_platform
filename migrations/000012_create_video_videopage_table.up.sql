CREATE TABLE IF NOT EXISTS "video_videopage" (
  "id" SERIAL PRIMARY KEY,
  "abstractpage_id" integer UNIQUE,
  "video_file_url" varchar(512) NOT NULL,
  "video_name" varchar(255),
  CONSTRAINT fk_abstractpage FOREIGN KEY ("abstractpage_id") REFERENCES "pages_abstractpages" ("id") ON DELETE CASCADE
);

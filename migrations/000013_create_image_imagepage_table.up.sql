CREATE TABLE IF NOT EXISTS "image_imagepage" (
  "id" SERIAL PRIMARY KEY,
  "abstractpage_id" integer UNIQUE,
  "image_file_url" varchar(512) NOT NULL,
  "image_name" varchar(255),
  CONSTRAINT fk_abstractpage FOREIGN KEY ("abstractpage_id") REFERENCES "pages_abstractpages" ("id") ON DELETE CASCADE
);

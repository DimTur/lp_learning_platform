CREATE TABLE IF NOT EXISTS "pdf_pdfpage" (
  "id" SERIAL PRIMARY KEY,
  "abstractpage_id" integer UNIQUE,
  "pdf_file_url" varchar(512) NOT NULL,
  "pdf_name" varchar(255),
  CONSTRAINT fk_abstractpage FOREIGN KEY ("abstractpage_id") REFERENCES "pages_abstractpages" ("id") ON DELETE CASCADE
);

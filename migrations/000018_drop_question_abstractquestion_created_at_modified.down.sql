ALTER TABLE "question_abstractquestion"
ADD COLUMN "created_at" timestamptz DEFAULT (now()),
ADD COLUMN "modified" timestamptz;
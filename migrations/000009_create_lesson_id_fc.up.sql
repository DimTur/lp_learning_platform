ALTER TABLE "attempt_lessonattempt" 
ADD CONSTRAINT fk_lesson FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("id") ON DELETE CASCADE;

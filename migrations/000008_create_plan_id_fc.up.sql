ALTER TABLE "attempt_lessonattempt" 
ADD CONSTRAINT fk_plan FOREIGN KEY ("plan_id") REFERENCES "plans" ("id") ON DELETE CASCADE;

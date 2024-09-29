ALTER TABLE "attempt_lessonattempt" 
ADD CONSTRAINT fk_channel FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON DELETE CASCADE;

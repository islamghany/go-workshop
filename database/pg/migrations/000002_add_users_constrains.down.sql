ALTER TABLE users
DROP CONSTRAINT IF EXISTS user_unique_email;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS user_age_check;

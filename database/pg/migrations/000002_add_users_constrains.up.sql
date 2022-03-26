ALTER TABLE users
ADD CONSTRAINT user_unique_email UNIQUE (id);


ALTER TABLE users
ADD CONSTRAINT user_age_check CHECK(age >= 18 AND age <= 60);

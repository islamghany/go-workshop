-- The hash column will contain a SHA-256 hash of the activation token.
--  It’s important to emphasize that we will only store a hash of the activation
--   token in our database — not the activation token itself.

-- We want to hash the token before storing it for the same reason that we bcrypt 
-- a user’s password — it provides an extra layer of protection if the database 
-- is ever compromised or leaked. Because our activation token is going to be a 
-- high-entropy random string (128 bits) — rather than something low entropy 
-- like a typical user password — it is sufficient to use a fast algorithm 
-- like SHA-256 to create the hash, instead of a slow algorithm like bcrypt.


-- the scope column will denote what purpose the token can be used for. 
-- Later in the book we’ll also need to create and store authentication tokens,
--  and most of the code and storage requirements for these is exactly the same as 
--  for our activation tokens. So instead of creating separate 
--  tables (and the code to interact with them), we’ll store them in one table 
--  with a value in the scope column to restrict the purpose that the token can 
--  be used for.
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
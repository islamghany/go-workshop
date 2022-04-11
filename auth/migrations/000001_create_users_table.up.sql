CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL, --  comparisons against the data are always case-insensitive
    password_hash bytea NOT NULL,  -- the type bytea (binary string). In this column we’ll store a one-way hash of the user’s password generated using bcrypt — not the plaintext password itself.  
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1 --  prevent race conditions when updating user records
);
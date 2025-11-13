-- Enable pgcrypto for bcrypt functions (crypt, gen_salt)
CREATE EXTENSION IF NOT EXISTS pgcrypto;
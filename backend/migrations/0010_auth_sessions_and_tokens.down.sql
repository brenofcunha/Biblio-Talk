DROP TABLE IF EXISTS auth_action_tokens;
DROP TABLE IF EXISTS auth_sessions;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;

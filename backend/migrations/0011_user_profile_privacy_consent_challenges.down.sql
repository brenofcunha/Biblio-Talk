DROP TABLE IF EXISTS annual_challenges;
DROP TABLE IF EXISTS user_consents;
ALTER TABLE users
    DROP COLUMN IF EXISTS privacy_settings,
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS banner_url,
    DROP COLUMN IF EXISTS display_name;

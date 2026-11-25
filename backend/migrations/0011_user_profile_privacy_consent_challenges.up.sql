ALTER TABLE users
    ADD COLUMN display_name text NOT NULL DEFAULT '',
    ADD COLUMN banner_url text,
    ADD COLUMN location text,
    ADD COLUMN privacy_settings jsonb NOT NULL DEFAULT
        '{"profile":"public","diary":"public","progress":"public","shelf":"public","notes":"private","comments":"public","lists":"public","stats":"public","favorites":"public"}'::jsonb;

CREATE TABLE user_consents (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type text NOT NULL CHECK (consent_type IN ('terms_of_service', 'privacy_policy', 'marketing')),
    version text NOT NULL,
    granted boolean NOT NULL,
    recorded_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX user_consents_user_recorded_idx ON user_consents (user_id, recorded_at DESC);

CREATE TABLE annual_challenges (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year int NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    target_books int NOT NULL CHECK (target_books BETWEEN 1 AND 1000),
    is_public boolean NOT NULL DEFAULT false,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, year)
);

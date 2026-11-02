CREATE TABLE IF NOT EXISTS log_likes (
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id     uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, log_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id       uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    comment_text text NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS shelves (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id    uuid NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status     text NOT NULL CHECK (status IN ('tbr','reading','read','abandoned')),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, book_id)
);

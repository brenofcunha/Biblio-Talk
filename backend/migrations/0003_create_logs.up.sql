CREATE TABLE IF NOT EXISTS logs (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id           uuid NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    rating            numeric(2,1) NOT NULL CHECK (rating >= 0.5 AND rating <= 5.0 AND (rating * 2) = trunc(rating * 2)),
    review_text       text,
    read_date         date,
    is_reread         boolean NOT NULL DEFAULT false,
    contains_spoilers boolean NOT NULL DEFAULT false,
    likes_count       int NOT NULL DEFAULT 0,
    created_at        timestamptz NOT NULL DEFAULT now()
);

-- Uma avaliação por livro quando não há data de leitura; por livro+data quando há
CREATE UNIQUE INDEX IF NOT EXISTS logs_user_book_null_read_idx ON logs (user_id, book_id) WHERE read_date IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS logs_user_book_read_idx ON logs (user_id, book_id, read_date) WHERE read_date IS NOT NULL;

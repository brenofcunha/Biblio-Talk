CREATE TABLE IF NOT EXISTS lists (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       text NOT NULL,
    description text,
    is_private  boolean NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS list_items (
    list_id  uuid NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    book_id  uuid NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    position int NOT NULL,
    notes    text,
    PRIMARY KEY (list_id, book_id),
    UNIQUE (list_id, position)
);

CREATE TABLE IF NOT EXISTS books (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    external_api_id text,
    title           text NOT NULL,
    authors         text[] NOT NULL DEFAULT '{}',
    isbn            text,
    publisher       text,
    publish_year    int,
    page_count      int,
    cover_url       text,
    genres          text[] NOT NULL DEFAULT '{}',
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS books_external_api_id_idx ON books (external_api_id);
CREATE INDEX IF NOT EXISTS books_title_lower_idx ON books (lower(title));

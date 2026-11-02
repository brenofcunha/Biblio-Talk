-- Cache das estatísticas anuais recalculadas de forma assíncrona pelo worker pool
CREATE TABLE IF NOT EXISTS user_year_stats (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year       int NOT NULL,
    payload    jsonb NOT NULL DEFAULT '{}',
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, year)
);

-- Índices de performance para feed, diário, página do livro e comentários
CREATE INDEX IF NOT EXISTS logs_user_created_idx ON logs (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS logs_book_created_idx ON logs (book_id, created_at DESC);
CREATE INDEX IF NOT EXISTS logs_created_include_id_idx ON logs (created_at DESC) INCLUDE (id);
CREATE INDEX IF NOT EXISTS comments_log_created_idx ON comments (log_id, created_at);
CREATE INDEX IF NOT EXISTS log_likes_log_idx ON log_likes (log_id, user_id);
CREATE INDEX IF NOT EXISTS list_items_list_pos_idx ON list_items (list_id, position);
CREATE INDEX IF NOT EXISTS lists_user_created_idx ON lists (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS books_isbn_idx ON books (isbn) WHERE isbn IS NOT NULL;

-- Busca por título com similaridade (pg_trgm)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS books_title_trgm_idx ON books USING gin (title gin_trgm_ops);

CREATE TABLE reading_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id uuid NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status text NOT NULL CHECK (status IN ('want_to_read', 'reading', 'read', 'abandoned')),
    current_page int NOT NULL DEFAULT 0 CHECK (current_page >= 0),
    total_pages int CHECK (total_pages IS NULL OR total_pages > 0),
    started_at date,
    finished_at date,
    is_reread boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (total_pages IS NULL OR current_page <= total_pages)
);

CREATE INDEX reading_sessions_user_updated_idx ON reading_sessions (user_id, updated_at DESC);
CREATE INDEX reading_sessions_book_user_idx ON reading_sessions (book_id, user_id);

CREATE TABLE reading_session_history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id uuid NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type text NOT NULL,
    from_status text,
    to_status text,
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX reading_session_history_session_created_idx
    ON reading_session_history (session_id, created_at DESC);

CREATE TABLE reading_progress (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id uuid NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    progress_date date NOT NULL,
    current_page int NOT NULL CHECK (current_page >= 0),
    pages_read int NOT NULL CHECK (pages_read >= 0),
    note text,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (session_id, progress_date)
);

CREATE INDEX reading_progress_user_date_idx ON reading_progress (user_id, progress_date DESC);

CREATE TABLE private_notes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id uuid REFERENCES books(id) ON DELETE CASCADE,
    reading_session_id uuid REFERENCES reading_sessions(id) ON DELETE CASCADE,
    content text NOT NULL CHECK (length(trim(content)) > 0),
    is_published boolean NOT NULL DEFAULT false,
    published_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((book_id IS NOT NULL)::int + (reading_session_id IS NOT NULL)::int = 1)
);

CREATE INDEX private_notes_user_updated_idx ON private_notes (user_id, updated_at DESC);

CREATE TABLE private_note_tags (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    note_id uuid NOT NULL REFERENCES private_notes(id) ON DELETE CASCADE,
    tag text NOT NULL CHECK (length(trim(tag)) BETWEEN 1 AND 50),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX private_note_tags_note_lower_idx
    ON private_note_tags (note_id, lower(tag));

CREATE TABLE private_note_quotes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    note_id uuid NOT NULL REFERENCES private_notes(id) ON DELETE CASCADE,
    quote_text text NOT NULL CHECK (length(trim(quote_text)) > 0),
    page int CHECK (page IS NULL OR page > 0),
    chapter text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE reviews (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id uuid NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    reading_session_id uuid REFERENCES reading_sessions(id) ON DELETE SET NULL,
    rating numeric(2,1) CHECK (rating IS NULL OR (rating >= 0.5 AND rating <= 5.0 AND (rating * 2) = trunc(rating * 2))),
    review_text text NOT NULL DEFAULT '',
    contains_spoilers boolean NOT NULL DEFAULT false,
    status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    allow_comments boolean NOT NULL DEFAULT true,
    allow_likes boolean NOT NULL DEFAULT true,
    published_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (status <> 'published' OR length(trim(review_text)) > 0),
    CHECK (status <> 'published' OR published_at IS NOT NULL)
);

CREATE INDEX reviews_user_status_updated_idx ON reviews (user_id, status, updated_at DESC);
CREATE INDEX reviews_book_published_idx ON reviews (book_id, published_at DESC) WHERE status = 'published';

CREATE TABLE review_history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id uuid NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action text NOT NULL,
    snapshot jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX review_history_review_created_idx ON review_history (review_id, created_at DESC);

CREATE TABLE reading_goals (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    goal_year int NOT NULL CHECK (goal_year BETWEEN 2000 AND 2100),
    goal_type text NOT NULL CHECK (goal_type IN ('monthly', 'annual')),
    goal_month int CHECK (goal_month IS NULL OR goal_month BETWEEN 1 AND 12),
    target_books int CHECK (target_books IS NULL OR target_books BETWEEN 1 AND 1000),
    target_pages int CHECK (target_pages IS NULL OR target_pages BETWEEN 1 AND 1000000),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((goal_type = 'monthly' AND goal_month IS NOT NULL) OR (goal_type = 'annual' AND goal_month IS NULL)),
    CHECK (target_books IS NOT NULL OR target_pages IS NOT NULL)
);

CREATE UNIQUE INDEX reading_goals_user_period_idx
    ON reading_goals (user_id, goal_year, goal_type, COALESCE(goal_month, 0));

CREATE TABLE reading_history_imports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status text NOT NULL CHECK (status IN ('completed', 'failed')),
    total_items int NOT NULL DEFAULT 0 CHECK (total_items >= 0),
    imported_items int NOT NULL DEFAULT 0 CHECK (imported_items >= 0),
    result jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz
);

CREATE INDEX reading_history_imports_user_created_idx
    ON reading_history_imports (user_id, created_at DESC);

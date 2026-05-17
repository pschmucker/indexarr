-- Rename tvdb_id column to tmdb_id in series table and add new tvdb_id column
ALTER TABLE series RENAME COLUMN tvdb_id TO tmdb_id;
ALTER TABLE series ADD COLUMN tvdb_id INTEGER;
ALTER TABLE series ADD COLUMN imdb_id_new TEXT;
ALTER TABLE series ADD COLUMN poster_new TEXT;

UPDATE series SET tvdb_id = (-1 * tmdb_id), imdb_id_new = imdb_id, poster_new = poster;
ALTER TABLE series DROP COLUMN imdb_id;
ALTER TABLE series DROP COLUMN poster;
ALTER TABLE series RENAME COLUMN imdb_id_new TO imdb_id;
ALTER TABLE series RENAME COLUMN poster_new TO poster;

CREATE UNIQUE INDEX IF NOT EXISTS idx_series_tvdb_id ON series(tvdb_id);


-- Add tvdb_tokens table for storing TVDB API bearer tokens
CREATE TABLE IF NOT EXISTS tvdb_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tvdb_tokens_singleton ON tvdb_tokens (id);

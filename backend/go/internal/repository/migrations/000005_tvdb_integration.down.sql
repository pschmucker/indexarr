DROP INDEX IF EXISTS idx_series_tvdb_id;


-- Remove tvdb_id column and rename tmdb_id column to tvdb_id in series table
ALTER TABLE series DROP COLUMN tvdb_id;
ALTER TABLE series RENAME COLUMN tmdb_id TO tvdb_id;


-- Remove tvdb_tokens table
DROP TABLE IF EXISTS tvdb_tokens;

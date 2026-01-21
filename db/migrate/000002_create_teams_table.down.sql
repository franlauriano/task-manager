-- Drop indexes
DROP INDEX IF EXISTS idx_teams_deleted_at;
DROP INDEX IF EXISTS idx_teams_uuid;

-- Drop teams table
DROP TABLE IF EXISTS teams;

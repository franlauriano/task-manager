-- Drop index
DROP INDEX IF EXISTS idx_tasks_team_id;

-- Remove team_id column from tasks table
ALTER TABLE tasks
DROP COLUMN IF EXISTS team_id;

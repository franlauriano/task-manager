-- Add team_id column to tasks table
ALTER TABLE tasks
ADD COLUMN team_id INTEGER REFERENCES teams(id);

-- Create index on team_id for performance
CREATE INDEX idx_tasks_team_id ON tasks(team_id);

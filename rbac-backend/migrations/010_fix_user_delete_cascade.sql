-- Fix foreign key constraints to allow user deletion with cascade
-- Need to disable foreign keys, recreate tables, then re-enable them

PRAGMA foreign_keys = OFF;

-- Step 1: Create backup of projects
CREATE TABLE projects_backup AS SELECT * FROM projects;

-- Step 2: Drop the old projects table and recreate with cascade
DROP TABLE projects;
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by TEXT NOT NULL,
    status TEXT,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 3: Copy data back
INSERT INTO projects SELECT * FROM projects_backup;
DROP TABLE projects_backup;

-- Step 4: Create backup of tasks and recreate with cascade
CREATE TABLE tasks_backup AS SELECT * FROM tasks;
DROP TABLE tasks;
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK(status IN ('TODO','IN_PROGRESS','REVIEW','DONE','ARCHIVED')) NOT NULL DEFAULT 'TODO',
    assignee TEXT,
    created_by TEXT NOT NULL,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 5: Recreate indexes
CREATE INDEX idx_tasks_project ON tasks(project_id);
CREATE INDEX idx_tasks_assignee ON tasks(assignee);
CREATE INDEX idx_tasks_status ON tasks(status);

-- Step 6: Copy data back
INSERT INTO tasks SELECT * FROM tasks_backup;
DROP TABLE tasks_backup;

-- Re-enable foreign keys
PRAGMA foreign_keys = ON;

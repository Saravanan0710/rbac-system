-- Add status column to projects table
ALTER TABLE projects ADD COLUMN status TEXT DEFAULT 'Active';

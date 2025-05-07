-- Drop task-related tables from the taskodex schema

-- Drop task history table
DROP TABLE IF EXISTS taskodex.task_history;

-- Drop task time entries table
DROP TABLE IF EXISTS taskodex.task_time_entries;

-- Drop task attachments table
DROP TABLE IF EXISTS taskodex.task_attachments;

-- Drop task comments table
DROP TABLE IF EXISTS taskodex.task_comments;

-- Drop task tags table
DROP TABLE IF EXISTS taskodex.task_tags;

-- Add task-related tables to the taskodex schema

-- Task tags table
CREATE TABLE IF NOT EXISTS taskodex.task_tags (
    task_id UUID NOT NULL,
    tag VARCHAR(50) NOT NULL,
    PRIMARY KEY (task_id, tag),
    CONSTRAINT fk_task_tags_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_tags_task_id ON taskodex.task_tags(task_id);
CREATE INDEX IF NOT EXISTS idx_task_tags_tag ON taskodex.task_tags(tag);

-- Task comments table
CREATE TABLE IF NOT EXISTS taskodex.task_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_task_comments_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_comments_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_comments_task_id ON taskodex.task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_user_id ON taskodex.task_comments(user_id);

-- Task attachments table
CREATE TABLE IF NOT EXISTS taskodex.task_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(1000) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_task_attachments_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_attachments_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_attachments_task_id ON taskodex.task_attachments(task_id);
CREATE INDEX IF NOT EXISTS idx_task_attachments_user_id ON taskodex.task_attachments(user_id);

-- Task time entries table
CREATE TABLE IF NOT EXISTS taskodex.task_time_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_task_time_entries_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_time_entries_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_time_entries_task_id ON taskodex.task_time_entries(task_id);
CREATE INDEX IF NOT EXISTS idx_task_time_entries_user_id ON taskodex.task_time_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_task_time_entries_start_time ON taskodex.task_time_entries(start_time);
CREATE INDEX IF NOT EXISTS idx_task_time_entries_end_time ON taskodex.task_time_entries(end_time);

-- Task history table for audit trail
CREATE TABLE IF NOT EXISTS taskodex.task_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_task_history_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_history_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_task_history_task_id ON taskodex.task_history(task_id);
CREATE INDEX IF NOT EXISTS idx_task_history_user_id ON taskodex.task_history(user_id);
CREATE INDEX IF NOT EXISTS idx_task_history_action ON taskodex.task_history(action);
CREATE INDEX IF NOT EXISTS idx_task_history_created_at ON taskodex.task_history(created_at);

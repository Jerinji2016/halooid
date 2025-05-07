-- Create task_time_entries table

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
    CONSTRAINT fk_time_entries_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_time_entries_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_time_entries_task_id ON taskodex.task_time_entries(task_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_user_id ON taskodex.task_time_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_start_time ON taskodex.task_time_entries(start_time);
CREATE INDEX IF NOT EXISTS idx_time_entries_end_time ON taskodex.task_time_entries(end_time);
CREATE INDEX IF NOT EXISTS idx_time_entries_running ON taskodex.task_time_entries(user_id, task_id) WHERE end_time IS NULL;

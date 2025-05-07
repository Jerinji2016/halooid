-- Create task_file_attachments table
CREATE TABLE IF NOT EXISTS taskodex.task_file_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    storage_path VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_file_attachments_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_file_attachments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_file_attachments_task_id ON taskodex.task_file_attachments(task_id);
CREATE INDEX IF NOT EXISTS idx_file_attachments_user_id ON taskodex.task_file_attachments(user_id);
CREATE INDEX IF NOT EXISTS idx_file_attachments_created_at ON taskodex.task_file_attachments(created_at);

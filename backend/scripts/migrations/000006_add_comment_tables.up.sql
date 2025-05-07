-- Create task_comments table
CREATE TABLE IF NOT EXISTS taskodex.task_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_comments_task FOREIGN KEY (task_id) REFERENCES taskodex.tasks(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create comment_mentions table
CREATE TABLE IF NOT EXISTS taskodex.comment_mentions (
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (comment_id, user_id),
    CONSTRAINT fk_mentions_comment FOREIGN KEY (comment_id) REFERENCES taskodex.task_comments(id) ON DELETE CASCADE,
    CONSTRAINT fk_mentions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_comments_task_id ON taskodex.task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON taskodex.task_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON taskodex.task_comments(created_at);
CREATE INDEX IF NOT EXISTS idx_mentions_user_id ON taskodex.comment_mentions(user_id);

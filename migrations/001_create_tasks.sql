CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'completed')),
    priority VARCHAR(50) NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

INSERT INTO tasks (id, title, description, status, priority, due_date, created_at, updated_at)
VALUES
    ('sample-1', 'Complete the project documentation', 'Write comprehensive README and setup instructions', 'active', 'high', NOW() + INTERVAL '3 days', NOW(), NOW()),
    ('sample-2', 'Review code changes', 'Review PR #123 and provide feedback', 'active', 'medium', NOW() + INTERVAL '1 day', NOW(), NOW()),
    ('sample-3', 'Setup CI/CD pipeline', 'Configure GitHub Actions for automated testing', 'completed', 'medium', NULL, NOW() - INTERVAL '2 days', NOW())
ON CONFLICT (id) DO NOTHING;
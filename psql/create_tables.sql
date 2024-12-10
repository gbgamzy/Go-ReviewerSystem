-- SQL script to create tables in the ReviewerDB database

-- Create the users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the tasks table
CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       status VARCHAR(50) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'In Progress', 'Approved', 'Cancelled')),
                       created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       required_approvals INT DEFAULT 3 NOT NULL,
                       current_approvals INT DEFAULT 0 NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the `updated_at` column automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create the task_approvers table
CREATE TABLE task_approvers (
                                id SERIAL PRIMARY KEY,
                                task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
                                approver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                approved_at TIMESTAMP,
                                UNIQUE (task_id, approver_id)
);

-- Create the task_approval_comments table
CREATE TABLE task_approval_comments (
                                        id SERIAL PRIMARY KEY,
                                        task_approver_id INT NOT NULL REFERENCES task_approvers(id) ON DELETE CASCADE,
                                        comment TEXT NOT NULL,
                                        commented_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimization
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_task_approvers_task_id ON task_approvers(task_id);
CREATE INDEX idx_task_approvers_approved_by ON task_approvers(approver_id);

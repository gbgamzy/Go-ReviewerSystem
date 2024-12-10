package repository

import (
	"ReviewerService/pkg/models"
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	DB *sql.DB
}

// Create a new task
func (r *TaskRepository) CreateTask(task *models.Task) (int, error) {
	query := `
		INSERT INTO tasks (title, description, status, created_by, required_approvals, current_approvals)
		VALUES ($1, $2, 'Pending', $3, $4, 0) RETURNING id;
	`
	err := r.DB.QueryRow(query, task.Title, task.Description, task.CreatedBy, task.RequiredApprovals).Scan(&task.ID)
	if err != nil {
		return 0, fmt.Errorf("error creating task: %v", err)
	}
	return task.ID, nil
}

// Mark a task as in-progress
func (r *TaskRepository) MarkTaskInProgress(taskID int) error {
	query := `UPDATE tasks SET status = 'In Progress' WHERE id = $1`
	_, err := r.DB.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("error marking task as in-progress: %v", err)
	}
	return nil
}

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

func (r *TaskRepository) GetTaskCreatorEmail(taskID int) (string, error) {
	var email string
	query := `
		SELECT u.email
		FROM tasks t
		JOIN users u ON t.created_by = u.id
		WHERE t.id = $1
	`
	err := r.DB.QueryRow(query, taskID).Scan(&email)
	return email, err
}

func (r *TaskRepository) GetAllTaskPartyEmails(taskID int) ([]string, error) {
	query := `
		SELECT u.email
		FROM users u
		JOIN task_approvers ta ON u.id = ta.approver_id
		WHERE ta.task_id = $1
		UNION
		SELECT u.email
		FROM users u
		JOIN tasks t ON u.id = t.created_by
		WHERE t.id = $1
	`
	rows, err := r.DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

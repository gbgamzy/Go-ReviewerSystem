package repository

import (
	"database/sql"
	"fmt"
)

type CommentRepository struct {
	DB *sql.DB
}

// Add a comment to a task
func (r *CommentRepository) AddComment(taskApproverID int, comment string) error {
	query := `INSERT INTO task_approval_comments (task_approver_id, comment) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, taskApproverID, comment)
	if err != nil {
		return fmt.Errorf("error adding comment to task approver %d: %v", taskApproverID, err)
	}
	return nil
}

// GetTaskApproverID retrieves the ID of the task approver record for a given task and approver.
func (r *CommentRepository) GetTaskApproverID(taskID, approverID int) (int, error) {
	var taskApproverID int
	query := `
		SELECT id
		FROM task_approvers
		WHERE task_id = $1 AND approver_id = $2
	`
	err := r.DB.QueryRow(query, taskID, approverID).Scan(&taskApproverID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("no task approver found for taskID: %d and approverID: %d", taskID, approverID)
	} else if err != nil {
		return 0, fmt.Errorf("error retrieving task approver: %v", err)
	}
	return taskApproverID, nil
}

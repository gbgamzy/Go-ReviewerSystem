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

package repository

import (
	"database/sql"
	"fmt"
)

type ApproverRepository struct {
	DB *sql.DB
}

// Assign approvers to a task
func (r *ApproverRepository) AssignApprovers(taskID int, approverIDs []int) error {
	query := `INSERT INTO task_approvers (task_id, approver_id) VALUES ($1, $2)`
	for _, approverID := range approverIDs {
		_, err := r.DB.Exec(query, taskID, approverID)
		if err != nil {
			return fmt.Errorf("error assigning approver %d to task %d: %v", approverID, taskID, err)
		}
	}
	return nil
}

// Check if a user is an approver for a task
func (r *ApproverRepository) IsApproverForTask(taskID, approverID int) (bool, error) {
	query := `SELECT COUNT(*) FROM task_approvers WHERE task_id = $1 AND approver_id = $2`
	var count int
	err := r.DB.QueryRow(query, taskID, approverID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking approver: %v", err)
	}
	return count > 0, nil
}

package repository

import (
	"database/sql"
	"fmt"
)

type ApprovalRepository struct {
	DB *sql.DB
}

// Approve a task
func (r *ApprovalRepository) ApproveTask(taskID, approverID int) error {
	// Update task_approvers to set approved_at
	query := `
		UPDATE task_approvers
		SET approved_at = CURRENT_TIMESTAMP
		WHERE task_id = $1 AND approver_id = $2
	`
	_, err := r.DB.Exec(query, taskID, approverID)
	if err != nil {
		return fmt.Errorf("error approving task %d by approver %d: %v", taskID, approverID, err)
	}

	// Increment the current_approvals count in tasks
	updateQuery := `
		UPDATE tasks
		SET current_approvals = current_approvals + 1
		WHERE id = $1
	`
	_, err = r.DB.Exec(updateQuery, taskID)
	if err != nil {
		return fmt.Errorf("error updating approval count for task %d: %v", taskID, err)
	}

	// Check if the task has received all required approvals
	checkQuery := `
		UPDATE tasks
		SET status = 'Approved'
		WHERE id = $1 AND current_approvals >= required_approvals
	`
	_, err = r.DB.Exec(checkQuery, taskID)
	if err != nil {
		return fmt.Errorf("error marking task %d as approved: %v", taskID, err)
	}

	return nil
}

func (r *ApprovalRepository) IsTaskFullyApproved(taskID int) (bool, error) {
	var requiredApprovals, currentApprovals int
	query := `SELECT required_approvals, current_approvals FROM tasks WHERE id = $1`
	err := r.DB.QueryRow(query, taskID).Scan(&requiredApprovals, &currentApprovals)
	if err != nil {
		return false, err
	}
	return currentApprovals >= requiredApprovals, nil
}

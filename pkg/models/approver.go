package models

import "time"

type TaskApprover struct {
	ID         int       `json:"id"`          // Unique identifier for the approver record
	TaskID     int       `json:"task_id"`     // ID of the associated task
	ApproverID int       `json:"approver_id"` // ID of the user assigned as the approver
	ApprovedAt time.Time `json:"approved_at"` // Timestamp when the approval was given (nullable)
}

type TaskApprovalComment struct {
	ID             int       `json:"id"`               // Unique identifier for the comment
	TaskApproverID int       `json:"task_approver_id"` // ID of the associated TaskApprover record
	Comment        string    `json:"comment"`          // The comment text
	CommentedAt    time.Time `json:"commented_at"`     // Timestamp when the comment was made
}

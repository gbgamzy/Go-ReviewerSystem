package models

import "time"

type Task struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	CreatedBy         int       `json:"created_by"`
	RequiredApprovals int       `json:"required_approvals"`
	CurrentApprovals  int       `json:"current_approvals"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

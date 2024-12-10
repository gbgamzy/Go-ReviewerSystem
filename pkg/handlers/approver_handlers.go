package handlers

import (
	"ReviewerService/pkg/repository"
	"ReviewerService/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApproverHandler struct {
	ApproverRepo *repository.ApproverRepository
	ApprovalRepo *repository.ApprovalRepository
	UserRepo     *repository.UserRepository
	TaskRepo     *repository.TaskRepository
}

// Assign Approvers to a Task
func (h *ApproverHandler) AssignApproversHandler(w http.ResponseWriter, r *http.Request) {
	// Parse task_id from URL
	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse approver IDs from request body
	var req struct {
		ApproverIDs []int `json:"approver_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Assign approvers to the task
	if err := h.ApproverRepo.AssignApprovers(taskID, req.ApproverIDs); err != nil {
		http.Error(w, "Failed to assign approvers", http.StatusInternalServerError)
		return
	}

	// Fetch and send email notifications to approvers
	for _, approverID := range req.ApproverIDs {
		go func(approverID int) {
			email, err := h.UserRepo.GetUserEmailByID(approverID)
			if err != nil {
				log.Printf("Failed to fetch email for approver ID %d: %v", approverID, err)
				return
			}

			err = utils.SendEmail(email, "Task Assignment Notification", fmt.Sprintf("You have been assigned to task ID %d.", taskID))
			if err != nil {
				log.Printf("Failed to send email to %s: %v", email, err)
			}
		}(approverID)
	}

	// Response to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Approvers assigned successfully."})
}

// Approve a Task
func (h *ApproverHandler) ApproveTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ApproverID int `json:"approver_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the approver is valid for the given task
	valid, err := h.ApproverRepo.IsApproverForTask(taskID, req.ApproverID)
	if err != nil {
		http.Error(w, "Failed to validate approver", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Approver is not assigned to this task", http.StatusForbidden)
		return
	}

	// Approve the task
	err = h.ApprovalRepo.ApproveTask(taskID, req.ApproverID)
	if err != nil {
		http.Error(w, "Failed to approve task", http.StatusInternalServerError)
		return
	}

	// Notify the task creator about the approval
	taskCreatorEmail, err := h.TaskRepo.GetTaskCreatorEmail(taskID)
	if err == nil {
		go utils.SendEmail(taskCreatorEmail, "Task Update", fmt.Sprintf("Approver ID %d has approved your task ID %d.", req.ApproverID, taskID))
	}

	// Check if all approvals are complete
	allApproved, err := h.ApprovalRepo.IsTaskFullyApproved(taskID)
	if err != nil {
		http.Error(w, "Failed to check task approval status", http.StatusInternalServerError)
		return
	}

	if allApproved {
		// Notify all parties involved
		partyEmails, err := h.TaskRepo.GetAllTaskPartyEmails(taskID)
		if err == nil {
			for _, email := range partyEmails {
				go utils.SendEmail(email, "Task Fully Approved", fmt.Sprintf("Task ID %d has been fully approved.", taskID))
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task approved successfully."})
}

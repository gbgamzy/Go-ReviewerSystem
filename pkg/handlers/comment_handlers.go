package handlers

import (
	"ReviewerService/pkg/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	CommentRepo *repository.CommentRepository
}

func (h *CommentHandler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ApproverID int    `json:"approver_id"`
		Comment    string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve taskApproverID
	taskApproverID, err := h.CommentRepo.GetTaskApproverID(taskID, req.ApproverID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve task approver ID: %v", err), http.StatusBadRequest)
		return
	}

	// Add comment to the task
	if err := h.CommentRepo.AddComment(taskApproverID, req.Comment); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add comment: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully."})
}

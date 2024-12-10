package handlers

import (
	"ReviewerService/pkg/models"
	"ReviewerService/pkg/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	TaskRepo *repository.TaskRepository
}

// Create a Task
func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		CreatedBy         int    `json:"created_by"`
		RequiredApprovals int    `json:"required_approvals"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newTask := &models.Task{
		Title:             task.Title,
		Description:       task.Description,
		CreatedBy:         task.CreatedBy,
		RequiredApprovals: task.RequiredApprovals,
	}
	taskID, err := h.TaskRepo.CreateTask(newTask)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	newTask.ID = taskID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// Publish a Task
func (h *TaskHandler) PublishTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID, _ := strconv.Atoi(mux.Vars(r)["task_id"])

	if err := h.TaskRepo.MarkTaskInProgress(taskID); err != nil {
		http.Error(w, "Failed to publish task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task marked as In Progress."})
}

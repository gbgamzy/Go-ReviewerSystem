package handlers

import (
	"ReviewerService/pkg/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserRepo *repository.UserRepository
}

// Login or Signup a User
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := h.UserRepo.LoginOrSignup(req.Email, req.Name)
	if err != nil {
		http.Error(w, "Failed to login/signup user", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":    userID,
		"email": req.Email,
		"name":  req.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Get Users with Pagination and Search
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		search = ""
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	users, total, err := h.UserRepo.GetUsers(search, page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

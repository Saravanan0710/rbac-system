package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"rbac-backend/internal/auth"
	repositories "rbac-backend/internal/repository"

	"github.com/google/uuid"
)

type AdminHandler struct {
	UserRepo *repositories.UserRepository
}

func NewAdminHandler(repo *repositories.UserRepository) *AdminHandler {
	return &AdminHandler{UserRepo: repo}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}
	for _, u := range users {
		response = append(response, map[string]interface{}{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"role":  u.Role,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		payload.Name = payload.Email
	}

	if payload.Role == "" {
		payload.Role = "VIEWER"
	} else {
		// Ensure role is uppercase
		payload.Role = strings.ToUpper(payload.Role)
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	userID := uuid.New().String()

	err = h.UserRepo.CreateUserWithPassword(userID, payload.Name, payload.Email, hashedPassword, payload.Role)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    userID,
		"name":  payload.Name,
		"email": payload.Email,
		"role":  payload.Role,
	})
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.GetUserByID(id)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"user": map[string]interface{}{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
	}})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		ID string `json:"id"`
	}
	// try decode body, otherwise read query param
	_ = json.NewDecoder(r.Body).Decode(&payload)
	id := payload.ID
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := h.UserRepo.DeleteUser(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if payload.ID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	// Prepare updates map
	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if payload.Role != "" {
		// Ensure role is uppercase
		updates["role"] = strings.ToUpper(payload.Role)
	}

	// If no fields to update, return error
	if len(updates) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	if err := h.UserRepo.UpdateUser(payload.ID, updates); err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	// Fetch updated user and return it
	u, err := h.UserRepo.GetUserByID(payload.ID)
	if err != nil || u == nil {
		http.Error(w, "failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
	})
}

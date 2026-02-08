package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rbac-backend/internal/auth"
	dbrepo "rbac-backend/internal/db"
	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/utils"

	"github.com/google/uuid"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ✅ Method check
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// ✅ Content-Type
		w.Header().Set("Content-Type", "application/json")

		var req struct {
			Email    string
			Password string
		}

		json.NewDecoder(r.Body).Decode(&req)

		var userID, hash, role, name string
		err := db.QueryRow(
			"SELECT id, password_hash, role, name FROM users WHERE email=? AND is_active=1",
			req.Email,
		).Scan(&userID, &hash, &role, &name)

		if err != nil || auth.CheckPassword(hash, req.Password) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateJWT(userID, role, name)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// Signup registers a new user and returns a JWT token.
func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "Name, email and password are required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}

		userID := uuid.New().String()
		role := "VIEWER"

		_, err = db.Exec(
			`INSERT INTO users (id, name, email, password_hash, role, is_active) 
			 VALUES (?, ?, ?, ?, ?, 1)`,
			userID, req.Name, req.Email, hashedPassword, role,
		)
		if err != nil {
			// Most likely a unique constraint on email
			http.Error(w, "Could not create user", http.StatusBadRequest)
			return
		}

		token, err := auth.GenerateJWT(userID, role, req.Name)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func ViewEmployees(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value(middleware.RoleKey).(string)

		perms, _ := dbrepo.GetPermissionsByRole(database, role)
		fieldPerms := perms["employees"].Fields

		employee := map[string]interface{}{
			"id":     "E101",
			"name":   "Ravi",
			"salary": 90000,
		}

		response := utils.FilterFields(employee, fieldPerms)
		json.NewEncoder(w).Encode(response)
	})
}

func EditEmployees(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		role := r.Context().Value(middleware.RoleKey).(string)

		perms, _ := dbrepo.GetPermissionsByRole(database, role)
		fieldPerms := perms["employees"].Fields

		var input map[string]interface{}
		json.NewDecoder(r.Body).Decode(&input)

		allowedUpdate := utils.FilterEditableFields(input, fieldPerms)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"updated_fields": allowedUpdate,
		})
	})
}

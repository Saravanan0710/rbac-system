package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rbac-backend/internal/db"
)

// StatsHandler returns dashboard stats (counts).
type StatsHandler struct {
	DB *sql.DB
}

func NewStatsHandler(database *sql.DB) *StatsHandler {
	return &StatsHandler{DB: database}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var totalProjects, totalUsers, totalTasks int
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&totalProjects)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE is_active = 1`).Scan(&totalUsers)
	_ = h.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&totalTasks)

	roles, _ := db.ListRoles(h.DB)
	totalRoles := len(roles) + 1 // +1 for ADMIN which is not in role_permissions

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_projects": totalProjects,
		"total_users":    totalUsers,
		"total_tasks":    totalTasks,
		"total_roles":   totalRoles,
	})
}

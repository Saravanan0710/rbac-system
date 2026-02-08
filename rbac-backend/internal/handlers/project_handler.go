package handlers

import (
	"encoding/json"
	"net/http"
	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
	repositories "rbac-backend/internal/repository"
	"rbac-backend/internal/utils"

	"github.com/google/uuid"
)

type ProjectHandler struct {
	Repo     *repositories.ProjectRepository
	UserRepo *repositories.UserRepository
}

func NewProjectHandler(repo *repositories.ProjectRepository, userRepo *repositories.UserRepository) *ProjectHandler {
	return &ProjectHandler{Repo: repo, UserRepo: userRepo}
}
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	// Role-level restriction: VIEWER and EDITOR cannot create projects
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == rbac.RoleViewer || role == rbac.RoleEditor {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	safe["id"] = uuid.New().String()
	safe["created_by"] = userID

	err := h.Repo.CreateProjectDynamic(safe)
	if err != nil {
		http.Error(w, "failed to create project", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(safe)
}

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	projects, err := h.Repo.GetProjects()
	if err != nil {
		http.Error(w, "failed to fetch projects", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}

	// Determine role and userID to apply project-level visibility rules
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	for _, p := range projects {
		// Only ADMIN sees all projects. MANAGER, EDITOR and VIEWER see projects they're assigned to.
		if role != rbac.RoleAdmin {
			allowed := false
			for _, uid := range p.AssignedEmployees {
				if uid == userID {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}
		row := map[string]interface{}{
			"id":                 p.ID,
			"name":               p.Name,
			"description":        p.Description,
			"status":             p.Status,
			"created_by":         p.CreatedBy,
			"assigned_employees": p.AssignedEmployees,
		}

		filtered := utils.FilterFields(row, tablePerm.Fields)

		// Ensure assigned_employees is included when the table view is allowed.
		// Field-level config may omit this field; include it here so frontend can
		// show assignees when the user can view projects at all.
		if len(p.AssignedEmployees) > 0 && tablePerm.View {
			filtered["assigned_employees"] = p.AssignedEmployees
		}
		response = append(response, filtered)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}
	p, err := h.Repo.GetProjectByID(id)
	if err != nil {
		http.Error(w, "failed to fetch project", http.StatusInternalServerError)
		return
	}
	if p == nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	// Only ADMIN can view any project; MANAGER, EDITOR and VIEWER must be assigned.
	if role != rbac.RoleAdmin {
		allowed := false
		for _, uid := range p.AssignedEmployees {
			if uid == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}
	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	row := map[string]interface{}{
		"id":                 p.ID,
		"name":               p.Name,
		"description":        p.Description,
		"status":             p.Status,
		"created_by":         p.CreatedBy,
		"assigned_employees": p.AssignedEmployees,
	}
	filtered := utils.FilterFields(row, tablePerm.Fields)
	if len(p.AssignedEmployees) > 0 && tablePerm.View {
		filtered["assigned_employees"] = p.AssignedEmployees
	}
	json.NewEncoder(w).Encode(filtered)
}

func (h *ProjectHandler) GetProjectEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		http.Error(w, "project_id required", http.StatusBadRequest)
		return
	}
	p, err := h.Repo.GetProjectByID(projectID)
	if err != nil || p == nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	if role != rbac.RoleAdmin {
		allowed := false
		for _, uid := range p.AssignedEmployees {
			if uid == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}
	var employees []map[string]interface{}
	for _, uid := range p.AssignedEmployees {
		u, err := h.UserRepo.GetUserByID(uid)
		if err != nil || u == nil {
			continue
		}
		employees = append(employees, map[string]interface{}{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"role":  u.Role,
		})
	}
	json.NewEncoder(w).Encode(employees)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// ✅ ID must be read before filtering
	id, ok := incoming["id"].(string)
	if !ok || id == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	// Remove ID from editable fields input
	delete(incoming, "id")

	// Filter only editable fields
	safeData := utils.FilterEditableFields(incoming, tablePerm.Fields)

	if len(safeData) == 0 {
		http.Error(w, "no editable fields", http.StatusForbidden)
		return
	}

	// Put ID back for repository WHERE clause
	safeData["id"] = id

	err := h.Repo.UpdateProjectDynamic(safeData)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	return // ✅ VERY IMPORTANT

}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "project id required", http.StatusBadRequest)
		return
	}

	err := h.Repo.DeleteProject(id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "project deleted",
	})
}

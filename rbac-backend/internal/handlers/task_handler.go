package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"rbac-backend/internal/middleware"
	"rbac-backend/internal/models"
	"rbac-backend/internal/rbac"
	repositories "rbac-backend/internal/repository"
	"rbac-backend/internal/utils"
)

type TaskHandler struct {
	Repo        *repositories.TaskRepository
	ProjectRepo *repositories.ProjectRepository
}

func NewTaskHandler(repo *repositories.TaskRepository, projectRepo *repositories.ProjectRepository) *TaskHandler {
	return &TaskHandler{Repo: repo, ProjectRepo: projectRepo}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Extract required fields from original incoming (before filtering)
	// These are allowed on creation even though they're not editable afterward
	pid, ok := incoming["project_id"].(string)
	if !ok || pid == "" {
		http.Error(w, "project_id required", http.StatusBadRequest)
		return
	}
	title, ok := incoming["title"].(string)
	if !ok || title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	// Filter optional/editable fields
	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	t := models.Task{
		ID:        uuid.New().String(),
		ProjectID: pid,
		Title:     title,
		CreatedBy: userID,
		Status:    "TODO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Editors and Viewers are not allowed to create tasks.
	if role == rbac.RoleViewer || role == rbac.RoleEditor {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// MANAGER may create tasks only for projects they are assigned to.
	if role == rbac.RoleManager {
		p, err := h.ProjectRepo.GetProjectByID(pid)
		if err != nil || p == nil {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
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

	if desc, ok := safe["description"].(string); ok {
		t.Description = desc
	}
	if s, ok := safe["status"].(string); ok && s != "" {
		// Normalize status to uppercase format expected by database
		statusMap := map[string]string{
			"todo":        "TODO",
			"in-progress": "IN_PROGRESS",
			"in_progress": "IN_PROGRESS",
			"review":      "REVIEW",
			"done":        "DONE",
			"archived":    "ARCHIVED",
		}
		if normalized, exists := statusMap[strings.ToLower(s)]; exists {
			t.Status = normalized
		} else {
			// If not in map, try uppercasing as fallback
			t.Status = strings.ToUpper(s)
		}
	}
	// accept assignees as []string or single string
	if arr, ok := safe["assignees"].([]interface{}); ok {
		for _, it := range arr {
			if sid, ok := it.(string); ok {
				t.Assignees = append(t.Assignees, sid)
			}
		}
	} else if a, ok := safe["assignee"].(string); ok && a != "" {
		t.Assignees = append(t.Assignees, a)
	}

	if err := h.Repo.CreateTask(t); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	projectID := r.URL.Query().Get("project_id")
	assignee := r.URL.Query().Get("assignee")

	var tasks []models.Task
	var err error

	if projectID != "" {
		tasks, err = h.Repo.ListTasksByProject(projectID)
	} else if assignee != "" {
		tasks, err = h.Repo.ListTasksByAssignee(assignee)
	} else {
		tasks, err = h.Repo.ListTasks()
	}

	if err != nil {
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	var out []map[string]interface{}
	for _, t := range tasks {
		// Project-level visibility: only ADMIN can bypass project assignment.
		// MANAGER, EDITOR and VIEWER must be assigned to the task's project to view tasks.
		if role != rbac.RoleAdmin {
			p, err := h.ProjectRepo.GetProjectByID(t.ProjectID)
			if err != nil || p == nil {
				continue
			}
			assigned := false
			for _, uid := range p.AssignedEmployees {
				if uid == userID {
					assigned = true
					break
				}
			}
			if !assigned {
				continue
			}
		}

		// Within assigned projects: VIEWER has a narrower view (only created or assigned tasks),
		// EDITOR can view all tasks within their assigned projects.
		if role == rbac.RoleViewer {
			allowed := t.CreatedBy == userID
			if !allowed {
				for _, a := range t.Assignees {
					if a == userID {
						allowed = true
						break
					}
				}
			}
			if !allowed {
				continue
			}
		}

		row := map[string]interface{}{
			"id":           t.ID,
			"project_id":   t.ProjectID,
			"title":        t.Title,
			"description":  t.Description,
			"status":       t.Status,
			"assignees":    t.Assignees,
			"created_by":   t.CreatedBy,
			"started_at":   t.StartedAt,
			"completed_at": t.CompletedAt,
			"created_at":   t.CreatedAt,
			"updated_at":   t.UpdatedAt,
		}

		filtered := utils.FilterFields(row, tablePerm.Fields)
		out = append(out, filtered)
	}

	json.NewEncoder(w).Encode(out)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	t, err := h.Repo.GetTaskByID(id)
	if err != nil {
		http.Error(w, "failed to fetch task", http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	// Require project assignment for non-ADMIN; MANAGER, EDITOR and VIEWER must be assigned.
	if role != rbac.RoleAdmin {
		p, err := h.ProjectRepo.GetProjectByID(t.ProjectID)
		if err != nil || p == nil {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		assigned := false
		for _, uid := range p.AssignedEmployees {
			if uid == userID {
				assigned = true
				break
			}
		}
		if !assigned {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// If VIEWER, additionally ensure they either created the task or are assigned to it.
		if role == rbac.RoleViewer {
			allowed := t.CreatedBy == userID
			if !allowed {
				for _, a := range t.Assignees {
					if a == userID {
						allowed = true
						break
					}
				}
			}
			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
	}

	row := map[string]interface{}{
		"id":           t.ID,
		"project_id":   t.ProjectID,
		"title":        t.Title,
		"description":  t.Description,
		"status":       t.Status,
		"assignees":    t.Assignees,
		"created_by":   t.CreatedBy,
		"started_at":   t.StartedAt,
		"completed_at": t.CompletedAt,
		"created_at":   t.CreatedAt,
		"updated_at":   t.UpdatedAt,
	}

	json.NewEncoder(w).Encode(utils.FilterFields(row, tablePerm.Fields))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tablePerm := r.Context().Value(middleware.TablePermKey).(models.ResourcePermission)

	var incoming map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	idVal, ok := incoming["id"].(string)
	if !ok || idVal == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	// Fetch existing
	existing, err := h.Repo.GetTaskByID(idVal)
	if err != nil || existing == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	safe := utils.FilterEditableFields(incoming, tablePerm.Fields)

	role, _ := r.Context().Value(middleware.RoleKey).(string)
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	// Detect whether incoming payload included a status field
	_, hasStatus := incoming["status"]

	// If the current user is an assignee on the task, allow them to update
	// only the `status` field regardless of project assignment.
	isAssignee := false
	for _, a := range existing.Assignees {
		if a == userID {
			isAssignee = true
			break
		}
	}

	if isAssignee {
		if hasStatus {
			safe["status"] = incoming["status"]
		}
		if len(safe) == 0 {
			http.Error(w, "no editable fields", http.StatusForbidden)
			return
		}
		for k := range safe {
			if k != "status" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
	} else if role != rbac.RoleAdmin {
		// Non-admins must be assigned to the project to modify tasks in it.
		p, err := h.ProjectRepo.GetProjectByID(existing.ProjectID)
		if err != nil || p == nil {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}
		assigned := false
		for _, uid := range p.AssignedEmployees {
			if uid == userID {
				assigned = true
				break
			}
		}
		if !assigned {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// VIEWER and EDITOR can only update `status` (and VIEWER must be creator if not assignee).
		if role == rbac.RoleViewer {
			allowed := existing.CreatedBy == userID
			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if hasStatus {
				safe["status"] = incoming["status"]
			}
			if len(safe) == 0 {
				http.Error(w, "no editable fields", http.StatusForbidden)
				return
			}
			for k := range safe {
				if k != "status" {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}
		}

		if role == rbac.RoleEditor {
			if hasStatus {
				safe["status"] = incoming["status"]
			}
			if len(safe) == 0 {
				http.Error(w, "no editable fields", http.StatusForbidden)
				return
			}
			for k := range safe {
				if k != "status" {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}
		}
		// MANAGER may update fields according to table permissions (no extra restriction here).
	}

	if title, ok := safe["title"].(string); ok {
		existing.Title = title
	}
	if desc, ok := safe["description"].(string); ok {
		existing.Description = desc
	}
	if status, ok := safe["status"].(string); ok {
		// Normalize status to uppercase format expected by database
		statusMap := map[string]string{
			"todo":        "TODO",
			"in-progress": "IN_PROGRESS",
			"in_progress": "IN_PROGRESS",
			"review":      "REVIEW",
			"done":        "DONE",
			"archived":    "ARCHIVED",
		}
		if normalized, exists := statusMap[strings.ToLower(status)]; exists {
			status = normalized
		} else {
			status = strings.ToUpper(status)
		}

		existing.Status = status
		if status == "IN_PROGRESS" && existing.StartedAt == nil {
			now := time.Now()
			existing.StartedAt = &now
		}
		if (status == "DONE" || status == "ARCHIVED") && existing.CompletedAt == nil {
			now := time.Now()
			existing.CompletedAt = &now
		}
	}
	if arr, ok := safe["assignees"].([]interface{}); ok {
		existing.Assignees = nil
		for _, it := range arr {
			if sid, ok := it.(string); ok {
				existing.Assignees = append(existing.Assignees, sid)
			}
		}
	} else if a, ok := safe["assignee"].(string); ok && a != "" {
		// maintain backward compatibility for single assignee field
		existing.Assignees = []string{a}
	}

	existing.UpdatedAt = time.Now()

	if err := h.Repo.UpdateTask(*existing); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		ID        string   `json:"id"`
		Assignee  string   `json:"assignee"`
		Assignees []string `json:"assignees"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if payload.ID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	// If assignees array provided, set as the assignees list
	if len(payload.Assignees) > 0 {
		t, err := h.Repo.GetTaskByID(payload.ID)
		if err != nil || t == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		t.Assignees = payload.Assignees
		if err := h.Repo.UpdateTask(*t); err != nil {
			http.Error(w, "assign failed", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "assigned"})
		return
	}

	if payload.Assignee == "" {
		http.Error(w, "assignee required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.AssignTask(payload.ID, payload.Assignee); err != nil {
		http.Error(w, "assign failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteTask(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "task deleted"})
}

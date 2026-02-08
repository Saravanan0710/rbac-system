package repositories

import (
	"database/sql"
	"errors"
	"rbac-backend/internal/models"
	"strings"
)

type ProjectRepository struct {
	DB *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}
func (r *ProjectRepository) CreateProject(project models.Project) error {

	_, err := r.DB.Exec(
		`INSERT INTO projects (id, name, description, created_by)
		 VALUES (?, ?, ?, ?)`,
		project.ID,
		project.Name,
		project.Description,
		project.CreatedBy,
	)

	return err
}
func (r *ProjectRepository) CreateProjectDynamic(data map[string]interface{}) error {

	if len(data) == 0 {
		return errors.New("no data provided")
	}

	columns := []string{}
	placeholders := []string{}
	args := []interface{}{}

	for col, val := range data {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		args = append(args, val)
	}

	_ = "INSERT INTO projects (" +
		strings.Join(columns, ",") +
		") VALUES (" +
		strings.Join(placeholders, ",") +
		")"

	// Extract assigned_employees if present
	var assignments []string
	if a, ok := data["assigned_employees"]; ok {
		// copy and remove from data used for projects table
		switch v := a.(type) {
		case []string:
			assignments = v
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					assignments = append(assignments, s)
				}
			}
		}
		// remove assignments from insert columns
		// rebuild columns/args without assigned_employees
		newCols := []string{}
		newPlaceholders := []string{}
		newArgs := []interface{}{}
		for i, col := range columns {
			if col == "assigned_employees" {
				continue
			}
			newCols = append(newCols, col)
			newPlaceholders = append(newPlaceholders, placeholders[i])
			newArgs = append(newArgs, args[i])
		}
		columns = newCols
		placeholders = newPlaceholders
		args = newArgs
	}

	_, err := r.DB.Exec("INSERT INTO projects ("+strings.Join(columns, ",")+") VALUES ("+strings.Join(placeholders, ",")+")", args...)
	if err != nil {
		return err
	}

	// If assignments exist, insert into project_assignments
	if len(assignments) > 0 {
		// expect project id present in data
		pid, _ := data["id"].(string)
		if pid == "" {
			return errors.New("project id required for assignments")
		}
		tx, err := r.DB.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT OR REPLACE INTO project_assignments (project_id, user_id) VALUES (?, ?)`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		for _, uid := range assignments {
			if _, err := stmt.Exec(pid, uid); err != nil {
				tx.Rollback()
				return err
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectRepository) GetProjectByID(id string) (*models.Project, error) {
	row := r.DB.QueryRow(`SELECT id, name, description, status, created_by FROM projects WHERE id = ?`, id)
	var p models.Project
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	assignRows, err := r.DB.Query(`SELECT user_id FROM project_assignments WHERE project_id = ?`, p.ID)
	if err == nil {
		defer assignRows.Close()
		for assignRows.Next() {
			var uid string
			if err := assignRows.Scan(&uid); err == nil {
				p.AssignedEmployees = append(p.AssignedEmployees, uid)
			}
		}
	}
	return &p, nil
}

func (r *ProjectRepository) GetProjects() ([]models.Project, error) {

	rows, err := r.DB.Query(`SELECT id, name, description, status, created_by FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedBy)
		if err != nil {
			return nil, err
		}

		// fetch assigned employees for this project
		assignRows, err := r.DB.Query(`SELECT user_id FROM project_assignments WHERE project_id = ?`, p.ID)
		if err == nil {
			defer assignRows.Close()
			var assigns []string
			for assignRows.Next() {
				var uid string
				if err := assignRows.Scan(&uid); err == nil {
					assigns = append(assigns, uid)
				}
			}
			p.AssignedEmployees = assigns
		}

		projects = append(projects, p)
	}

	return projects, nil
}
func (r *ProjectRepository) UpdateProjectDynamic(data map[string]interface{}) error {

	idVal, ok := data["id"]
	if !ok {
		return errors.New("id required for update")
	}
	id := idVal.(string)

	delete(data, "id") // do not update ID

	// Extract assigned_employees if present (need separate handling)
	var assignments []string
	assignmentsProvided := false
	if a, ok := data["assigned_employees"]; ok {
		assignmentsProvided = true
		switch v := a.(type) {
		case []string:
			assignments = v
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					assignments = append(assignments, s)
				}
			}
		}
		// Remove from data since it's handled separately
		delete(data, "assigned_employees")
	}

	// Update projects table with regular fields
	if len(data) > 0 {
		query := "UPDATE projects SET "
		args := []interface{}{}
		i := 0

		for field, value := range data {
			if i > 0 {
				query += ", "
			}
			query += field + "=?"
			args = append(args, value)
			i++
		}

		query += " WHERE id=?"
		args = append(args, id)

		_, err := r.DB.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	// Handle assigned_employees separately if provided (even if empty)
	if assignmentsProvided {
		// Validate that all user IDs exist
		for _, uid := range assignments {
			exists, err := r.UserExists(uid)
			if err != nil || !exists {
				return errors.New("invalid user id: " + uid)
			}
		}

		tx, err := r.DB.Begin()
		if err != nil {
			return err
		}

		// Delete existing assignments
		if _, err := tx.Exec(`DELETE FROM project_assignments WHERE project_id=?`, id); err != nil {
			tx.Rollback()
			return err
		}

		// Insert new assignments (skip if empty array)
		if len(assignments) > 0 {
			stmt, err := tx.Prepare(`INSERT INTO project_assignments (project_id, user_id) VALUES (?, ?)`)
			if err != nil {
				tx.Rollback()
				return err
			}
			defer stmt.Close()

			for _, uid := range assignments {
				if _, err := stmt.Exec(id, uid); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) UserExists(userID string) (bool, error) {
	row := r.DB.QueryRow(`SELECT id FROM users WHERE id = ?`, userID)
	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ProjectRepository) DeleteProject(id string) error {

	_, err := r.DB.Exec(`DELETE FROM projects WHERE id=?`, id)
	return err
}

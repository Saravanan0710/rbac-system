package repositories

import (
	"database/sql"
	"strings"

	"rbac-backend/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, name, email, password_hash, role, is_active)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive,
	)
	return err
}

func (r *UserRepository) CreateUserWithPassword(id, name, email, passwordHash, role string) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		id, name, email, passwordHash, role,
	)
	return err
}

func (r *UserRepository) ListUsers() ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, email, role, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var isActive int
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &isActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		u.IsActive = isActive == 1
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	return r.ListUsers()
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, email, role, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`, id)

	var u models.User
	var isActive int
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &isActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.IsActive = isActive == 1
	return &u, nil
}

func (r *UserRepository) UpdateUser(id string, updates map[string]interface{}) error {
	// Build the SET clause dynamically based on the updates provided
	var setClauses []string
	args := []interface{}{}

	// Add requested fields to update
	if name, ok := updates["name"]; ok {
		setClauses = append(setClauses, "name = ?")
		args = append(args, name)
	}
	if role, ok := updates["role"]; ok {
		setClauses = append(setClauses, "role = ?")
		args = append(args, role)
	}
	if email, ok := updates["email"]; ok {
		setClauses = append(setClauses, "email = ?")
		args = append(args, email)
	}

	// Always update the timestamp
	setClauses = append(setClauses, "updated_at = datetime('now')")

	// Build the final SET clause
	setClause := strings.Join(setClauses, ", ")

	// Append ID to args for WHERE clause
	args = append(args, id)

	_, err := r.DB.Exec(
		"UPDATE users SET "+setClause+" WHERE id = ?",
		args...,
	)
	return err
}

func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.DB.Exec(`UPDATE users SET is_active = 0 WHERE id = ?`, id)
	return err
}

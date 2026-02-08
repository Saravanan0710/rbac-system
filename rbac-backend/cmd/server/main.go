package main

import (
	"log"
	"net/http"

	"rbac-backend/internal/db"
	"rbac-backend/internal/handlers"
	"rbac-backend/internal/middleware"
	repositories "rbac-backend/internal/repository"
)

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	database := db.Connect()
	defer database.Close()

	// Root Greeting
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to RBAC System Backend 🚀"))
	})

	// AUTH ROUTES
	http.Handle("/login", handlers.Login(database))

	// ⭐ CREATE PROJECT HANDLER
	projectRepo := repositories.NewProjectRepository(database)
	userRepo := repositories.NewUserRepository(database)
	projectHandler := handlers.NewProjectHandler(projectRepo, userRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)
	statsHandler := handlers.NewStatsHandler(database)

	// TASK HANDLER
	taskRepo := repositories.NewTaskRepository(database)
	taskHandler := handlers.NewTaskHandler(taskRepo, projectRepo)

	// ⭐ PROJECT ROUTES

	http.Handle(
		"/projects/create",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "create",
				http.HandlerFunc(projectHandler.CreateProject),
			),
		),
	)

	http.Handle(
		"/projects",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "view",
				http.HandlerFunc(projectHandler.GetProjects),
			),
		),
	)
	http.Handle(
		"/projects/get",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "view",
				http.HandlerFunc(projectHandler.GetProject),
			),
		),
	)
	http.Handle(
		"/projects/employees",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "view",
				http.HandlerFunc(projectHandler.GetProjectEmployees),
			),
		),
	)
	http.Handle(
		"/admin/stats",
		middleware.AuthMiddleware(http.HandlerFunc(statsHandler.GetStats)),
	)

	http.Handle(
		"/projects/update",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "edit",
				http.HandlerFunc(projectHandler.UpdateProject),
			),
		),
	)

	http.Handle(
		"/projects/delete",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "projects", "create",
				http.HandlerFunc(projectHandler.DeleteProject),
			),
		),
	)

	// TASK ROUTES
	http.Handle(
		"/tasks/create",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "create",
				http.HandlerFunc(taskHandler.CreateTask),
			),
		),
	)

	http.Handle(
		"/tasks",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "view",
				http.HandlerFunc(taskHandler.ListTasks),
			),
		),
	)

	http.Handle(
		"/tasks/get",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "view",
				http.HandlerFunc(taskHandler.GetTask),
			),
		),
	)

	http.Handle(
		"/tasks/update",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "edit",
				http.HandlerFunc(taskHandler.UpdateTask),
			),
		),
	)

	http.Handle(
		"/tasks/assign",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "edit",
				http.HandlerFunc(taskHandler.AssignTask),
			),
		),
	)

	http.Handle(
		"/tasks/delete",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "tasks", "delete",
				http.HandlerFunc(taskHandler.DeleteTask),
			),
		),
	)

	http.Handle(
		"/admin/create-user",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "create",
				http.HandlerFunc(adminHandler.CreateUser),
			),
		),
	)

	http.Handle(
		"/admin/delete-user",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "create",
				http.HandlerFunc(adminHandler.DeleteUser),
			),
		),
	)

	http.Handle(
		"/admin/update-user",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "edit",
				http.HandlerFunc(adminHandler.UpdateUser),
			),
		),
	)

	// LIST USERS - protected route
	http.Handle(
		"/api/users",
		middleware.AuthMiddleware(
			middleware.RBACMiddleware(database, "users", "view",
				http.HandlerFunc(adminHandler.ListUsers),
			),
		),
	)

	// GET public user info by id (authenticated users)
	http.Handle(
		"/users/get",
		middleware.AuthMiddleware(
			http.HandlerFunc(adminHandler.GetUser),
		),
	)
	db.SeedAdmin(database)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", CORSMiddleware(http.DefaultServeMux)))
}

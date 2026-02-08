# RBAC System - Role-Based Access Control Platform

A comprehensive full-stack Role-Based Access Control (RBAC) system with modern UI/UX, real-time project management, and task tracking. The system provides secure user authentication, granular permission management, and an intuitive interface for managing users, projects, and tasks.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Authentication & Authorization](#authentication--authorization)
- [Roles & Permissions](#roles--permissions)
- [Development Guide](#development-guide)
- [Troubleshooting](#troubleshooting)

## ✨ Features

### Core Features
- **User Authentication**: Secure JWT-based authentication with password hashing
- **Role-Based Access Control**: Four-tier role system (ADMIN, MANAGER, EDITOR, VIEWER)
- **User Management**: Create, update, and manage users with role assignment
- **Project Management**: Create, update, and manage projects with team member assignment
- **Task Management**: Full task lifecycle management with status tracking (TODO → IN_PROGRESS → REVIEW → DONE)
- **Task Board**: Interactive Kanban board with drag-and-drop functionality
- **Dashboard**: Real-time overview of projects, users, and tasks

### Frontend Features
- **Dark Theme**: Professional dark interface with vibrant neon color scheme
- **Responsive Design**: Mobile-friendly, works on desktop, tablet, and mobile
- **User-Friendly Interface**: Intuitive navigation with sidebar and breadcrumbs
- **Real-time Updates**: Instant UI updates on data changes
- **Profile Management**: User profile with editable fields and photo upload

### Backend Features
- **RESTful API**: Comprehensive REST API with proper HTTP status codes
- **Database Migrations**: Version-controlled schema with SQL migrations
- **CORS Support**: Cross-origin resource sharing for frontend integration
- **Configuration-Driven RBAC**: Database-driven permission system
- **Transaction Support**: Database transactions for data consistency

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (React + Vite)                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Dashboard│  │ Projects │  │  Tasks   │  │ Profile  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│         │                │            │            │        │
└─────────┼────────────────┼────────────┼────────────┼────────┘
          │                │            │            │
      ┌───└────────────────┴────────────┴────────────┘─────┐
      │         Axios HTTP Client + JWT Interceptor        │
      └──────────────────────────────────────────────────────┘
                │
                │ HTTP/CORS
                ↓
┌─────────────────────────────────────────────────────────────┐
│                  Backend (Go + SQLite)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ HTTP Router & Middleware                            │   │
│  │ ├── Auth Middleware (JWT verification)             │   │
│  │ ├── RBAC Middleware (permission checking)          │   │
│  │ └── CORS Middleware                                │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ Handlers    │  │ Handlers     │  │ Handlers     │    │
│  │ (Auth, User,│  │ (Projects)   │  │ (Tasks, Stats)   │    │
│  │  Projects)  │  │              │  │                  │    │
│  └──────────────┘  └──────────────┘  └──────────────┘    │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Repositories (Data Access Layer)                      │  │
│  │ ├── User Repository    ├── Project Repository       │  │
│  │ └── Task Repository    └── Permission Repository    │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Database (SQLite)                                    │  │
│  │ ├── users       ├── projects   ├── tasks           │  │
│  │ ├── permissions ├── project_assignments             │  │
│  │ └── role_permissions                               │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ Tech Stack

### Frontend
- **React 18** - Modern UI library with hooks
- **Vite** - Fast build tool and dev server (port 5173)
- **React Router v6** - Client-side routing
- **Axios** - HTTP client with interceptors
- **CSS3** - Custom properties for theming

### Backend
- **Go 1.25.7** - Memory-efficient, fast execution (port 8080)
- **SQLite** - Embedded SQL database
- **JWT (golang-jwt/jwt/v5)** - Token-based authentication
- **bcrypt** - Password hashing (via golang.org/x/crypto)
- **UUID** - Unique identifier generation

## 📁 Project Structure

```
rbac-system/
├── rbac-backend/               # Go backend service
│   ├── cmd/
│   │   ├── server/            # Main server entry point
│   │   ├── migrate/           # Database migration runner
│   │   ├── fixadmin/          # Admin user setup utility
│   │   └── tables/            # Database schema inspection
│   ├── internal/
│   │   ├── auth/              # Authentication utilities
│   │   │   ├── jwt.go         # JWT token handling
│   │   │   └── password.go    # Password hashing & verification
│   │   ├── db/                # Database layer
│   │   │   ├── migrate.go     # Migration runner
│   │   │   ├── sqlite.go      # Database connection
│   │   │   └── permission.go  # Permission queries
│   │   ├── handlers/          # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── project_handler.go
│   │   │   ├── task_handler.go
│   │   │   ├── admin_handler.go
│   │   │   └── protected.go   # Protected routes helper
│   │   ├── middleware/        # HTTP middleware
│   │   │   ├── auth.go        # JWT verification middleware
│   │   │   ├── rbac.go        # Permission checking middleware
│   │   │   └── admin.go       # Admin-only routes middleware
│   │   ├── models/            # Data structures
│   │   │   ├── user.go
│   │   │   ├── project.go
│   │   │   ├── task.go
│   │   │   └── permission.go
│   │   ├── rbac/              # RBAC constants & logic
│   │   │   ├── constants.go   # Role & action definitions
│   │   │   └── permissions.go # Permission logic
│   │   ├── repository/        # Data access layer
│   │   │   ├── user_repository.go
│   │   │   ├── project_repository.go
│   │   │   └── task_repository.go
│   │   └── utils/             # Utility functions
│   │       └── field_filter.go # Field-level access control
│   ├── migrations/            # SQL migration files
│   │   ├── 001_create_user_table.sql
│   │   ├── 002_create_projects_table.sql
│   │   ├── 003_role_permissions.sql
│   │   ├── 004_seed_rbac_permissions.sql
│   │   ├── 005_config_driven_rbac.sql
│   │   ├── 006_create_projects_employees_table.sql
│   │   ├── 007_create_tasks_table.sql
│   │   ├── 008_add_tasks_permissions.sql
│   │   └── 009_add_status_to_projects.sql
│   ├── docs/                  # API documentation
│   ├── go.mod & go.sum        # Go dependencies
│   └── rbac.db                # SQLite database (generated)
│
├── rbac-frontend/             # React frontend application
│   ├── src/
│   │   ├── pages/            # Page components
│   │   │   ├── Login.jsx      # Authentication
│   │   │   ├── Dashboard.jsx  # Main dashboard
│   │   │   ├── Projects.jsx   # Project listing
│   │   │   ├── ProjectDetails.jsx # Project detail & team management
│   │   │   ├── CreateProject.jsx  # Project creation form
│   │   │   ├── Tasks.jsx      # Task management & Kanban board
│   │   │   ├── CreateUser.jsx # User creation & listing
│   │   │   └── Profile.jsx    # User profile page
│   │   ├── components/        # Reusable components
│   │   │   ├── Sidebar.jsx    # Navigation sidebar
│   │   │   ├── MainLayout.jsx # Layout wrapper
│   │   │   ├── Navbar.jsx     # Top navigation
│   │   │   ├── KanbanBoard.jsx # Drag-and-drop task board
│   │   ├── api/               # API client
│   │   │   └── client.js      # Axios client & endpoint helpers
│   │   ├── utils/             # Utility functions
│   │   │   └── jwt.js         # JWT token management
│   │   ├── App.jsx            # Root React component
│   │   ├── main.jsx           # React entry point
│   │   └── index.css          # Global styles
│   ├── public/                # Static assets
│   ├── index.html             # HTML template
│   ├── vite.config.js         # Vite configuration
│   ├── package.json           # Node dependencies
│   └── README.md              # Frontend documentation
│
└── README.md                  # This file
```

## 📋 Prerequisites

### System Requirements
- **For Backend**: Go 1.25.7 or higher
- **For Frontend**: Node.js 16+ and npm/yarn
- **Database**: SQLite (included with backend)
- **Port Availability**: Ports 8080 (backend) and 5173 (frontend)

### Optional
- Git for version control
- VS Code with Go and ESLint extensions
- Postman or Thunder Client for API testing

## 🚀 Installation & Setup

### Backend Setup

1. **Navigate to backend directory**:
```bash
cd rbac-system/rbac-backend
```

2. **Install Go dependencies**:
```bash
go mod tidy
```

3. **Run database migrations**:
```bash
go run cmd/migrate/main.go
```
This will create the SQLite database and apply all migrations.

4. **Seed admin user (optional)**:
```bash
go run cmd/fixadmin/main.go
```
Default admin credentials:
- **Email**: admin@example.com
- **Password**: admin123

### Frontend Setup

1. **Navigate to frontend directory**:
```bash
cd rbac-system/rbac-frontend
```

2. **Install Node dependencies**:
```bash
npm install
```

3. **Verify configuration** (optional):
Check `src/api/client.js` and ensure `API_BASE_URL` points to your backend:
```javascript
const API_BASE_URL = 'http://localhost:8080'
```

## ▶️ Running the Application

### Start Backend Server

```bash
cd rbac-system/rbac-backend
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

Expected output:
```
2026/02/09 14:30:00 Starting server on port :8080
```

### Start Frontend Development Server

In a new terminal:
```bash
cd rbac-system/rbac-frontend
npm run dev
```

The application will start on `http://localhost:5173`

Expected output:
```
  VITE v5.4.21  ready in 282 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### Access the Application

1. Open your browser
2. Navigate to `http://localhost:5173`
3. Log in with default credentials:
   - **Email**: admin@example.com
   - **Password**: admin123

## 📡 API Documentation

### Authentication Endpoints

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**Response (200 OK)**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@example.com",
    "name": "Super Admin",
    "role": "ADMIN"
  }
}
```

### Project Endpoints

#### List Projects
```http
GET /projects
Authorization: Bearer {token}
```

**Response (200 OK)**:
```json
[
  {
    "id": "project-uuid",
    "name": "TechyBrat",
    "description": "Project description",
    "status": "Active",
    "created_by": "user-uuid",
    "assigned_employees": ["user-id-1", "user-id-2"]
  }
]
```

#### Get Project Details
```http
GET /projects/get?id={project_id}
Authorization: Bearer {token}
```

#### Create Project
```http
POST /projects/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "New Project",
  "description": "Project description",
  "status": "Active"
}
```

#### Update Project
```http
PUT /projects/update
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": "project-uuid",
  "name": "Updated Name",
  "description": "Updated description",
  "status": "Completed",
  "assigned_employees": ["user-id-1", "user-id-2"]
}
```

#### Get Project Employees
```http
GET /projects/employees?project_id={project_id}
Authorization: Bearer {token}
```

### Task Endpoints

#### List Tasks
```http
GET /tasks?project_id={project_id}
Authorization: Bearer {token}
```

**Status Values**: `TODO`, `IN_PROGRESS`, `REVIEW`, `DONE`, `ARCHIVED`

#### Create Task
```http
POST /tasks/create
Authorization: Bearer {token}
Content-Type: application/json

{
  "project_id": "project-uuid",
  "title": "Task Title",
  "description": "Task description",
  "status": "TODO"
}
```

#### Update Task
```http
PUT /tasks/update
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": "task-uuid",
  "status": "IN_PROGRESS",
  "title": "Updated Title"
}
```

Status transition flow:
- `TODO` → `IN_PROGRESS` (sets `started_at` timestamp)
- `IN_PROGRESS` → `REVIEW` or `DONE`
- `DONE` (sets `completed_at` timestamp)
- `ARCHIVED` (final state)

### User Endpoints

#### List All Users
```http
GET /api/users
Authorization: Bearer {token}
```

#### Get User by ID
```http
GET /users/get?id={user_id}
Authorization: Bearer {token}
```

#### Create User (Admin Only)
```http
POST /admin/create-user
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure_password",
  "role": "MANAGER"
}
```

**Available Roles**: `ADMIN`, `MANAGER`, `EDITOR`, `VIEWER`

### Statistics Endpoints

#### Get Dashboard Stats
```http
GET /admin/stats
Authorization: Bearer {token}
```

**Response**:
```json
{
  "total_projects": 5,
  "total_users": 10,
  "total_tasks": 25,
  "active_users": 8
}
```

## 🗃️ Database Schema

### Users Table
```sql
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role TEXT CHECK(role IN ('ADMIN','MANAGER','EDITOR','VIEWER')),
  is_active BOOLEAN DEFAULT 1,
  last_login DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Projects Table
```sql
CREATE TABLE projects (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'Active',
  created_by TEXT NOT NULL,
  FOREIGN KEY (created_by) REFERENCES users(id)
);
```

### Tasks Table
```sql
CREATE TABLE tasks (
  id TEXT PRIMARY KEY,
  project_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT CHECK(status IN ('TODO','IN_PROGRESS','REVIEW','DONE','ARCHIVED')),
  assignee TEXT,
  created_by TEXT NOT NULL,
  started_at DATETIME,
  completed_at DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (project_id) REFERENCES projects(id),
  FOREIGN KEY (created_by) REFERENCES users(id)
);
```

### Project Assignments Table
```sql
CREATE TABLE project_assignments (
  project_id TEXT,
  user_id TEXT,
  PRIMARY KEY (project_id, user_id),
  FOREIGN KEY (project_id) REFERENCES projects(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Permissions Table
```sql
CREATE TABLE permissions (
  id TEXT PRIMARY KEY,
  role TEXT NOT NULL,
  resource_table TEXT NOT NULL,
  action TEXT CHECK(action IN ('view','create','edit','delete')),
  field_name TEXT,
  UNIQUE(role, resource_table, action, field_name)
);
```

## 🔐 Authentication & Authorization

### Token-Based Authentication (JWT)

1. **Login Request**: User submits credentials
2. **Token Generation**: Backend validates credentials and generates JWT
3. **Token Storage**: Frontend stores JWT in localStorage
4. **Token Usage**: JWT included in `Authorization: Bearer {token}` header
5. **Token Validation**: Backend verifies token signature and expiration

### JWT Token Structure
```
Header: { "alg": "HS256", "typ": "JWT" }
Payload: {
  "sub": "user_id",
  "email": "user@example.com",
  "role": "ADMIN",
  "exp": 1707533400
}
Signature: HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload))
```

### Password Security
- Passwords are hashed using bcrypt (cost factor 10)
- Passwords never stored in plain text
- Password verification uses constant-time comparison

## 👥 Roles & Permissions

### Role Hierarchy

#### ADMIN
- Full access to all resources
- Can create, read, update, delete users, projects, and tasks
- Can assign permissions to other roles
- No field restrictions

#### MANAGER
- Create and manage projects
- View all projects and team members
- Assign users to projects
- Create and update tasks
- Limited user creation permissions (if configured)

#### EDITOR
- Create and manage tasks
- Update task status
- View assigned projects
- Cannot create projects or manage users

#### VIEWER
- Read-only access to projects and tasks
- Cannot create, edit, or delete resources
- Can view task details and project information

### Configuration-Driven Permissions

Permissions are stored in the database and can be modified without code changes:

```sql
-- Example: Allow MANAGER to create projects
INSERT INTO permissions (id, role, resource_table, action, field_name)
VALUES (uuid(), 'MANAGER', 'projects', 'create', NULL);

-- Example: Restrict EDITOR from viewing certain fields
INSERT INTO permissions (id, role, resource_table, action, field_name)
VALUES (uuid(), 'EDITOR', 'projects', 'view', 'budget');
```

## 🔧 Development Guide

### Backend Development

#### Adding a New API Endpoint

1. **Create Handler** in `internal/handlers/new_handler.go`:
```go
func (h *NewHandler) GetData(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  // Implementation
}
```

2. **Register Route** in `cmd/server/main.go`:
```go
http.Handle(
  "/new-endpoint",
  middleware.AuthMiddleware(
    middleware.RBACMiddleware(database, "table", "action",
      http.HandlerFunc(handler.GetData),
    ),
  ),
)
```

3. **Add Middleware** for authentication and authorization

#### Running Tests
```bash
cd rbac-backend
go test ./...
```

### Frontend Development

#### Adding a New Page

1. **Create Page Component** in `src/pages/NewPage.jsx`:
```jsx
import { useState, useEffect } from 'react'
import { apiClient } from '../api/client'

function NewPage() {
  // Component implementation
  return <div>New Page</div>
}

export default NewPage
```

2. **Add Route** in `src/App.jsx`:
```jsx
<Route path="/new-page" element={<NewPage />} />
```

3. **Add Navigation Link** in sidebar

#### Building for Production
```bash
cd rbac-frontend
npm run build
```

Output will be in `dist/` directory for deployment.

### Development Workflow

1. **Start Backend**:
```bash
cd rbac-system/rbac-backend
go run cmd/server/main.go
```

2. **Start Frontend** (in new terminal):
```bash
cd rbac-system/rbac-frontend
npm run dev
```

3. **Make Changes**: Edit files and see hot-reload
4. **Test Locally**: Navigate to http://localhost:5173
5. **Test API** (optional): Use Postman or cURL

## 📝 Common Development Tasks

### Add a New Database Migration

1. Create SQL file in `migrations/` directory:
```bash
cd rbac-backend/migrations
# Create file: 010_add_new_table.sql
```

2. Add SQL statements:
```sql
CREATE TABLE IF NOT EXISTS new_table (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

3. Run migrations:
```bash
go run cmd/migrate/main.go
```

### Create a New User Role

Currently predefined roles: ADMIN, MANAGER, EDITOR, VIEWER

To add a custom role:
1. Update `internal/rbac/constants.go`
2. Add role to role validation in database
3. Configure permissions for the new role
4. Update frontend role selection dropdowns

## 🐛 Troubleshooting

### Backend Issues

#### Port 8080 Already in Use
```bash
# Find and kill process using port 8080
lsof -ti:8080 | xargs kill -9  # Linux/Mac
netstat -ano | findstr :8080   # Windows
taskkill /PID {PID} /F
```

#### Database Connection Error
- Ensure `rbac.db` exists in `rbac-backend/` directory
- Check file permissions
- Delete corrupted database and re-run migrations:
```bash
rm rbac.db
go run cmd/migrate/main.go
```

#### Migrations Failed
- Check migration file syntax (valid SQL)
- Ensure migrations run in order
- Check database file is writable
- Review logs for specific error

### Frontend Issues

#### Port 5173 Already in Use
```bash
# Kill process using port 5173
lsof -ti:5173 | xargs kill -9  # Linux/Mac
netstat -ano | findstr :5173   # Windows
```

#### API Connection Failed
- Check backend is running on port 8080
- Verify `API_BASE_URL` in `src/api/client.js`
- Check CORS headers in backend response
- Verify internet connectivity

#### Build Fails
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### Login Fails
- Verify admin user exists: `go run cmd/fixadmin/main.go`
- Check password is correct (admin123)
- Verify JWT secret matches between frontend and backend
- Check browser console for error details

### Common Error Messages

| Error | Solution |
|-------|----------|
| "unauthorized" | Check JWT token exists and is valid |
| "forbidden" | User doesn't have permission for this action |
| "project not found" | Verify project ID is correct |
| "invalid request" | Check request body JSON format |
| "CORS error" | Verify frontend port matches CORS config |
| "database is locked" | Close other connections or restart backend |

## 📊 Performance Tips

### Backend Optimization
- Database queries use indexes (check `migrations/`)
- Connection pooling handled by Go's `database/sql`
- JWT validation cached per request
- Field-level filtering reduces data transfer

### Frontend Optimization
- Code splitting with React Router
- Lazy loading for heavy components
- CSS variables for efficient theming
- Axios response interceptors cache tokens

## 🔒 Security Considerations

- ✅ All passwords hashed with bcrypt (cost factor 10)
- ✅ JWT tokens with expiration (configurable)
- ✅ CORS restricted to localhost:5173 (update for production)
- ✅ SQL injection prevention via parameterized queries
- ✅ RBAC enforced on every request
- ✅ Admin-only endpoints protected

### Security Recommendations
1. Change default admin password immediately after setup
2. Use environment variables for JWT secret (currently hardcoded)
3. Enable HTTPS in production
4. Implement rate limiting on authentication endpoints
5. Add audit logging for sensitive operations
6. Regular security updates for dependencies

## 📦 Deployment

### Backend Deployment
```bash
# Build executable
go build -o rbac-server cmd/server/main.go

# Run server
./rbac-server
```

### Frontend Deployment
```bash
# Build static files
npm run build

# Deploy contents of dist/ to web server
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/your-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/your-feature`
5. Submit pull request

## 📄 License

This project is provided as-is for educational and commercial use.

## 👥 Support

For issues or questions:
1. Check troubleshooting section
2. Review API documentation
3. Check frontend README.md
4. Contact development team

## 🎯 Roadmap

- [ ] Real-time notifications with WebSockets
- [ ] File upload and attachment support
- [ ] Activity audit logs
- [ ] User profile customization
- [ ] Email notifications
- [ ] Two-factor authentication (2FA)
- [ ] Advanced filtering and search
- [ ] Team invitations
- [ ] Project templates
- [ ] Mobile native app

## 📚 Additional Resources

- [Frontend README](./rbac-frontend/README.md)
- [Frontend API Documentation](./rbac-frontend/src/api/client.js)
- [Task API Documentation](./rbac-backend/docs/tasks.md)
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

---

**Version**: 1.0.0  
**Last Updated**: February 9, 2026  
**Status**: Production Ready ✅

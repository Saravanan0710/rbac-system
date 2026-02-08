# RBAC Frontend

A modern, dark-themed frontend for the Role-Based Access Control (RBAC) system with vibrant neon colors and responsive design.

## Features

- **Dark Theme**: Professional dark interface with vibrant neon accents
- **Vibrant UI**: Neon green, cyan, pink, and orange button colors for visual appeal
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile devices
- **Modern Stack**: React + Vite + React Router
- **Authentication**: JWT-based authentication system
- **Kanban Board**: Interactive drag-and-drop task management
- **User Management**: Create and manage users with roles
- **Project Management**: Create and manage projects
- **Dashboard**: Real-time statistics and project overview
- **Accessibility**: Semantic HTML and keyboard navigation support

## Tech Stack

- **React 18** - UI library
- **Vite** - Build tool and dev server
- **React Router DOM** - Client-side routing
- **Axios** - HTTP client
- **CSS3** - Modern styling with CSS variables

## Project Structure

```
src/
├── pages/           # Page components
│   ├── Login.jsx
│   ├── Dashboard.jsx
│   ├── CreateProject.jsx
│   ├── ProjectDetails.jsx
│   ├── Tasks.jsx
│   ├── CreateUser.jsx
│   └── Profile.jsx
├── components/      # Reusable components
│   ├── Navbar.jsx
│   └── KanbanBoard.jsx
├── api/             # API integration
│   └── client.js
├── utils/           # Utility functions
│   └── jwt.js
├── App.jsx
├── main.jsx
└── index.css        # Global styles
```

## Installation

### Prerequisites

- Node.js 16+ and npm

### Setup

1. Install dependencies:
```bash
npm install
```

2. Update API endpoint in `src/api/client.js` if needed:
```javascript
const API_BASE_URL = 'http://localhost:8080/api'
```

## Development

Run the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## Build

Create a production build:
```bash
npm run build
```

Preview production build:
```bash
npm run preview
```

## Color Theme

### Primary Colors (Vibrant Neon)
- **Green**: #00ff88 (Primary action button)
- **Cyan**: #00d4ff (Secondary accent)
- **Pink**: #ff006e (Danger/secondary actions)
- **Orange**: #ff8c42 (Warning)

### Background Colors (Dark Neon Theme)
- **Primary**: #0f0f0f
- **Secondary**: #1a1a1a
- **Tertiary**: #2a2a2a
- **Border**: #333333

### Text Colors
- **Primary**: #ffffff
- **Secondary**: #b0b0b0

## Pages

### Login
- Email and password authentication
- JWT token storage
- Redirect to dashboard on success

### Dashboard
- Overview of statistics (projects, users, tasks, roles)
- Project listing with quick access
- Create new project button

### Create Project
- Form to create new projects
- Project name, description, and status
- Form validation

### Project Details
- View project information
- List of assigned tasks
- Team members list
- Task deletion

### Tasks (Kanban Board)
- Drag-and-drop task management
- Four columns: To Do, In Progress, Review, Done
- Create new tasks
- Task priority and status indicators
- Task editing and deletion

### User Management
- List all users
- Create new users
- Delete users
- Display user roles and status
- User information table

### Profile
- View user information
- Edit profile details
- Security settings
- Role information
- Password management options

## Styling

All components use CSS custom properties for easy theming:

```css
:root {
  --bg-primary: #0f0f0f;
  --btn-primary: #00ff88;
  --accent-cyan: #00d4ff;
  /* ... and more */
}
```

## API Integration

The frontend connects to the backend API with the following endpoints:

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Users
- `GET /api/users` - List all users
- `POST /api/users` - Create user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Projects
- `GET /api/projects` - List all projects
- `POST /api/projects` - Create project
- `GET /api/projects/:id` - Get project details
- `GET /api/projects/:id/employees` - Get project employees

### Tasks
- `GET /api/tasks` - List all tasks
- `POST /api/tasks` - Create task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task

### Admin
- `GET /api/admin/stats` - Dashboard statistics
- `GET /api/admin/users` - List all users (admin)
- `POST /api/admin/users` - Create user (admin)

## Changelog (Recent)

- Enforced project-level RBAC on backend: non-ADMIN roles (MANAGER, EDITOR, VIEWER) can only view or modify projects and tasks when assigned to the project.
- Task-level assignees: tasks now support an `assignees` list and a single `assignee` field. Assignees can update task `status`.
- Role restrictions: VIEWER/EDITOR cannot create tasks; MANAGER may create tasks only for projects they're assigned to; ADMIN bypasses assignment checks.
- Frontend: added assignee dropdown in task creation, role-aware show/hide for create buttons, and JWT role helper in `src/utils/jwt.js`.

## Backend & Local Run (Quick)

1. Prepare DB and migrations (backend repo):
```powershell
cd rbac-backend
go run .\cmd\migrate\main.go
```

2. Build and run the server:
```powershell
cd rbac-backend
go build -o server .\cmd\server\main.go
.\server.exe
```

The API server listens on `http://localhost:8080` by default.

## RBAC Notes (developer)

- JWT includes role information used by the frontend to show/hide controls. See `src/utils/jwt.js` for helpers.
- Important endpoints touched by RBAC logic:
  - `POST /login` — returns JWT
  - `GET /projects`, `GET /projects/:id` — project visibility requires assignment for non-ADMIN
  - `POST /tasks`, `PUT /tasks` — creation and update rules depend on role and project assignment
  - `PUT /admin/update-user` — admin can update user role


## Responsive Breakpoints

- **Mobile**: < 480px
- **Tablet**: 480px - 768px
- **Desktop**: > 768px

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

MIT License - see LICENSE file for details

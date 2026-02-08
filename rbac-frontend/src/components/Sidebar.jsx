import { Link, useNavigate, useLocation } from 'react-router-dom'
import { getAuthToken, removeAuthToken, getUserFromToken } from '../utils/jwt'
import './Sidebar.css'

function Sidebar({ onLogout }) {
  const navigate = useNavigate()
  const location = useLocation()
  const token = getAuthToken()
  const user = token ? getUserFromToken(token) : null

  const handleLogout = () => {
    removeAuthToken()
    onLogout()
    navigate('/login')
  }

  const isActive = (path) => location.pathname === path

  // Generate avatar from first letter of name or email
  const avatarLetter = (user?.name || user?.email || 'U').charAt(0).toUpperCase()

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        <Link to="/dashboard" className="logo-text">
          RBAC
        </Link>
      </div>

      <nav className="sidebar-nav">
        <Link 
          to="/dashboard" 
          className={`nav-item ${isActive('/dashboard') ? 'active' : ''}`}
        >
          <span className="nav-icon">📊</span>
          <span className="nav-label">Dashboard</span>
        </Link>
        <Link 
          to="/tasks" 
          className={`nav-item ${isActive('/tasks') ? 'active' : ''}`}
        >
          <span className="nav-icon">✓</span>
          <span className="nav-label">Tasks</span>
        </Link>
        <Link 
          to="/projects" 
          className={`nav-item ${isActive('/projects') ? 'active' : ''}`}
        >
          <span className="nav-icon">📁</span>
          <span className="nav-label">Projects</span>
        </Link>
        <Link 
          to="/users" 
          className={`nav-item ${isActive('/users') ? 'active' : ''}`}
        >
          <span className="nav-icon">👥</span>
          <span className="nav-label">Users</span>
        </Link>
        <Link 
          to="/profile" 
          className={`nav-item ${isActive('/profile') ? 'active' : ''}`}
        >
          <span className="nav-icon">⚙️</span>
          <span className="nav-label">Profile</span>
        </Link>
      </nav>

      <div className="sidebar-footer">
        <div className="user-info">
          <div className="user-avatar">{avatarLetter}</div>
          <div className="user-details">
            <p className="user-name">{user?.name || 'User'}</p>
            <p className="user-email">{user?.email}</p>
          </div>
        </div>
        <button 
          className="logout-btn" 
          onClick={handleLogout}
          title="Logout"
          aria-label="Logout"
        >
          🚪
        </button>
      </div>
    </aside>
  )
}

export default Sidebar

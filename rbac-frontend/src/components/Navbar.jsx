import { Link, useNavigate } from 'react-router-dom'
import { getAuthToken, removeAuthToken, getUserFromToken } from '../utils/jwt'
import './Navbar.css'

function Navbar({ onLogout }) {
  const navigate = useNavigate()
  const token = getAuthToken()
  const user = token ? getUserFromToken(token) : null

  const handleLogout = () => {
    removeAuthToken()
    onLogout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <div className="nav-container">
        <Link to="/dashboard" className="nav-logo">
          RBAC System
        </Link>
        <ul className="nav-menu">
          <li>
            <Link to="/dashboard">Dashboard</Link>
          </li>
          <li>
            <Link to="/tasks">Tasks</Link>
          </li>
          <li>
            <Link to="/projects">Projects</Link>
          </li>
          <li>
            <Link to="/users">Users</Link>
          </li>
          <li>
            <Link to="/profile">Profile</Link>
          </li>
        </ul>
        <div className="nav-user">
          <span className="user-email">{user?.email}</span>
          <button className="btn btn-danger btn-small" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>
    </nav>
  )
}

export default Navbar

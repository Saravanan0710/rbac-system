import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { projectsAPI, adminAPI } from '../api/client'
import './Dashboard.css'

function Dashboard() {
  const [projects, setProjects] = useState([])
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [projectsRes, statsRes] = await Promise.all([
          projectsAPI.getAll(),
          adminAPI.getStats()
        ])
        const projectList = Array.isArray(projectsRes?.data) ? projectsRes.data : []
        const statsObj = statsRes?.data && typeof statsRes.data === 'object' ? statsRes.data : {}
        setProjects(projectList)
        setStats(statsObj)
      } catch (err) {
        setError('Failed to load dashboard')
        console.error(err)
        setProjects([])
        setStats({})
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="container">
        <div className="flex-center" style={{ minHeight: '400px' }}>
          <div className="spinner"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="dashboard">
      <div className="container">
        <div className="dashboard-header">
          <h1>Dashboard</h1>
        </div>

        {error && <div className="alert alert-error">{error}</div>}

        {stats && (
          <div className="stats-grid">
            <div className="stat-card">
              <div className="stat-value">{stats.total_projects || 0}</div>
              <div className="stat-label">Total Projects</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{stats.total_users || 0}</div>
              <div className="stat-label">Total Users</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{stats.total_tasks || 0}</div>
              <div className="stat-label">Total Tasks</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{stats.total_roles || 0}</div>
              <div className="stat-label">Total Roles</div>
            </div>
          </div>
        )}

        <h2>Projects</h2>
        {projects.length === 0 ? (
          <div className="alert alert-info">No projects found. Create one to get started!</div>
        ) : (
          <div className="projects-grid">
            {projects.map((project) => (
              <Link to={`/projects/${project.id}`} key={project.id} className="project-card">
                <h3>{project.name}</h3>
                <p>{project.description}</p>
                <div className="project-footer">
                  <span className="project-status">{project.status || 'Active'}</span>
                  <span className="project-arrow">→</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Dashboard

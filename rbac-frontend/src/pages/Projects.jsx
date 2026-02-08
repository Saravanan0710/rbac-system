import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { projectsAPI } from '../api/client'
import { getAuthToken, getUserRoleFromToken } from '../utils/jwt'
import './Dashboard.css'

function Projects() {
  const [projects, setProjects] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const res = await projectsAPI.getAll()
        const list = Array.isArray(res?.data) ? res.data : []
        setProjects(list)
      } catch (err) {
        setError('Failed to load projects')
        console.error(err)
        setProjects([])
      } finally {
        setLoading(false)
      }
    }
    fetchProjects()
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
          <h1>Projects</h1>
          {(['ADMIN','MANAGER'].includes(getUserRoleFromToken(getAuthToken())) ) && (
            <Link to="/create-project" className="btn btn-primary">
              Create New Project
            </Link>
          )}
        </div>

        {error && <div className="alert alert-error">{error}</div>}

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

export default Projects

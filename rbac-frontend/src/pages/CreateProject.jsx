import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { projectsAPI } from '../api/client'
import './CreateProject.css'

function CreateProject() {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    status: 'Active'
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await projectsAPI.create(formData)
      navigate('/dashboard')
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to create project')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="create-project-container">
      <div className="container">
        <div className="form-wrapper">
          <h1>Create New Project</h1>
          <p className="form-subtitle">Set up a new project for your team</p>

          {error && <div className="alert alert-error">{error}</div>}

          <form onSubmit={handleSubmit} className="create-form">
            <div className="form-group">
              <label htmlFor="name">Project Name *</label>
              <input
                id="name"
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="Enter project name"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="description">Description</label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                placeholder="Enter project description"
                rows="4"
              ></textarea>
            </div>

            <div className="form-group">
              <label htmlFor="status">Status</label>
              <select
                id="status"
                name="status"
                value={formData.status}
                onChange={handleChange}
              >
                <option value="Active">Active</option>
                <option value="Inactive">Inactive</option>
                <option value="Planning">Planning</option>
                <option value="Completed">Completed</option>
              </select>
            </div>

            <div className="form-actions">
              <button type="submit" className="btn btn-primary btn-large" disabled={loading}>
                {loading ? <span className="spinner"></span> : 'Create Project'}
              </button>
              <button
                type="button"
                className="btn btn-outline btn-large"
                onClick={() => navigate('/dashboard')}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

export default CreateProject

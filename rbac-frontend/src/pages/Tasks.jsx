import { useState, useEffect } from 'react'
import { tasksAPI, projectsAPI } from '../api/client'
import { getAuthToken, getUserRoleFromToken } from '../utils/jwt'
import KanbanBoard from '../components/KanbanBoard'
import './Tasks.css'

function Tasks() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [projects, setProjects] = useState([])
  const [selectedProject, setSelectedProject] = useState('')
  const [projectEmployees, setProjectEmployees] = useState([])
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    priority: 'Medium',
    status: 'TODO'
  })

  useEffect(() => {
    const init = async () => {
      await fetchProjects()
    }
    init()
  }, [])

  useEffect(() => {
    fetchTasks()
    // fetch employees for selected project
    if (selectedProject) fetchProjectEmployees(selectedProject)
  }, [selectedProject])

  const fetchTasks = async () => {
    try {
      let response
      if (selectedProject) {
        response = await tasksAPI.getByProject(selectedProject)
      } else {
        response = await tasksAPI.getAll()
      }
      const list = Array.isArray(response?.data) ? response.data : []
      setTasks(list)
    } catch (err) {
      setError('Failed to load tasks')
      console.error(err)
      setTasks([])
    } finally {
      setLoading(false)
    }
  }

  const fetchProjects = async () => {
    try {
      const res = await projectsAPI.getAll()
      const list = Array.isArray(res?.data) ? res.data : []
      setProjects(list)
      if (list.length > 0 && !selectedProject) setSelectedProject(list[0].id)
    } catch (err) {
      console.error('Failed to load projects', err)
      setProjects([])
    }
  }

  const fetchProjectEmployees = async (projectId) => {
    try {
      const res = await projectsAPI.getEmployees(projectId)
      const list = Array.isArray(res?.data) ? res.data : []
      setProjectEmployees(list)
    } catch (err) {
      console.error('Failed to load project employees', err)
      setProjectEmployees([])
    }
  }

  const handleCreateTask = async (e) => {
    e.preventDefault()
    try {
      // get project_id from form or selected filter
      const projectId = formData.project_id || selectedProject
      if (!projectId) {
        setError('Please select a project for the task')
        return
      }
      
      // ensure assignee selected
      const assignee = formData.assignee || ''
      if (!assignee) {
        setError('Please select an assignee for the task')
        return
      }

      const payload = { ...formData, project_id: projectId, assignee }
      const response = await tasksAPI.create(payload)
      
      if (response.status === 200 || response.status === 201) {
        // Reset form
        setFormData({
          title: '',
          description: '',
          priority: 'Medium',
          status: 'TODO',
          assignee: ''
        })
        setError('')
        setShowForm(false)
        
        // Refresh tasks
        await fetchTasks()
      }
    } catch (err) {
      const errorMsg = err.response?.data?.message || err.message || 'Failed to create task'
      setError(errorMsg)
      console.error('Task creation error:', err)
    }
  }

  const handleUpdateTask = async (updatedTask) => {
    try {
      await tasksAPI.update(updatedTask.id, updatedTask)
      setTasks((prev) =>
        prev.map((t) => (t.id === updatedTask.id ? updatedTask : t))
      )
    } catch (err) {
      setError('Failed to update task')
      console.error(err)
    }
  }

  const handleDeleteTask = async (taskId) => {
    if (!window.confirm('Are you sure you want to delete this task?')) return
    try {
      await tasksAPI.delete(taskId)
      setTasks((prev) => prev.filter((t) => t.id !== taskId))
    } catch (err) {
      setError('Failed to delete task')
      console.error(err)
    }
  }

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
    <div className="tasks-page">
      <div className="tasks-header">
        <div className="container">
          <div className="header-top">
            <h1>Tasks & Kanban Board</h1>
            <div className="header-actions">
              <div className="filter-group">
                <label htmlFor="filter-project">Filter by Project:</label>
                <select
                  id="filter-project"
                  value={selectedProject}
                  onChange={(e) => setSelectedProject(e.target.value)}
                >
                  <option value="">-- All Projects --</option>
                  {projects.map((p) => (
                    <option key={p.id} value={p.id}>{p.name}</option>
                  ))}
                </select>
              </div>
              {(['ADMIN','MANAGER'].includes(getUserRoleFromToken(getAuthToken())) ) && (
                <button
                  className="btn btn-primary"
                  onClick={() => setShowForm(!showForm)}
                >
                  {showForm ? 'Cancel' : 'Create Task'}
                </button>
              )}
            </div>
          </div>

          {error && <div className="alert alert-error">{error}</div>}

          {showForm && (
            <div className="create-task-form">
              <form onSubmit={handleCreateTask}>
                <div className="form-group">
                  <label htmlFor="project">Project *</label>
                  <select
                    id="project"
                    value={formData.project_id || selectedProject}
                    onChange={(e) => {
                      const pid = e.target.value
                      setFormData((prev) => ({ ...prev, project_id: pid }))
                      setSelectedProject(pid)
                    }}
                    required
                  >
                    <option value="">-- All Projects --</option>
                    {projects.map((p) => (
                      <option key={p.id} value={p.id}>{p.name}</option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label htmlFor="assignee">Assignee *</label>
                  <select
                    id="assignee"
                    value={formData.assignee || ''}
                    onChange={(e) => setFormData((prev) => ({ ...prev, assignee: e.target.value }))}
                    required
                  >
                    <option value="">-- Select Assignee --</option>
                    {projectEmployees.map((u) => (
                      <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                    ))}
                  </select>
                </div>
                <div className="form-row">
                  <div className="form-group">
                    <label htmlFor="title">Task Title *</label>
                    <input
                      id="title"
                      type="text"
                      value={formData.title}
                      onChange={(e) =>
                        setFormData((prev) => ({ ...prev, title: e.target.value }))
                      }
                      placeholder="Enter task title"
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label htmlFor="priority">Priority</label>
                    <select
                      id="priority"
                      value={formData.priority}
                      onChange={(e) =>
                        setFormData((prev) => ({ ...prev, priority: e.target.value }))
                      }
                    >
                      <option value="Low">Low</option>
                      <option value="Medium">Medium</option>
                      <option value="High">High</option>
                    </select>
                  </div>
                </div>
                <div className="form-group">
                  <label htmlFor="description">Description</label>
                  <textarea
                    id="description"
                    value={formData.description}
                    onChange={(e) =>
                      setFormData((prev) => ({ ...prev, description: e.target.value }))
                    }
                    placeholder="Enter task description"
                    rows="3"
                  ></textarea>
                </div>
                <button type="submit" className="btn btn-success">
                  Create Task
                </button>
              </form>
            </div>
          )}
        </div>
      </div>

      <KanbanBoard
        tasks={tasks}
        onTaskUpdate={handleUpdateTask}
        onTaskDelete={handleDeleteTask}
      />
    </div>
  )
}

export default Tasks

import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { projectsAPI, tasksAPI, usersAPI } from '../api/client'
import './ProjectDetails.css'

function ProjectDetails() {
  const { id } = useParams()
  const [project, setProject] = useState(null)
  const [tasks, setTasks] = useState([])
  const [employees, setEmployees] = useState([])
  const [allUsers, setAllUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [editingStatus, setEditingStatus] = useState(false)
  const [newStatus, setNewStatus] = useState('')
  const [updating, setUpdating] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    if (!id) {
      setLoading(false)
      return
    }
    const fetchProjectData = async () => {
      try {
        const [projectRes, tasksRes, employeesRes, usersRes] = await Promise.all([
          projectsAPI.getById(id),
          tasksAPI.getByProject(id),
          projectsAPI.getEmployees(id),
          usersAPI.getAll()
        ])
        const proj = projectRes?.data && typeof projectRes.data === 'object' ? projectRes.data : null
        const taskList = Array.isArray(tasksRes?.data) ? tasksRes.data : []
        const empList = Array.isArray(employeesRes?.data) ? employeesRes.data : []
        const userList = Array.isArray(usersRes?.data) ? usersRes.data : []
        setProject(proj)
        setTasks(taskList)
        setEmployees(empList)
        setAllUsers(userList)
        if (proj?.status) {
          setNewStatus(proj.status)
        }
      } catch (err) {
        setError('Failed to load project details')
        console.error(err)
        setProject(null)
        setTasks([])
        setEmployees([])
        setAllUsers([])
      } finally {
        setLoading(false)
      }
    }

    fetchProjectData()
  }, [id])

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

  const handleUpdateStatus = async () => {
    if (!newStatus.trim() || newStatus === project?.status) {
      setEditingStatus(false)
      return
    }

    setUpdating(true)
    try {
      await projectsAPI.update(id, { status: newStatus })
      setProject((prev) => ({ ...prev, status: newStatus }))
      setEditingStatus(false)
    } catch (err) {
      setError('Failed to update project status')
      console.error(err)
    } finally {
      setUpdating(false)
    }
  }

  const handleAddEmployee = async (userId) => {
    const employeeIds = employees.map((e) => e.id)
    if (employeeIds.includes(userId)) {
      setError('User already assigned to project')
      return
    }

    const newEmployeeIds = [...employeeIds, userId]
    setUpdating(true)
    try {
      await projectsAPI.update(id, { assigned_employees: newEmployeeIds })
      const updatedEmployees = allUsers.filter((u) => newEmployeeIds.includes(u.id))
      setEmployees(updatedEmployees)
    } catch (err) {
      setError('Failed to add employee to project')
      console.error(err)
    } finally {
      setUpdating(false)
    }
  }

  const handleRemoveEmployee = async (userId) => {
    if (!window.confirm('Remove this employee from the project?')) return

    const employeeIds = employees.map((e) => e.id).filter((id) => id !== userId)

    setUpdating(true)
    try {
      await projectsAPI.update(id, { assigned_employees: employeeIds })
      setEmployees((prev) => prev.filter((e) => e.id !== userId))
    } catch (err) {
      setError('Failed to remove employee from project')
      console.error(err)
    } finally {
      setUpdating(false)
    }
  }

  if (loading) {
    return (
      <div className="container" style={{ minHeight: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner"></div>
      </div>
    )
  }

  if (!project) {
    return (
      <div className="container">
        <div className="alert alert-error">Project not found</div>
        <button className="btn btn-primary" onClick={() => navigate('/dashboard')}>
          Back to Dashboard
        </button>
      </div>
    )
  }

  const availableUsers = allUsers.filter((u) => !employees.some((e) => e.id === u.id))

  return (
    <div className="project-details">
      <div className="container">
        <button className="btn btn-outline btn-small" onClick={() => navigate('/dashboard')}>
          ← Back to Dashboard
        </button>

        <div className="project-header">
          <div>
            <h1>{project?.name ?? 'Project'}</h1>
            <p>{project?.description ?? ''}</p>
          </div>
          <div className="project-info">
            {editingStatus ? (
              <div className="status-edit">
                <select value={newStatus} onChange={(e) => setNewStatus(e.target.value)} disabled={updating}>
                  <option value="Active">Active</option>
                  <option value="Inactive">Inactive</option>
                  <option value="Planning">Planning</option>
                  <option value="Completed">Completed</option>
                </select>
                <button
                  className="btn btn-primary btn-small"
                  onClick={handleUpdateStatus}
                  disabled={updating || !newStatus.trim()}
                >
                  Save
                </button>
                <button
                  className="btn btn-outline btn-small"
                  onClick={() => {
                    setEditingStatus(false)
                    setNewStatus(project?.status || '')
                  }}
                  disabled={updating}
                >
                  Cancel
                </button>
              </div>
            ) : (
              <div onClick={() => setEditingStatus(true)} style={{ cursor: 'pointer' }}>
                <span className="status-badge">{project?.status ?? 'Active'}</span>
              </div>
            )}
          </div>
        </div>

        {error && <div className="alert alert-error">{error}</div>}

        <div className="grid grid-2">
          <div className="section">
            <h2>Tasks ({tasks.length})</h2>
            {tasks.length === 0 ? (
              <p className="empty-message">No tasks yet</p>
            ) : (
              <div className="tasks-list">
                {tasks.map((task) => (
                  <div key={task.id} className="task-item">
                    <div className="task-item-header">
                      <h4>{task.title ?? 'Untitled'}</h4>
                      <span className={`priority-${(task.priority || task.status || 'medium').toLowerCase()}`}>
                        {task.status ?? task.priority ?? '—'}
                      </span>
                    </div>
                    <p>{task.description ?? ''}</p>
                    <div className="task-item-footer">
                      <span>{task.status ?? '—'}</span>
                      <button className="btn btn-danger btn-small" onClick={() => handleDeleteTask(task.id)}>
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="section">
            <h2>Team Members ({employees.length})</h2>
            {employees.length === 0 ? (
              <p className="empty-message">No team members assigned</p>
            ) : (
              <div className="employees-list">
                {employees.map((employee) => (
                  <div key={employee.id} className="employee-card">
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                      <div>
                        <h4>{employee.name}</h4>
                        <p>{employee.email}</p>
                        <span className="role-badge">{employee.role}</span>
                      </div>
                      <button
                        className="btn btn-danger btn-small"
                        onClick={() => handleRemoveEmployee(employee.id)}
                        disabled={updating}
                        title="Remove from project"
                      >
                        ✕
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {availableUsers.length > 0 && (
              <div style={{ marginTop: '20px' }}>
                <h3 style={{ fontSize: '14px', marginBottom: '10px' }}>Add Team Member</h3>
                <div className="add-employee">
                  <select defaultValue="" onChange={(e) => {
                    if (e.target.value) {
                      handleAddEmployee(e.target.value)
                      e.target.value = ''
                    }
                  }} disabled={updating}>
                    <option value="">Select user to add...</option>
                    {availableUsers.map((user) => (
                      <option key={user.id} value={user.id}>
                        {user.name} ({user.email})
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default ProjectDetails

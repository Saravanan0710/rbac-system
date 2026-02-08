import { useState, useEffect } from 'react'
import { adminAPI, usersAPI } from '../api/client'
import { getAuthToken, getUserRoleFromToken } from '../utils/jwt'
import './CreateUser.css'

function CreateUser() {
  const [users, setUsers] = useState([])
  const [showForm, setShowForm] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [editingUserId, setEditingUserId] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    role: 'VIEWER'
  })
  const [editFormData, setEditFormData] = useState({
    name: '',
    role: 'VIEWER'
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await adminAPI.getUsers()
      // API returns array directly, not wrapped in object
      const list = Array.isArray(response?.data) ? response.data : []
      setUsers(list)
    } catch (err) {
      setError('Failed to load users')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateUser = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      const role = ['ADMIN', 'MANAGER', 'EDITOR', 'VIEWER'].includes(formData.role) ? formData.role : 'VIEWER'
      await adminAPI.createUser({
        name: formData.name,
        email: formData.email,
        password: formData.password,
        role
      })
      setSuccess('User created successfully')
      setFormData({
        name: '',
        email: '',
        password: '',
        role: 'VIEWER'
      })
      setShowForm(false)
      fetchUsers()
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to create user')
      console.error(err)
    }
  }

  const handleEditUser = (user) => {
    setEditingUserId(user.id)
    setEditFormData({
      name: user.name,
      role: user.role
    })
    setShowEditModal(true)
    setError('')
    setSuccess('')
  }

  const handleUpdateUser = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      const role = ['ADMIN', 'MANAGER', 'EDITOR', 'VIEWER'].includes(editFormData.role) ? editFormData.role : 'VIEWER'
      await usersAPI.update(editingUserId, {
        name: editFormData.name,
        role
      })
      setSuccess('User updated successfully')
      setShowEditModal(false)
      setEditingUserId(null)
      fetchUsers()
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to update user')
      console.error(err)
    }
  }

  const handleDeleteUser = async (userId) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return

    try {
      await usersAPI.delete(userId)
      setSuccess('User deleted successfully')
      fetchUsers()
    } catch (err) {
      setError('Failed to delete user')
      console.error(err)
    }
  }

  if (loading) {
    return (
      <div className="container" style={{ minHeight: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner"></div>
      </div>
    )
  }

  return (
    <div className="users-page">
      <div className="container">
        <div className="users-header">
          <h1>User Management</h1>
          {getUserRoleFromToken(getAuthToken()) === 'ADMIN' && (
            <button
              className="btn btn-primary"
              onClick={() => setShowForm(!showForm)}
            >
              {showForm ? 'Cancel' : 'Create User'}
            </button>
          )}
        </div>

        {error && <div className="alert alert-error">{error}</div>}
        {success && <div className="alert alert-success">{success}</div>}

        {showForm && (
          <div className="create-user-form">
            <h2>Create New User</h2>
            <form onSubmit={handleCreateUser}>
              <div className="form-group">
                <label htmlFor="name">Full Name *</label>
                <input
                  id="name"
                  type="text"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, name: e.target.value }))
                  }
                  placeholder="Enter full name"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="email">Email Address *</label>
                <input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, email: e.target.value }))
                  }
                  placeholder="Enter email address"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="password">Password *</label>
                <input
                  id="password"
                  type="password"
                  value={formData.password}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, password: e.target.value }))
                  }
                  placeholder="Enter password"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="role">Role *</label>
                <select
                  id="role"
                  value={formData.role}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, role: e.target.value }))
                  }
                >
                  <option value="ADMIN">Admin</option>
                  <option value="MANAGER">Manager</option>
                  <option value="EDITOR">Editor</option>
                  <option value="VIEWER">Viewer</option>
                </select>
              </div>

              <button type="submit" className="btn btn-success btn-large">
                Create User
              </button>
            </form>
          </div>
        )}

        {showEditModal && (
          <div className="modal-overlay" onClick={() => setShowEditModal(false)}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h2>Edit User</h2>
                <button className="modal-close" onClick={() => setShowEditModal(false)}>×</button>
              </div>
              <form onSubmit={handleUpdateUser}>
                <div className="form-group">
                  <label htmlFor="edit-name">Full Name *</label>
                  <input
                    id="edit-name"
                    type="text"
                    value={editFormData.name}
                    onChange={(e) =>
                      setEditFormData((prev) => ({ ...prev, name: e.target.value }))
                    }
                    placeholder="Enter full name"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="edit-role">Role *</label>
                  <select
                    id="edit-role"
                    value={editFormData.role}
                    onChange={(e) =>
                      setEditFormData((prev) => ({ ...prev, role: e.target.value }))
                    }
                  >
                    <option value="ADMIN">Admin</option>
                    <option value="MANAGER">Manager</option>
                    <option value="EDITOR">Editor</option>
                    <option value="VIEWER">Viewer</option>
                  </select>
                </div>

                <div className="modal-actions">
                  <button type="submit" className="btn btn-success">
                    Update User
                  </button>
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={() => setShowEditModal(false)}
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        <h2>Users List</h2>
        {users.length === 0 ? (
          <div className="alert alert-info">No users found</div>
        ) : (
          <div className="users-table-wrapper">
            <table className="users-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id}>
                    <td>
                      <strong>{user.name}</strong>
                    </td>
                    <td>{user.email}</td>
                    <td>
                      <span className="role-badge">{user.role || 'User'}</span>
                    </td>
                    <td>
                      <span className="status-badge active">{user.status || 'Active'}</span>
                    </td>
                    <td>
                      <div className="action-buttons">
                        <button className="btn btn-tertiary btn-small" onClick={() => handleEditUser(user)}>Edit</button>
                        <button
                          className="btn btn-danger btn-small"
                          onClick={() => handleDeleteUser(user.id)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default CreateUser

import axios from 'axios'
import { getAuthToken } from '../utils/jwt'

const API_BASE_URL = 'http://localhost:8080'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

apiClient.interceptors.request.use((config) => {
  const token = getAuthToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth API
export const authAPI = {
  login: (email, password) => 
    apiClient.post('/login', { email, password }),
  
  register: (email, password) => 
    apiClient.post('/register', { email, password }),
  
  logout: () => 
    apiClient.post('/logout')
}

// Users API (backend: /api/users for list, /users/get?id= for one, /admin/create-user, /admin/delete-user)
export const usersAPI = {
  getAll: () => 
    apiClient.get('/api/users'),
  
  getById: (id) => 
    apiClient.get('/users/get', { params: { id } }),
  
  getCurrentUser: (userId) => 
    apiClient.get('/users/get', { params: { id: userId } }),
  
  create: (userData) => 
    apiClient.post('/admin/create-user', userData),
  
  update: (id, userData) => 
    apiClient.put('/admin/update-user', { ...userData, id }),
  
  delete: (id) => 
    apiClient.delete('/admin/delete-user', { params: { id } })
}

// Projects API (backend uses /projects, /projects/get?id=, /projects/employees?project_id=, /projects/create, /projects/update, /projects/delete)
export const projectsAPI = {
  getAll: () => 
    apiClient.get('/projects'),
  
  getById: (id) => 
    apiClient.get('/projects/get', { params: { id } }),
  
  create: (projectData) => 
    apiClient.post('/projects/create', projectData),
  
  update: (id, projectData) => 
    apiClient.put('/projects/update', { ...projectData, id }),
  
  delete: (id) => 
    apiClient.delete('/projects/delete', { params: { id } }),
  
  getEmployees: (id) => 
    apiClient.get('/projects/employees', { params: { project_id: id } })
}

// Tasks API (backend: GET /tasks?project_id= or ?assignee=, /tasks/get?id=, /tasks/create, /tasks/update, /tasks/delete?id=)
export const tasksAPI = {
  getAll: (params = {}) => 
    apiClient.get('/tasks', { params }),
  
  getById: (id) => 
    apiClient.get('/tasks/get', { params: { id } }),
  
  getByProject: (projectId) => 
    apiClient.get('/tasks', { params: { project_id: projectId } }),
  
  create: (taskData) => 
    apiClient.post('/tasks/create', taskData),
  
  update: (id, taskData) => 
    apiClient.put('/tasks/update', { ...taskData, id }),
  
  delete: (id) => 
    apiClient.delete('/tasks/delete', { params: { id } })
}

// Admin API (backend: /admin/stats, /api/users for list, /admin/create-user, /admin/delete-user)
export const adminAPI = {
  getStats: () => 
    apiClient.get('/admin/stats'),
  
  getUsers: () => 
    apiClient.get('/api/users'),
  
  createUser: (userData) => 
    apiClient.post('/admin/create-user', userData),
  
  assignRole: (userId, roleId) => 
    apiClient.post(`/admin/users/${userId}/role`, { roleId })
}

export default apiClient

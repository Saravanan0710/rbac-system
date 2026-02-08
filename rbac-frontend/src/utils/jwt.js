const TOKEN_KEY = 'auth_token'

export function setAuthToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function getAuthToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function removeAuthToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export function isTokenValid(token) {
  if (!token) return false
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const expirationTime = payload.exp * 1000
    return Date.now() < expirationTime
  } catch (error) {
    console.error('Invalid token:', error)
    return false
  }
}

export function decodeToken(token) {
  if (!token) return null
  
  try {
    return JSON.parse(atob(token.split('.')[1]))
  } catch (error) {
    console.error('Failed to decode token:', error)
    return null
  }
}

export function getUserFromToken(token) {
  const payload = decodeToken(token)
  if (!payload) return null
  
  return {
    id: payload.sub || payload.user_id,
    email: payload.email,
    name: payload.name,
    roles: payload.roles || [],
    role: payload.role || (payload.roles && payload.roles[0]) || null
  }
}

export function getUserRoleFromToken(token) {
  const user = getUserFromToken(token)
  if (!user) return null
  return user.role || (Array.isArray(user.roles) && user.roles[0]) || null
}

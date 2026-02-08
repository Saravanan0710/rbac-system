import { useState, useEffect } from 'react'
import { getAuthToken, getUserFromToken } from '../utils/jwt'
import { usersAPI } from '../api/client'
import './Profile.css'

function Profile() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [profileImage, setProfileImage] = useState(null)
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const token = getAuthToken()
        if (!token) {
          setError('Not authenticated')
          return
        }

        // Get user info from JWT and fetch full profile from API
        const tokenData = getUserFromToken(token)
        setUser(tokenData)

        if (tokenData?.id) {
          try {
            // Try to get full profile from backend API
            const response = await usersAPI.getCurrentUser(tokenData.id)
            if (response?.data?.user) {
              // Merge API data with JWT data
              setUser((prev) => ({ ...prev, ...response.data.user }))
            }
          } catch (err) {
            // If API call fails, use JWT data (already set above)
            console.log('Using JWT data fallback for profile')
          }
        }
      } catch (err) {
        setError('Failed to load user profile')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchUserProfile()
  }, [])

  const handleImageUpload = async (e) => {
    const file = e.target.files?.[0]
    if (!file) return

    setUploading(true)
    try {
      // Create FormData for image upload
      const formData = new FormData()
      formData.append('image', file)
      
      // For now, just show preview
      const reader = new FileReader()
      reader.onload = (event) => {
        setProfileImage(event.target.result)
      }
      reader.readAsDataURL(file)
      
      setError('')
    } catch (err) {
      setError('Failed to upload image')
      console.error(err)
    } finally {
      setUploading(false)
    }
  }

  if (loading) {
    return (
      <div className="container" style={{ minHeight: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner"></div>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="container">
        <div className="alert alert-error">User information not available</div>
      </div>
    )
  }

  return (
    <div className="profile-page">
      <div className="container">
        <h1>My Profile</h1>

        <div className="profile-container">
          <div className="profile-card">
            <div className="profile-header">
              <div className="profile-image-section">
                <div className="profile-image">
                  {profileImage ? (
                    <img src={profileImage} alt="Profile" />
                  ) : (
                    <div className="avatar-placeholder">
                      {user?.name?.charAt(0).toUpperCase() || 'U'}
                    </div>
                  )}
                </div>
                <label className="image-upload-btn">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleImageUpload}
                    disabled={uploading}
                    style={{ display: 'none' }}
                  />
                  📸 {uploading ? 'Uploading...' : 'Change Photo'}
                </label>
              </div>
              <div className="profile-info">
                <h2>{user?.name || 'User'}</h2>
                <p className="email">{user?.email || 'N/A'}</p>
              </div>
            </div>

            {error && <div className="alert alert-error">{error}</div>}

            <div className="profile-details">
              <div className="detail-section">
                <h3>Personal Information</h3>
                <div className="details-grid">
                  <div className="detail-item">
                    <label>Email</label>
                    <p>{user?.email || 'N/A'}</p>
                  </div>
                  {user?.date_of_birth && (
                    <div className="detail-item">
                      <label>Date of Birth</label>
                      <p>{new Date(user.date_of_birth).toLocaleDateString()}</p>
                    </div>
                  )}
                  {user?.social && (
                    <div className="detail-item">
                      <label>Social Media</label>
                      <p>{user.social}</p>
                    </div>
                  )}
                </div>
              </div>

              <div className="detail-section">
                <h3>Work Information</h3>
                <div className="details-grid">
                  {user?.employee_id && (
                    <div className="detail-item">
                      <label>Employee ID</label>
                      <p>{user.employee_id}</p>
                    </div>
                  )}
                  {user?.date_of_joining && (
                    <div className="detail-item">
                      <label>Date of Joining</label>
                      <p>{new Date(user.date_of_joining).toLocaleDateString()}</p>
                    </div>
                  )}
                </div>
              </div>

              {user?.roles && user.roles.length > 0 && (
                <div className="detail-section">
                  <h3>Roles & Permissions</h3>
                  <div className="roles-list">
                    {user.roles.map((role) => (
                      <span key={role} className="role-badge">
                        {role}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Profile

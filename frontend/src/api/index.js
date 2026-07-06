import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export const authAPI = {
  login: (username, password) => api.post('/auth/login', { username, password }),
  register: (username, password, nickname) => api.post('/auth/register', { username, password, nickname }),
}

export const userAPI = {
  getProfile: () => api.get('/user/profile'),
  updateProfile: (data) => api.put('/user/profile', data),
}

export const matchAPI = {
  start: () => api.post('/match/start'),
  cancel: () => api.post('/match/cancel'),
  status: () => api.get('/match/status'),
}

export const roomAPI = {
  create: (name) => api.post('/rooms', { name }),
  list: () => api.get('/rooms'),
  get: (id) => api.get(`/rooms/${id}`),
  join: (id) => api.post(`/rooms/${id}/join`),
  leave: (id) => api.post(`/rooms/${id}/leave`),
}

export const lobbyAPI = {
  get: () => api.get('/lobby'),
}

export const historyAPI = {
  list: () => api.get('/history'),
  get: (id) => api.get(`/history/${id}`),
}

export default api

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '../api'
import wsClient from '../ws/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(username, password) {
    const res = await authAPI.login(username.trim(), password)
    token.value = res.token
    user.value = res.user
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    wsClient.connect()
    return res
  }

  async function register(username, password, nickname) {
    return await authAPI.register(username.trim(), password, nickname.trim())
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    wsClient.disconnect()
  }

  return { token, user, isLoggedIn, login, register, logout }
})

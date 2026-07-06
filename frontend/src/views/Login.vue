<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="brand">
        <h1>Liar's Bar</h1>
        <p class="subtitle">骗子酒馆 · 登录</p>
      </div>
      <form @submit.prevent="handleLogin">
        <div class="field">
          <label>用户名</label>
          <input v-model="username" type="text" placeholder="输入用户名" required />
        </div>
        <div class="field">
          <label>密码</label>
          <input v-model="password" type="password" placeholder="输入密码" required />
        </div>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '登录中...' : '进入酒馆' }}
        </button>
      </form>
      <p class="link">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </p>
    </div>

    <Transition name="toast">
      <div v-if="toast.show" class="toast" :class="toast.type">
        <span class="toast-icon">{{ toast.type === 'error' ? '⚠️' : '✓' }}</span>
        <span class="toast-msg">{{ toast.msg }}</span>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const toast = ref({ show: false, msg: '', type: 'error' })
let toastTimer = null

function showToast(msg, type = 'error') {
  toast.value = { show: true, msg, type }
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toast.value.show = false }, 3000)
}

async function handleLogin() {
  if (!username.value.trim() || !password.value) {
    showToast('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await authStore.login(username.value, password.value)
    router.push('/')
  } catch (e) {
    const msg = e.response?.data?.msg || '登录失败，请稍后重试'
    // 后端返回 "invalid credentials" 统一提示为更友好的文案
    const friendly = msg.includes('credential') ? '用户名或密码错误' : msg
    showToast(friendly)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  background:
    linear-gradient(rgba(15, 8, 4, 0.78), rgba(15, 8, 4, 0.92)),
    url('/images/login.png') center/cover no-repeat;
  background-color: #0d0805;
  position: relative;
}

.auth-page::before {
  content: '';
  position: fixed;
  inset: 0;
  background: radial-gradient(ellipse 70% 50% at 50% 40%,
    transparent 0%, transparent 40%, rgba(0,0,0,0.6) 100%);
  pointer-events: none;
}

.auth-card {
  position: relative;
  z-index: 1;
  background: linear-gradient(160deg, rgba(34, 22, 14, 0.92) 0%, rgba(20, 12, 8, 0.95) 100%);
  border: 1px solid #5a3a1e;
  border-radius: 14px;
  padding: 44px 40px 36px;
  width: 420px;
  max-width: 92%;
  text-align: center;
  box-shadow:
    0 20px 60px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(203, 151, 103, 0.15),
    inset 0 1px 0 rgba(214, 192, 169, 0.12);
}

.brand {
  margin-bottom: 28px;
}

h1 {
  color: #cb9767;
  font-size: 38px;
  font-weight: 800;
  letter-spacing: 2px;
  margin-bottom: 6px;
  text-shadow: 0 2px 12px rgba(203, 151, 103, 0.35);
}

.subtitle {
  color: #8a6a4a;
  font-size: 14px;
  letter-spacing: 4px;
}

.field {
  text-align: left;
  margin-bottom: 18px;
}

.field label {
  display: block;
  color: #d6c0a9;
  font-size: 13px;
  margin-bottom: 6px;
  letter-spacing: 1px;
}

input {
  width: 100%;
  padding: 12px 14px;
  background: rgba(13, 8, 5, 0.6);
  border: 1px solid #3a2616;
  border-radius: 8px;
  color: #d6c0a9;
  font-size: 14px;
  transition: border-color 0.2s, box-shadow 0.2s;
  box-sizing: border-box;
}

input::placeholder { color: #5a4434; }

input:focus {
  outline: none;
  border-color: #9a3827;
  box-shadow: 0 0 0 3px rgba(154, 56, 39, 0.2);
}

.btn-primary {
  width: 100%;
  background: linear-gradient(180deg, #9a3827 0%, #7a2a1d 100%);
  color: #f5e6d3;
  padding: 14px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
  border: 1px solid #b04a36;
  border-radius: 8px;
  margin-top: 8px;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 14px rgba(154, 56, 39, 0.35);
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(180deg, #b04a36 0%, #8a3322 100%);
  box-shadow: 0 6px 18px rgba(154, 56, 39, 0.5);
}

.btn-primary:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.error {
  color: #e94560;
  margin-top: 14px;
  font-size: 13px;
}

.link {
  margin-top: 22px;
  color: #8a6a4a;
  font-size: 13px;
}

.link a {
  color: #cb9767;
  text-decoration: none;
}
.link a:hover { text-decoration: underline; }

/* ---- Toast notification ---- */
.toast {
  position: fixed;
  top: 32px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 500;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 22px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  max-width: 90vw;
}

.toast.error {
  background: linear-gradient(160deg, rgba(120, 30, 25, 0.95) 0%, rgba(80, 18, 14, 0.95) 100%);
  border: 1px solid #e94560;
  color: #ffe4e0;
}

.toast.success {
  background: linear-gradient(160deg, rgba(40, 90, 55, 0.95) 0%, rgba(28, 60, 38, 0.95) 100%);
  border: 1px solid #4ade80;
  color: #d8ffe4;
}

.toast-icon { font-size: 18px; }

.toast-enter-active, .toast-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -16px);
}
</style>

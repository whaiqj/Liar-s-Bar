<template>
  <div class="match-page">
    <div class="match-card">
      <div class="spinner"></div>
      <h1>正在寻找对手...</h1>
      <p class="elapsed">等待时间: {{ elapsed }}s</p>
      <p class="hint" v-if="elapsed > 10">正在为您匹配AI对手</p>
      <button class="cancel-btn" @click="cancelMatch">取消匹配</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { matchAPI } from '../api'
import wsClient from '../ws/client'

const router = useRouter()

const elapsed = ref(0)
let timer = null

function onMatchFound(payload) {
  clearInterval(timer)
  router.push(`/room/${payload.room_id}`)
}

function onGameStarted() {
  clearInterval(timer)
}

onMounted(() => {
  timer = setInterval(() => elapsed.value++, 1000)

  wsClient.connect()
  wsClient.on('MATCH_FOUND', onMatchFound)
  wsClient.on('GAME_STARTED', onGameStarted)
})

onUnmounted(() => {
  clearInterval(timer)
  wsClient.off('MATCH_FOUND', onMatchFound)
  wsClient.off('GAME_STARTED', onGameStarted)
})

async function cancelMatch() {
  try {
    await matchAPI.cancel()
    clearInterval(timer)
    router.push('/')
  } catch (e) {
    console.error('Cancel failed:', e)
  }
}
</script>

<style scoped>
.match-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.match-card {
  text-align: center;
  padding: 60px 40px;
  background: rgba(22, 33, 62, 0.9);
  border: 1px solid #333;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.spinner {
  width: 60px;
  height: 60px;
  border: 4px solid #333;
  border-top-color: #e94560;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 24px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

h1 {
  color: #e0e0e0;
  margin-bottom: 8px;
}

.elapsed {
  color: #888;
  font-size: 18px;
  margin-bottom: 8px;
}

.hint {
  color: #e94560;
  margin-bottom: 24px;
}

.cancel-btn {
  background: #333;
  color: #e0e0e0;
  padding: 12px 32px;
}
</style>

<template>
  <div class="history-page">
    <header class="page-header">
      <button @click="$router.push('/')" class="back-btn">← 返回</button>
      <h1>历史战绩</h1>
    </header>

    <div class="history-list">
      <div v-if="games.length === 0" class="empty">暂无战绩记录</div>
      <div v-for="game in games" :key="game.id" class="game-card">
        <div class="game-time">{{ formatTime(game.start_time) }}</div>
        <div class="game-info">
          <span>{{ game.total_rounds }} 回合</span>
          <span>{{ game.total_turns }} 回合</span>
          <span v-if="game.ai_count > 0">{{ game.ai_count }} 个AI</span>
        </div>
        <div class="game-result">
          <span v-if="game.winner_id === authStore.user?.id" class="win">胜利</span>
          <span v-else class="lose">失败</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { historyAPI } from '../api'

const authStore = useAuthStore()
const games = ref([])

onMounted(async () => {
  try {
    const res = await historyAPI.list()
    if (res.data && Array.isArray(res.data)) {
      games.value = res.data
    }
  } catch (e) {
    console.error('Failed to load history:', e)
  }
})

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}
</script>

<style scoped>
.history-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-header h1 { color: #e94560; }

.back-btn {
  background: none;
  color: #888;
  padding: 6px 12px;
}

.empty {
  text-align: center;
  color: #666;
  padding: 40px;
}

.game-card {
  background: rgba(22, 33, 62, 0.8);
  border: 1px solid #333;
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.game-time { color: #888; font-size: 14px; }
.game-info { color: #aaa; font-size: 13px; display: flex; gap: 12px; }
.game-result .win { color: #4ade80; font-weight: 700; }
.game-result .lose { color: #f59e0b; font-weight: 700; }
</style>

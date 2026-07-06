<template>
  <div class="profile-page">
    <header class="page-header">
      <button @click="$router.push('/')" class="back-btn">← 返回</button>
      <h1>个人中心</h1>
    </header>

    <div class="profile-card" v-if="profile">
      <div class="avatar-section">
        <div class="avatar">{{ profile.nickname?.charAt(0) || '?' }}</div>
        <h2>{{ profile.nickname }}</h2>
        <p class="username">@{{ profile.username }}</p>
      </div>

      <div class="stats-grid">
        <div class="stat-item">
          <span class="stat-num">{{ profile.elo_rating }}</span>
          <span class="stat-lbl">ELO 评分</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ profile.total_games }}</span>
          <span class="stat-lbl">总局数</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ winRate }}%</span>
          <span class="stat-lbl">胜率</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ profile.total_wins }}</span>
          <span class="stat-lbl">胜场</span>
        </div>
      </div>

      <div class="detail-stats">
        <h3>详细数据</h3>
        <div class="detail-row">
          <span>总游戏</span><span>{{ profile.total_games }}</span>
        </div>
        <div class="detail-row">
          <span>胜利</span><span>{{ profile.total_wins }}</span>
        </div>
        <div class="detail-row">
          <span>失败</span><span>{{ profile.total_losses }}</span>
        </div>
        <div class="detail-row">
          <span>撒谎次数</span><span>{{ profile.total_lies }}</span>
        </div>
        <div class="detail-row">
          <span>质疑次数</span><span>{{ profile.total_challenges }}</span>
        </div>
        <div class="detail-row">
          <span>质疑成功</span><span>{{ profile.total_successful_challenges }}</span>
        </div>
      </div>

      <div class="edit-section">
        <h3>修改资料</h3>
        <input v-model="editNickname" placeholder="新昵称" />
        <button @click="updateProfile" :disabled="!editNickname">保存</button>
        <p v-if="updateMsg">{{ updateMsg }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { userAPI } from '../api'

const profile = ref(null)
const editNickname = ref('')
const updateMsg = ref('')

const winRate = computed(() => {
  if (!profile.value) return 0
  const total = profile.value.total_games
  if (total === 0) return 0
  return Math.round((profile.value.total_wins / total) * 100)
})

onMounted(async () => {
  try {
    const res = await userAPI.getProfile()
    profile.value = res.data
    editNickname.value = res.data?.nickname || ''
  } catch (e) {
    console.error('Failed to load profile:', e)
  }
})

async function updateProfile() {
  try {
    await userAPI.updateProfile({ nickname: editNickname.value })
    updateMsg.value = '保存成功'
    setTimeout(() => updateMsg.value = '', 2000)
  } catch (e) {
    updateMsg.value = '保存失败'
  }
}
</script>

<style scoped>
.profile-page {
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
  margin-bottom: 24px;
}

.page-header h1 { color: #e94560; }

.back-btn {
  background: none;
  color: #888;
  padding: 6px 12px;
}

.profile-card {
  background: rgba(22, 33, 62, 0.8);
  border: 1px solid #333;
  border-radius: 16px;
  padding: 32px;
}

.avatar-section {
  text-align: center;
  margin-bottom: 24px;
}

.avatar {
  width: 80px; height: 80px;
  background: #e94560;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0 auto 12px;
}

.avatar-section h2 {
  color: #e0e0e0;
  margin-bottom: 4px;
}

.username { color: #888; }

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 24px;
}

.stat-item {
  background: rgba(0,0,0,0.3);
  border-radius: 8px;
  padding: 12px 8px;
  text-align: center;
}

.stat-num {
  display: block;
  color: #e94560;
  font-size: 24px;
  font-weight: 700;
}

.stat-lbl {
  color: #888;
  font-size: 12px;
}

.detail-stats {
  margin-bottom: 24px;
}

.detail-stats h3, .edit-section h3 {
  color: #e0e0e0;
  margin-bottom: 12px;
  font-size: 16px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  color: #aaa;
  border-bottom: 1px solid #222;
}

.edit-section input {
  width: 100%;
  margin-bottom: 8px;
}

.edit-section button {
  background: #e94560;
  color: white;
}

.edit-section button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.edit-section p {
  color: #4ade80;
  margin-top: 8px;
  font-size: 14px;
}
</style>

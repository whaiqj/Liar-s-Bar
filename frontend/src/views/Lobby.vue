<template>
  <div class="lobby">
    <header class="lobby-header">
      <div class="brand">
        <h1>Liar's Bar</h1>
        <span class="sub">骗子酒馆</span>
      </div>
      <div class="header-info">
        <span class="online-badge">{{ lobbyData.online_count || 0 }} 在线</span>
        <span class="user-name">{{ authStore.user?.nickname || authStore.user?.username }}</span>
        <button class="rules-btn" @click="showRules = true">📖 规则</button>
        <button class="btn-sm" @click="logout">退出</button>
      </div>
    </header>

    <div class="stats-bar">
      <div class="stat">
        <span class="stat-value">{{ lobbyData.active_rooms?.length || 0 }}</span>
        <span class="stat-label">活跃房间</span>
      </div>
      <div class="stat">
        <span class="stat-value">{{ lobbyData.online_count || 0 }}</span>
        <span class="stat-label">在线玩家</span>
      </div>
      <div class="stat">
        <span class="stat-value">{{ lobbyData.queue_length || 0 }}</span>
        <span class="stat-label">匹配中</span>
      </div>
    </div>

    <div class="main-actions">
      <button class="action-btn match-btn" @click="startMatch">
        <span class="icon">⚔️</span>
        <span>快速匹配</span>
      </button>
      <button class="action-btn create-btn" @click="createRoom">
        <span class="icon">🏠</span>
        <span>创建房间</span>
      </button>
      <button class="action-btn profile-btn" @click="$router.push('/profile')">
        <span class="icon">👤</span>
        <span>个人中心</span>
      </button>
    </div>

    <div class="rooms-section">
      <h2>活跃房间</h2>
      <div v-if="rooms.length === 0" class="empty">暂无活跃房间</div>
      <div
        v-for="room in rooms"
        :key="room.id"
        class="room-card"
        :class="{ disabled: !roomCanJoin(room) }"
        @click="roomCanJoin(room) && joinRoom(room.id)"
      >
        <div class="room-main">
          <div class="room-name">{{ roomTitle(room) }}</div>
          <div class="room-ready">{{ room.ready_count || 0 }} 人已准备</div>
        </div>
        <div class="room-info">
          <span class="room-players">{{ roomPlayers(room) }}/{{ room.max_players || 4 }}</span>
          <span class="room-status" :class="roomPhase(room).toLowerCase()">
            {{ roomPhase(room) === 'PLAYING' ? '游戏中' : '等待中' }}
          </span>
          <button class="join-btn" :disabled="!roomCanJoin(room)" @click.stop="joinRoom(room.id)">
            {{ roomCanJoin(room) ? '加入' : '不可加入' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Game Rules modal -->
    <div v-if="showRules" class="rules-overlay" @click.self="showRules = false">
      <div class="rules-modal">
        <button class="rules-close" @click="showRules = false">×</button>
        <h2>📖 骗子酒馆 规则</h2>
        <div class="rules-body">
          <section>
            <h3>目标</h3>
            <p>4 人对战，通过出牌、撒谎、质疑和心理博弈，成为最后存活的玩家。</p>
          </section>
          <section>
            <h3>牌组</h3>
            <p>A / K / Q / J 各 6 张，共 24 张。每人发 6 张。</p>
          </section>
          <section>
            <h3>目标牌</h3>
            <p>每一轮有指定目标牌（A→K→Q→J 循环）。出牌时必须声明打出的全部是当前目标牌。</p>
          </section>
          <section>
            <h3>出牌</h3>
            <p>当前回合玩家选 1~3 张手牌打出。可以真出（全是目标牌），也可以虚张声势（含非目标牌）。其他玩家只能看到声明，看不到真实牌面。</p>
          </section>
          <section>
            <h3>质疑</h3>
            <p>下一名玩家可选择质疑上家。系统翻开上家刚打出的牌：</p>
            <ul>
              <li>含非目标牌（撒谎）→ 质疑成功，撒谎者受罚</li>
              <li>全是目标牌（说真话）→ 质疑失败，质疑者受罚</li>
            </ul>
          </section>
          <section>
            <h3>俄罗斯轮盘惩罚</h3>
            <p>每次受罚累计子弹数：第 1 次 1 发，第 2 次 2 发……第 6 次必死。被击中即淘汰出局，存活则继续。</p>
          </section>
          <section>
            <h3>手牌耗尽</h3>
            <p>所有存活玩家手牌为空时，重新洗牌发牌，目标牌切换到下一种。</p>
          </section>
          <section>
            <h3>操作</h3>
            <ul>
              <li><b>出牌</b>：选 1~3 张牌，点击「出牌」</li>
              <li><b>质疑上家</b>：仅当上家有出牌时可点</li>
              <li><b>过</b>：不质疑，交由下家</li>
            </ul>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { lobbyAPI, roomAPI, matchAPI } from '../api'
import wsClient from '../ws/client'

const router = useRouter()
const authStore = useAuthStore()

const lobbyData = ref({ online_count: 0, active_rooms: [] })
const rooms = ref([])
const showRules = ref(false)
let pollTimer = null

onMounted(async () => {
  wsClient.on('connected', fetchLobby)
  wsClient.connect()
  await fetchLobby()
  pollTimer = setInterval(fetchLobby, 5000)
})

onUnmounted(() => {
  wsClient.off('connected', fetchLobby)
  if (pollTimer) clearInterval(pollTimer)
})

async function fetchLobby() {
  try {
    const res = await lobbyAPI.get()
    lobbyData.value = res.data
    rooms.value = res.data?.active_rooms || []
  } catch (e) {
    console.error('Failed to fetch lobby:', e)
  }
}

async function startMatch() {
  try {
    await matchAPI.start()
    router.push('/match')
  } catch (e) {
    console.error('Match start failed:', e)
  }
}

async function createRoom() {
  try {
    const res = await roomAPI.create('New Room')
    router.push(`/room/${res.data.id}`)
  } catch (e) {
    const msg = e.response?.data?.msg || '创建房间失败'
    alert(msg)
    console.error('Create room failed:', e)
  }
}

async function joinRoom(roomId) {
  try {
    await roomAPI.join(roomId)
    router.push(`/room/${roomId}`)
  } catch (e) {
    const msg = e.response?.data?.msg || '加入房间失败'
    alert(msg)
    console.error('Join room failed:', e)
    fetchLobby()
  }
}

function roomTitle(room) {
  return room.name || room.room_name || '无名房间'
}

function roomPlayers(room) {
  return room.players ?? room.current_players ?? 0
}

function roomPhase(room) {
  return room.phase || room.room_status || 'WAITING'
}

function roomCanJoin(room) {
  return room.can_join ?? (roomPlayers(room) < (room.max_players || 4) && roomPhase(room) !== 'PLAYING')
}

function logout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.lobby {
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

/* Full-viewport background image layer */
.lobby::after {
  content: '';
  position: fixed;
  inset: 0;
  background:
    linear-gradient(rgba(13, 8, 4, 0.82), rgba(13, 8, 4, 0.94)),
    url('/images/lobby.png') center/cover no-repeat;
  background-color: #0d0805;
  z-index: -2;
}

.lobby::before {
  content: '';
  position: fixed;
  inset: 0;
  background: radial-gradient(ellipse 75% 60% at 50% 40%,
    transparent 0%, transparent 45%, rgba(0,0,0,0.55) 100%);
  pointer-events: none;
  z-index: -1;
}

/* Center content column on wide screens */
.lobby-header,
.stats-bar,
.main-actions,
.rooms-section {
  max-width: 900px;
  margin-left: auto;
  margin-right: auto;
}

.lobby > * { position: relative; z-index: 1; }

.lobby-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #3a2616;
  margin-bottom: 20px;
}

.brand { display: flex; align-items: baseline; gap: 10px; }

.lobby-header h1 {
  color: #cb9767;
  font-size: 30px;
  font-weight: 800;
  letter-spacing: 2px;
  text-shadow: 0 2px 12px rgba(203, 151, 103, 0.35);
}

.sub {
  font-size: 14px;
  color: #8a6a4a;
  letter-spacing: 3px;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.online-badge {
  background: rgba(90, 46, 23, 0.5);
  border: 1px solid #5a2e17;
  color: #cb9767;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
}

.user-name {
  color: #d6c0a9;
  font-size: 14px;
}

.rules-btn {
  background: rgba(154, 56, 39, 0.2);
  border: 1px solid #9a3827;
  color: #cb9767;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.rules-btn:hover {
  background: rgba(154, 56, 39, 0.4);
  border-color: #b04a36;
}

.btn-sm {
  background: rgba(58, 38, 22, 0.6);
  border: 1px solid #3a2616;
  color: #d6c0a9;
  padding: 6px 16px;
  font-size: 13px;
  border-radius: 8px;
  cursor: pointer;
}
.btn-sm:hover { background: rgba(154, 56, 39, 0.3); }

.stats-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.stat {
  flex: 1;
  background: linear-gradient(160deg, rgba(34, 22, 14, 0.85) 0%, rgba(20, 12, 8, 0.9) 100%);
  border: 1px solid #3a2616;
  border-radius: 12px;
  padding: 16px;
  text-align: center;
  box-shadow: inset 0 1px 0 rgba(214, 192, 169, 0.08);
}

.stat-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: #cb9767;
}

.stat-label {
  color: #8a6a4a;
  font-size: 14px;
  letter-spacing: 1px;
}

.main-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 32px;
}

.action-btn {
  background: linear-gradient(160deg, rgba(34, 22, 14, 0.85) 0%, rgba(20, 12, 8, 0.9) 100%);
  border: 1px solid #3a2616;
  border-radius: 12px;
  padding: 22px 12px;
  color: #d6c0a9;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: inset 0 1px 0 rgba(214, 192, 169, 0.08);
}

.action-btn:hover {
  border-color: #9a3827;
  background: linear-gradient(160deg, rgba(90, 46, 23, 0.55) 0%, rgba(34, 22, 14, 0.9) 100%);
  transform: translateY(-2px);
}

.icon { font-size: 26px; }

.rooms-section h2 {
  color: #cb9767;
  margin-bottom: 12px;
  letter-spacing: 2px;
}

.empty {
  color: #6a5238;
  text-align: center;
  padding: 40px;
  background: rgba(20, 12, 8, 0.6);
  border: 1px dashed #3a2616;
  border-radius: 12px;
}

.room-card {
  background: linear-gradient(160deg, rgba(34, 22, 14, 0.85) 0%, rgba(20, 12, 8, 0.9) 100%);
  border: 1px solid #3a2616;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: inset 0 1px 0 rgba(214, 192, 169, 0.08);
}

.room-card:hover {
  border-color: #9a3827;
  background: linear-gradient(160deg, rgba(90, 46, 23, 0.5) 0%, rgba(20, 12, 8, 0.9) 100%);
}

.room-card.disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.room-main { min-width: 0; }

.room-name {
  color: #d6c0a9;
  font-weight: 600;
  overflow-wrap: anywhere;
}

.room-ready {
  color: #6a5238;
  font-size: 12px;
  margin-top: 4px;
}

.room-info {
  display: flex;
  gap: 12px;
  align-items: center;
}

.room-players { color: #8a6a4a; font-size: 14px; }

.room-status {
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
}

.room-status.playing {
  background: rgba(154, 56, 39, 0.25);
  color: #e94560;
  border: 1px solid #9a3827;
}

.room-status.waiting {
  background: rgba(203, 151, 103, 0.15);
  color: #cb9767;
  border: 1px solid #5a3a1e;
}

.join-btn {
  background: linear-gradient(180deg, #9a3827 0%, #7a2a1d 100%);
  color: #f5e6d3;
  padding: 7px 16px;
  font-size: 13px;
  border: 1px solid #b04a36;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.join-btn:hover:not(:disabled) {
  background: linear-gradient(180deg, #b04a36 0%, #8a3322 100%);
}

.join-btn:disabled {
  background: #2a1a10;
  color: #5a4434;
  border-color: #3a2616;
  cursor: not-allowed;
}

/* ---- Rules modal ---- */
.rules-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.rules-modal {
  position: relative;
  background: linear-gradient(160deg, #1f1610 0%, #2a1c12 100%);
  border: 1px solid #5a3a1e;
  border-radius: 14px;
  max-width: 520px;
  width: 100%;
  max-height: 85vh;
  overflow-y: auto;
  padding: 28px 28px 24px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.6), 0 0 0 1px rgba(203, 151, 103, 0.2);
}
.rules-modal h2 {
  color: #cb9767;
  font-size: 22px;
  margin-bottom: 16px;
  text-align: center;
}
.rules-close {
  position: absolute;
  top: 10px;
  right: 14px;
  background: none;
  border: none;
  color: #8a6a4a;
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
}
.rules-close:hover { color: #e94560; }
.rules-body section {
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid #3a2616;
}
.rules-body section:last-child { border-bottom: none; margin-bottom: 0; }
.rules-body h3 {
  color: #f5a623;
  font-size: 15px;
  margin-bottom: 6px;
}
.rules-body p, .rules-body ul {
  color: #d0c8c0;
  font-size: 13px;
  line-height: 1.6;
  margin: 0;
}
.rules-body ul { padding-left: 18px; }
.rules-body li { margin-bottom: 3px; }
.rules-body b { color: #e94560; }

@media (max-width: 600px) {
  .main-actions {
    grid-template-columns: repeat(2, 1fr);
  }
  .lobby-header h1 { font-size: 24px; }
  .sub { display: none; }
}
</style>

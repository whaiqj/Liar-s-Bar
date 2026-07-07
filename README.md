#  Liar's Bar Online

> 基于 Go + WebSocket + 强化学习 AI 的多人在线欺骗博弈游戏平台

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)](https://vuejs.org/)
[![Python](https://img.shields.io/badge/Python-3.10+-3776AB?logo=python)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)](https://docs.docker.com/compose/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 项目简介

**Liar's Bar Online**（骗子酒馆在线）是一个 4 人制多人在线实时欺骗博弈游戏平台，参考知名派对游戏《骗子酒馆（Liar's Bar）》设计。

玩家轮流声明出牌，可以**说真话**也可以**撒谎**，下家可选择**质疑**揭穿谎言。质疑失败或撒谎被揭穿的一方将面临**俄罗斯轮盘**惩罚——累计 6 发子弹即被淘汰。最后存活者获胜。

系统支持**真人玩家**与**强化学习 AI** 混合对局，AI 补位机制确保玩家无需等待即可开局。

---

##  核心玩法

| 要素 | 说明 |
|------|------|
| **牌组** | A / K / Q / J 各 6 张，共 24 张 |
| **人数** | 固定 4 人，不足时 AI 补位 |
| **出牌** | 每回合 1-3 张，必须**声明**为目标牌（可以说谎） |
| **质疑** | 下家可质疑上家声明 → 翻牌验证 → 撒谎方受罚 |
| **惩罚** | 俄罗斯轮盘：第 N 次 N 发子弹，6 发必死 |
| **目标牌** | 按 A→K→Q→J 循环，手牌耗尽后重发切换 |
| **胜利** | 最后存活者获胜 |

###  四大角色

| 角色 | 特殊能力 |
|------|----------|
| **Scubby** | 每轮额外获得 1 张 WILD 万能牌 |
| **Foxy** | 技能：偷看一名玩家所有手牌（整局 1 次） |
| **Bristle** | 每轮可质疑 2 次（其余角色限 1 次） |
| **Tor** | 50% 概率跳过惩罚累加；30% 概率免疫致命子弹 |

---

##  技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Nginx :8081                          │
│              反向代理 (HTTP + WebSocket)                      │
├──────────────┬──────────────────┬────────────────────────────┤
│   Frontend   │     Backend      │        AI Service          │
│   Vue 3      │     Go + Gin     │     Python + FastAPI       │
│   :3000      │     :8080        │        :8000               │
│              │                  │                             │
│   Vite       │  ┌────────────┐  │   ┌───────────────────┐   │
│   Pinia      │  │ Game Engine │  │   │    PPO Agent      │   │
│   Router     │  │ (状态机)    │◄─┼───┤ 64D obs → 5D act  │   │
│   Axios      │  └────────────┘  │   └───────────────────┘   │
│   WebSocket  │  ┌────────────┐  │                             │
│              │  │ Hub + Room  │  │                             │
│              │  │ (Goroutine) │  │                             │
│              │  └────────────┘  │                             │
│              │  ┌────────────┐  │                             │
│              │  │ Matchmaker │  │                             │
│              │  └────────────┘  │                             │
├──────────────┴──────────────────┴────────────────────────────┤
│              MySQL :3306           Redis :6379                │
│           (持久化数据)           (在线状态/队列)                │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈详情

| 层次 | 技术 | 关键组件 |
|------|------|----------|
| **前端** | Vue 3 + Vite + Pinia + Vue Router + Axios + 原生 WebSocket | 7 个页面视图 |
| **后端** | Go 1.21 + Gin + Gorilla WebSocket + GORM + JWT + bcrypt | Controller → Service → Repository |
| **数据库** | MySQL 8.0 (GORM AutoMigrate) + Redis 7 | 8 张数据表 + 在线状态缓存 |
| **AI** | Python + FastAPI + NumPy + 自实现 PPO | Actor-Critic 双网络 (128→64) |
| **部署** | Docker Compose (6 容器) + Nginx Alpine | 一键启动 |

### 后端架构关键决策

| 决策 | 选择 | 权衡 |
|------|------|------|
| 并发模型 | 每房间独立 Goroutine + Channel 事件驱动 | 串行化简单可靠，无锁竞争 |
| 通信协议 | WebSocket 长连接 + HTTP REST | 实时性极好，断线需重连 |
| 游戏权威 | 服务端判定所有逻辑 | 防作弊，但服务端压力大 |
| AI 架构 | Python 独立服务 + HTTP 调用 | 与 Go 解耦，增加网络延迟 |

---

##  快速启动

### 前置要求

- Docker Desktop 20.10+
- 可用端口：8081 (Nginx)、3307 (MySQL)、6379 (Redis)

### 一键启动

```bash
# 启动所有服务
./run.sh start

# 查看状态
./run.sh status

# 查看日志
./run.sh logs

# 停止
./run.sh stop

# 清理（含数据库数据）
./run.sh clean
```

或直接使用 Docker Compose：

```bash
docker compose up -d
```

### 访问地址

| 服务 | 地址 |
|------|------|
| 游戏前端 | http://localhost:8081 |
| 后端 API | http://localhost:8081/api/v1 |
| WebSocket | ws://localhost:8081/ws |
| AI 服务 | http://localhost:8000 |
| 直接后端 | http://localhost:8082 |

---

##  项目结构

```
Liar-s-Bar/
├── frontend/                     # Vue 3 前端
│   ├── src/
│   │   ├── views/                # 页面: Login/Register/Lobby/MatchWait/GameRoom/Profile/History
│   │   ├── stores/               # Pinia: auth.js, game.js
│   │   ├── api/index.js          # Axios HTTP 封装
│   │   ├── ws/client.js          # WebSocket 单例客户端
│   │   └── router/index.js       # Vue Router 路由守卫
│   ├── package.json
│   └── vite.config.js
│
├── backend/                      # Go 后端
│   ├── cmd/server/main.go        # 入口: 初始化+路由注册
│   └── internal/
│       ├── config/               # 环境变量配置
│       ├── controller/           # HTTP 控制器 (Auth/User/Room/Match/Lobby/Admin)
│       ├── service/              # 业务逻辑 (UserService/RoomService)
│       ├── repository/           # 数据访问层
│       ├── model/                # GORM 数据模型
│       ├── middleware/           # JWT + CORS
│       ├── websocket/            # Hub (连接管理) + Room (事件处理)
│       ├── game/engine.go        # 游戏状态机 (出牌/质疑/惩罚/角色技能)
│       ├── match/matchmaker.go   # 匹配引擎 (真人优先 → AI补位)
│       └── utils/redis.go        # Redis 客户端
│
├── ai-service/                   # Python AI 服务
│   ├── main.py                   # FastAPI + PPOAgent + LiarsBarEnv
│   └── requirements.txt
│
├── deploy/                       # Docker 构建文件
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   ├── Dockerfile.ai
│   └── nginx.conf
│
├── docker-compose.yml            # 6 容器编排
├── run.sh                        # 管理脚本
│
├── requirements.md               # 完整需求规格
├── architecture.md               # 系统架构设计
├── game_rules.md                 # 游戏规则详细说明
├── api.md                        # API 接口文档
├── database.md                   # 数据库设计
├── mobile-api.md                 # 移动端 API
└── android-ui-prompts.md         # Android UI 设计提示
```

---

##  当前进度

### 已完成功能

#### 用户系统
- [x] 账号密码注册 / 登录（JWT + bcrypt）
- [x] 个人信息查看 / 修改（昵称、头像）
- [x] ELO 评分系统（胜 +20 / 负 -15）
- [x] 玩家统计（总局数/胜场/败场/撒谎/质疑）

#### 匹配系统
- [x] 快速匹配队列（1 秒扫描周期）
- [x] 真人优先匹配（0-10 秒纯真人）
- [x] AI 补位机制（10 秒后允许 AI + 15 秒强制补齐）
- [x] 房间创建 / 加入 / 离开
- [x] 4 人满员自动开局

#### 游戏引擎
- [x] 完整状态机（WAITING → PLAYING → CHALLENGE → PUNISHMENT → GAME_OVER）
- [x] 发牌 / 洗牌 / 出牌合法性校验
- [x] 质疑判定（翻牌验证 / 撒谎识别）
- [x] 俄罗斯轮盘惩罚（累计子弹 / 淘汰）
- [x] 4 角色 & 差异化技能（Scubby/Foxy/Bristle/Tor）
- [x] 手牌耗尽自动重发 + 目标牌切换
- [x] 服务端权威判定（客户端不可作弊）

#### WebSocket 通信
- [x] 全局 Hub 连接管理 + 心跳检测
- [x] 房间独立 Goroutine 事件驱动
- [x] 游戏状态实时广播（每人仅见自己手牌）
- [x] 聊天消息广播
- [x] 断线自动退出 & 游戏终止
- [x] 重连恢复

#### AI 玩家
- [x] 规则基础 AI（本地备用策略）
- [x] PPO 强化学习推理接口（Python FastAPI）
- [x] 自博弈训练接口（LiarsBarEnv）
- [x] AI 聊天消息生成

#### 前端
- [x] 暗黑酒馆主题 UI（登录/注册/大厅/游戏/个人中心/历史）
- [x] 实时手牌展示 + 选牌出牌操作
- [x] 对手面板（角色/手牌数/惩罚子弹/存活状态）
- [x] 质疑动画（成功红闪 + 失败蓝闪）
- [x] 淘汰动画（💀 头骨下坠 + 屏幕震动）
- [x] 俄罗斯轮盘动画
- [x] 游戏规则弹窗
- [x] 聊天面板
- [x] 游戏结束结算覆盖层
- [x] 准备房间（角色选择 / 4 人准备）
- [x] 响应式布局（桌面 + 移动端适配）

#### 部署
- [x] Docker Compose 6 容器编排（MySQL/Redis/Backend/AI/Frontend/Nginx）
- [x] Nginx 反向代理（HTTP + WebSocket）
- [x] 一键启动/停止/状态管理脚本
- [x] 数据库健康检查等待依赖

---

###  部分实现 / 待完善

| 功能 | 当前状态 | 待完成 |
|------|----------|--------|
| **AI 服务集成** | `AIProxy.GetAction()` 硬编码返回 PASS，从未真正调用 Python AI 服务 | Go ↔ Python HTTP 调用未接通 |
| **游戏历史** | `GET /history` 返回硬编码空数组 `[]`，`GameRepo` 方法存在但未被调用 | 需实现 Handler 查询 + 前端详情页 |
| **游戏行为日志** | `RecordAction()` 方法体为空，`game_actions` 表已有但无数据写入 | 需实现操作日志持久化 |
| **聊天记录持久化** | 消息只广播不存储，`chat_records` 表 + `CreateChat` 方法已定义 | 需在 `handleChat` 中调用写入 |
| **对局记录持久化** | `GameRecordID` 字段存在但从未赋值，`GameRepo.Create` 未被调用 | 需在游戏结束时写入 `games` / `game_players` 表 |
| **匹配队列持久化** | 纯内存 `[]MatchEntry` 切片，服务重启丢失 | 需利用 Redis 或 `matchmaking_queue` 表 |
| **AI 模型版本管理** | `ai_models` 表已定义但未使用，PPO 权重随机初始化 | 训练后保存模型 + 版本切换 |
| **ELO 算法** | 当前为简化 ±20/±15 固定值，非标准 ELO（期望胜率×K 因子） | 可改为正式 ELO 公式 |
| **前端开发环境** | `vite.config.js` 存在但使用默认配置 | 确认代理 /api → :8080, /ws → ws://localhost:8080 等

---

###  尚未实现

#### 用户模块
- [ ] 邮箱注册（需求中提及，DB 字段已预留）
- [ ] 修改密码
- [ ] 头像上传（目前仅支持 URL 文本）
- [ ] 存活率 / 撒谎成功率 / 质疑成功率 统计展示（数据已记录但未计算比率）
- [ ] 防重复登录

#### 房间模块
- [ ] 私人房间（密码保护）
- [ ] 好友邀请
- [ ] 踢出玩家
- [ ] 观战模式
- [ ] 房间搜索 / 过滤

#### 游戏模块
- [ ] 游戏回放系统
- [ ] 骰子模式（扩展玩法）

#### 社交模块
- [ ] 好友系统
- [ ] 语音聊天

#### 管理后台
- [ ] 封禁用户
- [ ] 删除房间
- [ ] 查看日志
- [ ] 管理 AI 模型版本
- [ ] 服务器状态监控

#### 安全
- [ ] HTTPS 证书配置
- [ ] 防重复登录
- [ ] 接口限流
- [ ] 防作弊深度校验

#### 运维
- [ ] 日志聚合系统
- [ ] 监控告警（Prometheus/Grafana）
- [ ] CI/CD 流水线

---

##  API 概览

### REST API (`/api/v1`)

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| POST | `/auth/register` | 注册 | ✅ |
| POST | `/auth/login` | 登录 | ✅ |
| GET | `/user/profile` | 个人信息 | ✅ |
| PUT | `/user/profile` | 修改资料 | ✅ |
| POST | `/match/start` | 开始匹配 | ✅ |
| POST | `/match/cancel` | 取消匹配 | ✅ |
| GET | `/match/status` | 匹配状态 | ✅ |
| POST | `/rooms` | 创建房间 | ✅ |
| GET | `/rooms` | 房间列表 | ✅ |
| GET | `/rooms/:id` | 房间详情 | ✅ |
| POST | `/rooms/:id/join` | 加入房间 | ✅ |
| POST | `/rooms/:id/leave` | 离开房间 | ✅ |
| GET | `/lobby` | 大厅数据 | ✅ |
| GET | `/history` | 历史对局 | 🔶 空实现 |
| GET | `/admin/online` | 在线人数 | ✅ |
| GET | `/admin/rooms` | 房间列表 | ✅ |

### WebSocket 事件

**客户端 → 服务端：** `PLAYER_JOIN` / `SET_CHARACTER` / `PLAYER_READY` / `PLAY_CARD` / `CHALLENGE` / `PASS` / `CHAT` / `USE_SKILL`

**服务端 → 客户端：** `ROOM_STATE` / `GAME_STARTED` / `GAME_STATE` / `CHALLENGE_RESULT` / `RUSSIAN_ROULETTE` / `PLAYER_ELIMINATED` / `GAME_OVER` / `PLAYER_LEFT` / `CHAT` / `SKILL_RESULT` / `ERROR`

### AI 接口

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| GET | `/health` | 健康检查 | ✅ |
| GET | `/ai/status` | 模型状态 | ✅ |
| POST | `/ai/inference` | 获取 AI 动作 | ✅ |
| POST | `/ai/train` | 自博弈训练 | ✅ |

---

##  AI 强化学习详情

AI 服务使用 **PPO（Proximal Policy Optimization）** 算法：

- **观测空间**（64 维）：目标牌 (4D) + 手牌分布 (20D) + 游戏状态 (7D) + 玩家统计 (4D) + 补齐
- **动作空间**（5 维）：`TruthPlay` / `LiePlay` / `Challenge` / `Pass` / `Chat`
- **网络结构**：Actor (64→128→64→5) + Critic (64→128→64→1)
- **奖励函数**：胜利 +100 / 成功撒谎 +20 / 识破谎言 +20 / 错误质疑 -15 / 撒谎被揭穿 -20 / 淘汰 -100

> ⚠️ 当前 PPO 使用随机权重初始化，需运行训练获得有意义的策略模型。

---

##  超出需求的实现亮点

以下功能在原始需求文档中未定义，但在代码中已完整实现：

- **四大角色系统**：Scubby (WILD 牌)、Foxy (偷看手牌)、Bristle (双倍质疑)、Tor (惩罚免疫) — 每个角色都有差异化技能
- **暗黑酒馆主题 UI**：完整的视觉风格（木纹墙面 / 绿绒桌面 / 吊灯光晕闪烁 / 屏幕震动 / 头骨动画）
- **质疑 & 淘汰动画**：成功红闪光环 / 失败蓝闪 / 💀 头骨下坠 + 屏幕震动
- **游戏规则弹窗**：前端内置完整规则说明

---

##  贡献

本项目处于活跃开发阶段，欢迎贡献代码或提出 Issue。

### 本地开发

```bash
# 后端
cd backend && go run cmd/server/main.go

# 前端
cd frontend && npm install && npm run dev

# AI 服务
cd ai-service && pip install -r requirements.txt && python main.py
```

---

##  许可证

MIT License

---

##  相关文档

- [需求规格](requirements.md) — 完整功能需求
- [架构设计](architecture.md) — 系统架构决策
- [游戏规则](game_rules.md) — 详细游戏机制
- [API 文档](api.md) — REST + WebSocket 接口
- [数据库设计](database.md) — MySQL Schema
- [移动端 API](mobile-api.md) — 移动客户端接口

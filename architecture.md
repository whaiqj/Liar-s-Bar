# architecture.md

# 系统架构设计文档

## 1. 总体架构

本项目采用前后端分离架构。

前端使用Vue3实现页面交互。

后端使用Go语言实现核心服务。

HTTP接口用于用户、房间、战绩等普通业务。

WebSocket用于游戏内实时通信。

AI训练服务独立部署，使用Python实现强化学习训练与推理。

MySQL用于持久化数据。

Redis用于在线状态、匹配队列和临时缓存。

---

## 2. 技术栈

前端：

Vue3

Vite

Pinia

Vue Router

Axios

WebSocket API

后端：

Go

Gin

Gorilla WebSocket

GORM

JWT

MySQL

Redis

AI服务：

Python

Gymnasium

Stable-Baselines3

PPO

gRPC或HTTP

部署：

Docker

Docker Compose

Nginx

---

## 3. 系统模块

系统分为以下模块：

前端模块

用户认证模块

匹配模块

房间模块

游戏状态机模块

WebSocket网关模块

AI服务模块

数据库模块

日志模块

---

## 4. 前端架构

前端页面包括：

登录页

注册页

游戏大厅页

匹配等待页

游戏房间页

个人中心页

历史战绩页

前端主要职责是：

展示游戏状态

发送用户操作

维护WebSocket连接

展示聊天消息

展示战绩数据

前端不能保存关键游戏逻辑。

---

## 5. 后端架构

后端分为三层：

Controller层

Service层

Repository层

Controller负责处理HTTP请求。

Service负责业务逻辑。

Repository负责数据库读写。

WebSocket消息不直接写在Controller中，而是交给Hub和Room处理。

---

## 6. 推荐后端目录结构

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── controller/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── middleware/
│   ├── websocket/
│   ├── game/
│   ├── match/
│   ├── ai/
│   └── utils/
├── pkg/
├── migrations/
├── go.mod
└── go.sum
```

---

## 7. WebSocket架构

系统维护一个全局WebSocket Hub。

Hub负责管理所有连接。

每个连接对应一个Client对象。

Client包含：

UserID

Conn

SendChannel

CurrentRoomID

OnlineStatus

每个房间拥有独立Room对象。

Room负责管理房间内玩家和游戏状态。

---

## 8. 房间并发模型

每个房间对应一个独立Goroutine。

Room Goroutine只处理本房间事件。

所有玩家操作都转化为事件发送到Room Channel。

Room串行处理事件，避免复杂锁竞争。

事件类型包括：

JoinRoomEvent

LeaveRoomEvent

PlayerReadyEvent

StartGameEvent

PlayCardEvent

ChallengeEvent

ChatEvent

ReconnectEvent

AIActionEvent

GameOverEvent

---

## 9. Room结构示例

```go
type Room struct {
    ID        string
    Players   map[string]*Player
    State     *GameState
    Events    chan GameEvent
    Broadcast chan Message
}
```

Room不允许外部直接修改状态。

外部只能通过Events通道发送事件。

---

## 10. 游戏状态机

游戏状态包括：

WAITING

MATCHED

PLAYING

CHALLENGE

PUNISHMENT

ROUND_END

GAME_OVER

状态流转如下：

WAITING → MATCHED → PLAYING → CHALLENGE → PUNISHMENT → PLAYING

当只剩一名玩家存活时：

PLAYING → GAME_OVER

当游戏进行中有玩家主动退出时：

PLAYING → GAME_OVER（无赢家或剩余存活者获胜，房间在所有真人离开后销毁）

---

## 11. 匹配系统架构

匹配系统维护一个全局匹配队列。

玩家点击快速匹配后加入队列。

匹配服务定时扫描队列。

优先匹配真人玩家。

若等待时间超过10秒，允许加入AI。

若等待时间超过15秒，自动补齐AI并创建房间。

匹配完成后，系统通过 `RoomService.CreateRoom` 持久化房间到 MySQL（获得唯一自增 ID），再创建内存中的 `GameRoom` 实例并启动 Room Goroutine，避免房间 ID 冲突。

---

## 12. AI服务架构

AI服务分为训练模式和推理模式。

训练模式用于自我博弈。

推理模式用于线上游戏。

Go游戏服务器向AI服务发送当前观察状态。

AI服务返回动作。

请求格式包含：

RoomID

AIPlayerID

Observation

LegalActions

返回格式包含：

ActionType

Cards

ChatMessage

Confidence

---

## 13. AI调用流程

当前回合玩家为AI时：

Go服务构造Observation。

Go服务调用AI推理接口。

AI服务返回Action。

Go服务验证Action合法性。

Go服务将Action转化为Room事件。

Room处理事件并广播结果。

---

## 14. 数据持久化

MySQL保存长期数据。

包括：

用户信息

房间记录

游戏记录

玩家战绩

聊天记录

AI模型版本

Redis保存短期数据。

包括：

在线状态

匹配队列

WebSocket连接映射

房间临时状态缓存

---

## 15. 服务端权威设计

游戏中的所有关键判定由Go后端完成。

包括：

发牌

洗牌

出牌合法性

质疑判定

轮盘惩罚

淘汰判定

胜负判定

客户端不能上传真实牌面。

客户端只能提交自己选择的手牌ID。

服务端根据已保存手牌验证合法性。

---

## 16. 玩家退出与断线设计

游戏未开始（WAITING/MATCHED）时玩家离开：

从房间移除该玩家并广播 PLAYER_LEFT。

若房间变空则销毁房间。

游戏进行中（PLAYING）玩家主动离开或断线：

本局立即结束，状态置为 GAME_OVER。

广播 PLAYER_LEFT（携带 game_over=true、reason、winner_id）。

记录所有玩家战绩（胜者由剩余存活真人玩家判定，若无则无赢家）。

房间在所有真人玩家离开后销毁。

注意：当前实现不再使用 30 秒 AI 托管机制，玩家退出即结束本局，避免剩余玩家被动与 AI 对战。

---

## 17. 日志系统

系统需要记录：

用户登录日志

匹配日志

房间创建日志

游戏行为日志

AI决策日志

异常日志

日志用于调试、统计和训练数据生成。

---

## 18. 部署架构

推荐使用Docker Compose部署。

包含：

frontend

backend

mysql

redis

ai-service

nginx

Nginx负责：

静态资源代理

HTTP反向代理

WebSocket反向代理

HTTPS配置

---

## 19. 可扩展方向

系统后续可以扩展：

骗子骰子模式

观战模式

好友系统

语音聊天

AI模型排行榜

回放系统

恶意Agent检测实验模块

多智能体安全研究接口

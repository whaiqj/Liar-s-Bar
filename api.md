# api.md

# REST API 与 WebSocket 接口设计

## API版本

```text
/api/v1
```

---

# 用户模块

## 注册

```http
POST /api/v1/auth/register
```

Request

```json
{
  "username":"even",
  "password":"123456",
  "nickname":"Even"
}
```

Response

```json
{
  "code":0,
  "msg":"success"
}
```

---

## 登录

```http
POST /api/v1/auth/login
```

Request

```json
{
  "username":"even",
  "password":"123456"
}
```

Response

```json
{
  "code": 0,
  "token": "jwt-token",
  "user": {
    "id": 1,
    "nickname": "Even",
    "username": "even"
  }
}
```

---

## 获取个人信息

```http
GET /api/v1/user/profile
```

---

## 修改资料

```http
PUT /api/v1/user/profile
```

---

# 匹配模块

## 开始匹配

```http
POST /api/v1/match/start
```

Response

```json
{
  "status":"WAITING"
}
```

---

## 取消匹配

```http
POST /api/v1/match/cancel
```

---

## 查询匹配状态

```http
GET /api/v1/match/status
```

Response

```json
{
  "code": 0,
  "status": "MATCHED"
}
```

匹配成功后，服务端通过 WebSocket 推送 `MATCH_FOUND` 消息携带 `room_id`。

---

# 房间模块

## 创建房间

```http
POST /api/v1/rooms
```

---

## 房间列表

```http
GET /api/v1/rooms
```

---

## 房间详情

```http
GET /api/v1/rooms/{id}
```

---

## 加入房间

```http
POST /api/v1/rooms/{id}/join
```

---

## 退出房间

```http
POST /api/v1/rooms/{id}/leave
```

---

# 战绩模块

## 获取历史对局

```http
GET /api/v1/history
```

---

# WebSocket

## 建立连接

```http
GET /ws
```

Header

```text
Authorization: Bearer JWT
```

Browser clients may also pass the token as a query parameter because the WebSocket API cannot set custom headers:

```http
GET /ws?token=JWT
```

---

# WebSocket消息结构

客户端 -> 服务端

```json
{
  "type":"PLAY_CARD",
  "payload":{}
}
```

服务端 -> 客户端

```json
{
  "type":"GAME_STATE",
  "payload":{}
}
```

---

# 房间事件

## PLAYER_JOINED

```json
{
  "type": "PLAYER_JOINED",
  "payload": {
    "player_id": 1,
    "nickname": "Even"
  }
}
```

---

## PLAYER_LEFT

游戏未开始时玩家离开：

```json
{
  "type": "PLAYER_LEFT",
  "payload": {
    "player_id": 3
  }
}
```

游戏进行中玩家离开（直接结束本局）：

```json
{
  "type": "PLAYER_LEFT",
  "payload": {
    "player_id": 3,
    "nickname": "zhangsan",
    "game_over": true,
    "reason": "玩家退出，游戏结束",
    "winner_id": 1
  }
}
```

---

## GAME_STARTED

```json
{
  "type":"GAME_STARTED"
}
```

---

# 游戏事件

## PLAY_CARD

客户端

```json
{
  "type": "PLAY_CARD",
  "payload": {
    "card_ids": [1, 3],
    "claim": "A"
  }
}
```

---

## CHALLENGE

```json
{
  "type": "CHALLENGE",
  "payload": {
    "target_player_id": 2
  }
}
```

---

## PASS

```json
{
  "type":"PASS"
}
```

---

## CHAT

```json
{
  "type":"CHAT",
  "payload":{
      "content":"我保证是真的"
  }
}
```

---

# 服务端广播事件

## GAME_STATE

完整状态同步。

```json
{
  "type": "GAME_STATE",
  "payload": {
    "phase": "PLAYING",
    "current_player": 2,
    "current_round": 1,
    "current_turn": 5,
    "target_card": "A",
    "players": [],
    "alive_count": 4,
    "last_play": {
      "player_id": 1,
      "count": 2,
      "claim": "A"
    }
  }
}
```

---

## CHALLENGE_RESULT

```json
{
  "type": "CHALLENGE_RESULT",
  "payload": {
    "success": true,
    "truthful": false,
    "liar_id": 2,
    "loser_id": 2,
    "challenger_id": 1,
    "challenged_cards": ["K", "K"]
  }
}
```

---

## RUSSIAN_ROULETTE

```json
{
  "type": "RUSSIAN_ROULETTE",
  "payload": {
    "player_id": 2,
    "bullet_count": 3,
    "survived": false
  }
}
```

---

## PLAYER_ELIMINATED

```json
{
  "type": "PLAYER_ELIMINATED",
  "payload": {
    "player_id": 2
  }
}
```

---

## GAME_OVER

```json
{
  "type": "GAME_OVER",
  "payload": {
    "winner_id": 1
  }
}
```

---

## CHAT

```json
{
  "type": "CHAT",
  "payload": {
    "sender_id": 1,
    "sender_name": "Even",
    "content": "我保证是真的",
    "is_ai": false
  }
}
```

---

# AI接口

Go服务器调用AI服务。

## 获取动作

```http
POST /ai/inference
```

Request

```json
{
  "gameState":{},
  "observation":{},
  "legalActions":[]
}
```

Response

```json
{
  "action":"CHALLENGE",
  "confidence":0.92
}
```

---

# 管理员接口

## 查看在线人数

```http
GET /api/v1/admin/online
```

---

## 查看房间

```http
GET /api/v1/admin/rooms
```

# 移动端接口文档

本文档面向安卓和 iOS 客户端。服务端统一提供 HTTP API 和 WebSocket，移动端不实现游戏判定逻辑，只根据服务端返回的状态刷新界面。

## 1. 服务端地址

课程演示阶段可使用：

```text
http://服务器IP:8081
```

正式部署建议使用域名和 HTTPS：

```text
https://你的域名
```

移动端只连接统一入口，不建议直接访问后端容器端口。

```text
HTTP API:  {BASE_URL}/api/v1
WebSocket: ws://服务器IP:8081/ws?token={JWT}
HTTPS 后:  wss://你的域名/ws?token={JWT}
```

网关健康检查：

```http
GET /health
```

响应：

```json
{"status":"ok","service":"liarsbar-gateway"}
```

## 2. 认证方式

登录成功后服务端返回 JWT。移动端需要保存 JWT：

- Android：建议用 Jetpack DataStore 或 EncryptedSharedPreferences。
- iOS：建议用 Keychain。

除注册、登录外，HTTP 请求都带：

```http
Authorization: Bearer <JWT>
```

WebSocket 连接使用 query 参数：

```text
/ws?token=<JWT>
```

## 3. HTTP 通用约定

请求体使用 JSON：

```http
Content-Type: application/json
Accept: application/json
```

常见响应结构：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

`code = 0` 表示成功。`401` 表示 token 缺失或失效，客户端应清除本地登录态并跳回登录页。

## 4. 账号接口

### 注册

```http
POST /api/v1/auth/register
```

请求：

```json
{
  "username": "even",
  "password": "123456",
  "nickname": "Even"
}
```

响应：

```json
{
  "code": 0,
  "msg": "success"
}
```

### 登录

```http
POST /api/v1/auth/login
```

请求：

```json
{
  "username": "even",
  "password": "123456"
}
```

响应：

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

### 获取个人资料

```http
GET /api/v1/user/profile
Authorization: Bearer <JWT>
```

### 修改个人资料

```http
PUT /api/v1/user/profile
Authorization: Bearer <JWT>
```

请求：

```json
{
  "nickname": "NewName",
  "avatar_url": "https://example.com/avatar.png"
}
```

## 5. 大厅与匹配接口

### 获取大厅状态

```http
GET /api/v1/lobby
Authorization: Bearer <JWT>
```

响应示例：

```json
{
  "code": 0,
  "data": {
    "online_count": 1,
    "queue_length": 0,
    "active_rooms": []
  }
}
```

### 开始匹配

```http
POST /api/v1/match/start
Authorization: Bearer <JWT>
```

响应：

```json
{
  "code": 0,
  "status": "WAITING"
}
```

匹配成功后，服务端会通过 WebSocket 推送 `MATCH_FOUND`。

### 取消匹配

```http
POST /api/v1/match/cancel
Authorization: Bearer <JWT>
```

### 查询匹配状态

```http
GET /api/v1/match/status
Authorization: Bearer <JWT>
```

## 6. 房间接口

### 创建房间

```http
POST /api/v1/rooms
Authorization: Bearer <JWT>
```

请求：

```json
{
  "name": "New Room"
}
```

### 获取房间列表

```http
GET /api/v1/rooms
Authorization: Bearer <JWT>
```

### 获取房间详情

```http
GET /api/v1/rooms/{id}
Authorization: Bearer <JWT>
```

### 加入房间

```http
POST /api/v1/rooms/{id}/join
Authorization: Bearer <JWT>
```

加入房间接口成功后，客户端还需要连接 WebSocket 并发送 `PLAYER_JOIN`，服务端才会把玩家加入实时房间状态。

### 离开房间

```http
POST /api/v1/rooms/{id}/leave
Authorization: Bearer <JWT>
```

## 7. WebSocket 协议

连接地址：

```text
ws://服务器IP:8081/ws?token=<JWT>
```

消息结构：

```json
{
  "type": "消息类型",
  "payload": {}
}
```

### 客户端发送

加入实时房间：

```json
{
  "type": "PLAYER_JOIN",
  "payload": {
    "room_id": 1
  }
}
```

准备：

```json
{
  "type": "PLAYER_READY",
  "payload": {}
}
```

出牌：

```json
{
  "type": "PLAY_CARD",
  "payload": {
    "card_ids": [1, 2],
    "claim": "A"
  }
}
```

质疑：

```json
{
  "type": "CHALLENGE",
  "payload": {
    "target_player_id": 2
  }
}
```

跳过：

```json
{
  "type": "PASS",
  "payload": {}
}
```

聊天：

```json
{
  "type": "CHAT",
  "payload": {
    "content": "hello"
  }
}
```

### 服务端推送

常见消息类型：

```text
ROOM_STATE
GAME_STATE
GAME_STARTED
MATCH_FOUND
PLAYER_JOINED
PLAYER_LEFT
PLAYER_READY
CHALLENGE_RESULT
PLAYER_ELIMINATED
RUSSIAN_ROULETTE
CHAT
GAME_OVER
ERROR
```

移动端应以 `GAME_STATE` 和 `ROOM_STATE` 为主要渲染依据，不要在本地自行推演游戏状态。

## 8. 客户端开发建议

### Android

推荐：

- Kotlin
- Jetpack Compose
- Retrofit 或 Ktor Client 处理 HTTP
- OkHttp WebSocket 处理实时消息
- DataStore / EncryptedSharedPreferences 保存 token

### iOS

推荐：

- Swift
- SwiftUI
- URLSession 处理 HTTP
- URLSessionWebSocketTask 处理实时消息
- Keychain 保存 token

## 9. 联调检查

服务端提供冒烟测试脚本：

```bash
python3 scripts/mobile_smoke_test.py
```

指定服务器地址：

```bash
python3 scripts/mobile_smoke_test.py http://服务器IP:8081
```

脚本会依次验证：

1. 注册
2. 登录
3. 获取个人资料
4. 获取大厅状态
5. 创建房间
6. 加入房间
7. WebSocket 升级握手
8. 离开房间

全部通过会输出：

```text
mobile_smoke_test=passed
```

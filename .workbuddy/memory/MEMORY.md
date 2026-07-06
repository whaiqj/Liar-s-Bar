# Liar's Bar Online - 项目记忆

## 项目概述
基于 Go + Vue3 + WebSocket + Python PPO 的多人在线欺骗博弈游戏平台（骗子酒馆）

## 技术栈
- 前端: Vue3 + Vite + Pinia + Vue Router
- 后端: Go + Gin + Gorilla WebSocket + GORM
- AI: Python + FastAPI + PPO (自实现)
- 数据库: MySQL + Redis
- 部署: Docker Compose + Nginx

## 项目结构
- `/backend` - Go 后端服务
- `/frontend` - Vue3 前端
- `/ai-service` - Python AI 服务
- `/deploy` - Docker 部署配置
- `docker-compose.yml` - 编排配置

## 关键设计决策
- 每个游戏房间对应一个独立 Goroutine，串行处理事件
- 服务端权威判定，客户端不可修改游戏状态
- 匹配超时后自动补充 AI 玩家 (10s 允许AI, 15s 自动补齐)
- 断线30秒内重连恢复，超时AI托管
- ELO 算法进行排名
- WebSocket 用于实时游戏通信

## 运行方式
Docker: `docker-compose up -d`
本地开发: 分别启动 backend, frontend, ai-service

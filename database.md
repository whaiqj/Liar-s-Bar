# database.md

# 数据库设计文档

## 数据库名称

```sql
liars_bar
```

---

# users 用户表

保存玩家账号信息。

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    username VARCHAR(50) UNIQUE NOT NULL,

    password_hash VARCHAR(255) NOT NULL,

    nickname VARCHAR(50) NOT NULL,

    avatar_url VARCHAR(255),

    email VARCHAR(100),

    elo_rating INT DEFAULT 1000,

    total_games INT DEFAULT 0,

    total_wins INT DEFAULT 0,

    total_losses INT DEFAULT 0,

    total_lies INT DEFAULT 0,

    total_challenges INT DEFAULT 0,

    total_successful_challenges INT DEFAULT 0,

    status ENUM(
        'ONLINE',
        'OFFLINE',
        'IN_GAME'
    ) DEFAULT 'OFFLINE',

    created_at DATETIME,

    updated_at DATETIME
);
```

---

# rooms 房间表

保存房间基础信息。

```sql
CREATE TABLE rooms (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    room_uuid VARCHAR(64) UNIQUE,

    host_user_id BIGINT,

    room_name VARCHAR(100),

    max_players INT DEFAULT 4,

    current_players INT DEFAULT 0,

    room_status ENUM(
        'WAITING',
        'MATCHED',
        'PLAYING',
        'FINISHED'
    ),

    created_at DATETIME,

    started_at DATETIME,

    finished_at DATETIME
);
```

---

# room_players 房间玩家表

保存玩家与房间关系。

```sql
CREATE TABLE room_players (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    room_id BIGINT,

    user_id BIGINT,

    is_ai BOOLEAN DEFAULT FALSE,

    seat_index INT,

    join_time DATETIME
);
```

---

# games 游戏记录表

每局游戏对应一条记录。

```sql
CREATE TABLE games (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    game_uuid VARCHAR(64) UNIQUE,

    room_id BIGINT,

    winner_user_id BIGINT,

    total_rounds INT,

    total_turns INT,

    ai_count INT,

    start_time DATETIME,

    end_time DATETIME
);
```

---

# game_players 玩家对局表

记录玩家在某一局的表现。

```sql
CREATE TABLE game_players (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    game_id BIGINT,

    user_id BIGINT,

    is_ai BOOLEAN,

    final_rank INT,

    survived BOOLEAN,

    lie_count INT,

    challenge_count INT,

    challenge_success_count INT,

    punishment_count INT,

    bullets_fired INT,

    score_change INT
);
```

---

# game_actions 行为日志表

记录每一步操作。

```sql
CREATE TABLE game_actions (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    game_id BIGINT,

    player_id BIGINT,

    round_no INT,

    turn_no INT,

    action_type VARCHAR(50),

    action_data JSON,

    created_at DATETIME
);
```

action_type:

```text
PLAY_CARD
CHALLENGE
PASS
CHAT
PUNISHMENT
ELIMINATED
GAME_OVER
```

---

# chat_records 聊天记录表

```sql
CREATE TABLE chat_records (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    room_id BIGINT,

    sender_id BIGINT,

    is_ai BOOLEAN,

    content TEXT,

    created_at DATETIME
);
```

---

# matchmaking_queue 匹配队列表

```sql
CREATE TABLE matchmaking_queue (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    user_id BIGINT,

    joined_at DATETIME,

    status ENUM(
        'WAITING',
        'MATCHED',
        'CANCELLED'
    )
);
```

---

# ai_models AI模型表

保存强化学习模型版本。

```sql
CREATE TABLE ai_models (

    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    model_name VARCHAR(100),

    version VARCHAR(50),

    win_rate DOUBLE,

    avg_reward DOUBLE,

    model_path VARCHAR(255),

    deployed BOOLEAN,

    created_at DATETIME
);
```

---

# 推荐索引

```sql
CREATE INDEX idx_user_username
ON users(username);

CREATE INDEX idx_room_status
ON rooms(room_status);

CREATE INDEX idx_game_winner
ON games(winner_user_id);

CREATE INDEX idx_action_game
ON game_actions(game_id);

CREATE INDEX idx_chat_room
ON chat_records(room_id);
```

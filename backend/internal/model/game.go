package model

import "time"

type Game struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	GameUUID     string     `gorm:"uniqueIndex;size:64" json:"game_uuid"`
	RoomID       uint       `json:"room_id"`
	WinnerUserID *uint      `json:"winner_user_id"`
	TotalRounds  int        `json:"total_rounds"`
	TotalTurns   int        `json:"total_turns"`
	AICount      int        `json:"ai_count"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
}

func (Game) TableName() string {
	return "games"
}

type GamePlayer struct {
	ID                    uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	GameID                uint  `json:"game_id"`
	UserID                uint  `json:"user_id"`
	IsAI                  bool  `json:"is_ai"`
	FinalRank             int   `json:"final_rank"`
	Survived              bool  `json:"survived"`
	LieCount              int   `json:"lie_count"`
	ChallengeCount        int   `json:"challenge_count"`
	ChallengeSuccessCount int   `json:"challenge_success_count"`
	PunishmentCount       int   `json:"punishment_count"`
	BulletsFired          int   `json:"bullets_fired"`
	ScoreChange           int   `json:"score_change"`
}

func (GamePlayer) TableName() string {
	return "game_players"
}

type GameAction struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GameID     uint      `json:"game_id"`
	PlayerID   uint      `json:"player_id"`
	RoundNo    int       `json:"round_no"`
	TurnNo     int       `json:"turn_no"`
	ActionType string    `gorm:"size:50" json:"action_type"`
	ActionData string    `gorm:"type:json" json:"action_data"`
	CreatedAt  time.Time `json:"created_at"`
}

func (GameAction) TableName() string {
	return "game_actions"
}

type ChatRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID    uint      `json:"room_id"`
	SenderID  uint      `json:"sender_id"`
	IsAI      bool      `json:"is_ai"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (ChatRecord) TableName() string {
	return "chat_records"
}

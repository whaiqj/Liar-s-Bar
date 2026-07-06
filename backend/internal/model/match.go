package model

import "time"

type MatchStatus string

const (
	MatchStatusWaiting   MatchStatus = "WAITING"
	MatchStatusMatched   MatchStatus = "MATCHED"
	MatchStatusCancelled MatchStatus = "CANCELLED"
)

type MatchmakingQueue struct {
	ID       uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint        `json:"user_id"`
	JoinedAt time.Time   `json:"joined_at"`
	Status   MatchStatus `gorm:"size:20" json:"status"`
}

func (MatchmakingQueue) TableName() string {
	return "matchmaking_queue"
}

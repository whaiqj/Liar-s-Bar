package model

import "time"

type RoomStatus string

const (
	RoomStatusWaiting  RoomStatus = "WAITING"
	RoomStatusMatched  RoomStatus = "MATCHED"
	RoomStatusPlaying  RoomStatus = "PLAYING"
	RoomStatusFinished RoomStatus = "FINISHED"
)

type Room struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomUUID       string     `gorm:"uniqueIndex;size:64" json:"room_uuid"`
	HostUserID     uint       `json:"host_user_id"`
	RoomName       string     `gorm:"size:100" json:"room_name"`
	MaxPlayers     int        `gorm:"default:4" json:"max_players"`
	CurrentPlayers int        `gorm:"default:0" json:"current_players"`
	RoomStatus     RoomStatus `gorm:"size:20" json:"room_status"`
	CreatedAt      time.Time  `json:"created_at"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
}

func (Room) TableName() string {
	return "rooms"
}

type RoomPlayer struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	IsAI      bool      `gorm:"default:false" json:"is_ai"`
	SeatIndex int       `json:"seat_index"`
	JoinTime  time.Time `json:"join_time"`
}

func (RoomPlayer) TableName() string {
	return "room_players"
}

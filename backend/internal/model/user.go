package model

import "time"

type UserStatus string

const (
	UserStatusOnline  UserStatus = "ONLINE"
	UserStatusOffline UserStatus = "OFFLINE"
	UserStatusInGame  UserStatus = "IN_GAME"
)

type User struct {
	ID                      uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username                string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash            string     `gorm:"size:255;not null" json:"-"`
	Nickname                string     `gorm:"size:50;not null" json:"nickname"`
	AvatarURL               string     `gorm:"size:255" json:"avatar_url"`
	Email                   string     `gorm:"size:100" json:"email"`
	EloRating               int        `gorm:"default:1000" json:"elo_rating"`
	TotalGames              int        `gorm:"default:0" json:"total_games"`
	TotalWins               int        `gorm:"default:0" json:"total_wins"`
	TotalLosses             int        `gorm:"default:0" json:"total_losses"`
	TotalLies               int        `gorm:"default:0" json:"total_lies"`
	TotalChallenges         int        `gorm:"default:0" json:"total_challenges"`
	TotalSuccessfulChallenges int       `gorm:"default:0" json:"total_successful_challenges"`
	Status                  UserStatus `gorm:"size:20;default:OFFLINE" json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

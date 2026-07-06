package model

import "time"

type AIModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ModelName  string    `gorm:"size:100" json:"model_name"`
	Version    string    `gorm:"size:50" json:"version"`
	WinRate    float64   `json:"win_rate"`
	AvgReward  float64   `json:"avg_reward"`
	ModelPath  string    `gorm:"size:255" json:"model_path"`
	Deployed   bool      `json:"deployed"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AIModel) TableName() string {
	return "ai_models"
}

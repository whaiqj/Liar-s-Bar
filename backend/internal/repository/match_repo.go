package repository

import (
	"liars-bar/internal/database"
	"liars-bar/internal/model"
)

type MatchRepo struct{}

func NewMatchRepo() *MatchRepo { return &MatchRepo{} }

func (r *MatchRepo) AddToQueue(userID uint) error {
	mq := &model.MatchmakingQueue{
		UserID: userID,
		Status: model.MatchStatusWaiting,
	}
	return database.DB.Create(mq).Error
}

func (r *MatchRepo) RemoveFromQueue(userID uint) error {
	return database.DB.Where("user_id = ? AND status = ?", userID, model.MatchStatusWaiting).
		Update("status", model.MatchStatusCancelled).Error
}

func (r *MatchRepo) GetWaiting() ([]model.MatchmakingQueue, error) {
	var queue []model.MatchmakingQueue
	err := database.DB.Where("status = ?", model.MatchStatusWaiting).
		Order("joined_at ASC").Find(&queue).Error
	return queue, err
}

func (r *MatchRepo) MarkMatched(ids []uint) error {
	return database.DB.Model(&model.MatchmakingQueue{}).
		Where("id IN ?", ids).
		Update("status", model.MatchStatusMatched).Error
}

func (r *MatchRepo) IsInQueue(userID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&model.MatchmakingQueue{}).
		Where("user_id = ? AND status = ?", userID, model.MatchStatusWaiting).
		Count(&count).Error
	return count > 0, err
}

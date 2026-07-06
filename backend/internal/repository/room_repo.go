package repository

import (
	"time"

	"liars-bar/internal/database"
	"liars-bar/internal/model"
)

type RoomRepo struct{}

func NewRoomRepo() *RoomRepo { return &RoomRepo{} }

func (r *RoomRepo) Create(room *model.Room) error {
	return database.DB.Create(room).Error
}

func (r *RoomRepo) FindByID(id uint) (*model.Room, error) {
	var room model.Room
	err := database.DB.First(&room, id).Error
	return &room, err
}

func (r *RoomRepo) FindByUUID(uuid string) (*model.Room, error) {
	var room model.Room
	err := database.DB.Where("room_uuid = ?", uuid).First(&room).Error
	return &room, err
}

func (r *RoomRepo) FindActive() ([]model.Room, error) {
	var rooms []model.Room
	err := database.DB.Where("room_status IN ?", []string{"WAITING", "MATCHED", "PLAYING"}).Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepo) Update(room *model.Room) error {
	return database.DB.Save(room).Error
}

func (r *RoomRepo) Delete(id uint) error {
	if err := database.DB.Where("room_id = ?", id).Delete(&model.RoomPlayer{}).Error; err != nil {
		return err
	}
	return database.DB.Delete(&model.Room{}, id).Error
}

func (r *RoomRepo) UpdateStatus(id uint, status model.RoomStatus) error {
	return database.DB.Model(&model.Room{}).Where("id = ?", id).Update("room_status", status).Error
}

func (r *RoomRepo) AddPlayer(roomID, userID uint, seat int, isAI bool) error {
	rp := &model.RoomPlayer{
		RoomID:    roomID,
		UserID:    userID,
		IsAI:      isAI,
		SeatIndex: seat,
		JoinTime:  time.Now(),
	}
	return database.DB.Create(rp).Error
}

func (r *RoomRepo) RemovePlayer(roomID, userID uint) error {
	return database.DB.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&model.RoomPlayer{}).Error
}

func (r *RoomRepo) GetPlayers(roomID uint) ([]model.RoomPlayer, error) {
	var players []model.RoomPlayer
	err := database.DB.Where("room_id = ?", roomID).Order("seat_index ASC").Find(&players).Error
	return players, err
}

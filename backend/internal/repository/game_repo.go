package repository

import (
	"liars-bar/internal/database"
	"liars-bar/internal/model"
)

type GameRepo struct{}

func NewGameRepo() *GameRepo { return &GameRepo{} }

func (r *GameRepo) Create(game *model.Game) error {
	return database.DB.Create(game).Error
}

func (r *GameRepo) FindByID(id uint) (*model.Game, error) {
	var game model.Game
	err := database.DB.First(&game, id).Error
	return &game, err
}

func (r *GameRepo) FindByUUID(uuid string) (*model.Game, error) {
	var game model.Game
	err := database.DB.Where("game_uuid = ?", uuid).First(&game).Error
	return &game, err
}

func (r *GameRepo) FindByUserID(userID uint, limit int) ([]model.Game, error) {
	var games []model.Game
	err := database.DB.
		Joins("JOIN game_players ON game_players.game_id = games.id").
		Where("game_players.user_id = ?", userID).
		Order("games.start_time DESC").
		Limit(limit).
		Find(&games).Error
	return games, err
}

func (r *GameRepo) Update(game *model.Game) error {
	return database.DB.Save(game).Error
}

func (r *GameRepo) CreateGamePlayer(gp *model.GamePlayer) error {
	return database.DB.Create(gp).Error
}

func (r *GameRepo) GetGamePlayers(gameID uint) ([]model.GamePlayer, error) {
	var players []model.GamePlayer
	err := database.DB.Where("game_id = ?", gameID).Find(&players).Error
	return players, err
}

func (r *GameRepo) CreateAction(action *model.GameAction) error {
	return database.DB.Create(action).Error
}

func (r *GameRepo) GetActions(gameID uint) ([]model.GameAction, error) {
	var actions []model.GameAction
	err := database.DB.Where("game_id = ?", gameID).Order("created_at ASC").Find(&actions).Error
	return actions, err
}

func (r *GameRepo) CreateChat(chat *model.ChatRecord) error {
	return database.DB.Create(chat).Error
}

func (r *GameRepo) GetChats(roomID uint, limit int) ([]model.ChatRecord, error) {
	var chats []model.ChatRecord
	err := database.DB.Where("room_id = ?", roomID).Order("created_at DESC").Limit(limit).Find(&chats).Error
	return chats, err
}

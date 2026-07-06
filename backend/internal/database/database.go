package database

import (
	"liars-bar/internal/config"
	"liars-bar/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.RoomPlayer{},
		&model.Game{},
		&model.GamePlayer{},
		&model.GameAction{},
		&model.ChatRecord{},
		&model.MatchmakingQueue{},
		&model.AIModel{},
	); err != nil {
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	return nil
}

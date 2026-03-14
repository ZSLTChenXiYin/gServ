package model

import (
	"time"

	"gorm.io/datatypes"
)

type Player struct {
	ModelHeader

	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,max=255,email"`
	PasswordHash string `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	Nickname     string `gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	ExpiredAt    time.Time

	ModelTail
}

func (Player) TableName() string { return "players" }

type PlayerDataArchive struct {
	ModelHeader

	GameID   uint
	PlayerID uint

	Data datatypes.JSON

	ModelTail

	Game   Game   `gorm:"foreignKey:GameID"`
	Player Player `gorm:"foreignKey:PlayerID"`
}

func (PlayerDataArchive) TableName() string { return "player_data_archives" }

package model

import (
	"time"

	"gorm.io/gorm"
)

type ModelHeader struct {
	ID uint `gorm:"primary_key"`
}

type ModelTail struct {
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

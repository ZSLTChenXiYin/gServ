package model

type Game struct {
	ModelHeader

	Name string `gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,min=1,max=255"`

	ModelTail
}

func (Game) TableName() string { return "games" }

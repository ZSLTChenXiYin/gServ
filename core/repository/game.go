package repository

import "gServ/pkg/model"

func CreateGame(name string) (*model.Game, error) {
	game := &model.Game{
		Name: name,
	}
	err := database.Create(game).Error
	return game, err
}

func FindGamesWithIndexAndLimit(index int, limit int) ([]model.Game, error) {
	games := []model.Game{}

	// 计算偏移量
	offset := (index - 1) * limit

	err := database.Model(&model.Game{}).
		Offset(offset).
		Limit(limit).
		Find(&games).Error

	return games, err
}

func FindGames() ([]model.Game, error) {
	games := []model.Game{}
	err := database.Model(&model.Game{}).Find(&games).Error
	return games, err
}

func DeleteGame(gameID uint) error {
	return database.Where("id = ?", gameID).Delete(&model.Game{}).Error
}

package gameserv

import (
	"gServ/core/repository"
)

func CreateGame(name string) (uint, error) {
	game, err := repository.CreateGame(name)
	if err != nil {
		return 0, err
	}

	games_name[game.ID] = game.Name

	return game.ID, nil
}

func DeleteGame(game_id uint) error {
	err := repository.DeleteGame(game_id)
	if err != nil {
		return err
	}

	delete(games_name, game_id)

	return nil
}

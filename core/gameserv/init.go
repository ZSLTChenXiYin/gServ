package gameserv

import (
	"gServ/core/repository"
	"gServ/pkg/gserv"
	"sync"
)

var (
	room_id_generator = gserv.NewRoomIDGenerator()
	games_name        = make(map[uint]string)
	online_players    = make(map[uint]*sync.Map)              // map[uint]map[uint]*gserv.Player [game_id][player_id]*Player
	game_rooms        = make(map[uint]map[uint64]*gserv.Room) // [game_id][room_id]*Room
)

func Init() error {
	games, err := repository.FindGames()
	if err != nil {
		return err
	}

	for _, game := range games {
		games_name[game.ID] = game.Name
		online_players[game.ID] = &sync.Map{}
	}

	return nil
}

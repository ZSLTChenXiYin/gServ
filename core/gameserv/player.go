package gameserv

import (
	"gServ/core/repository"
	"gServ/pkg/gserv"
)

func PlayerOnline(game_id uint, player_id uint) error {
	player_any, ok := online_players[game_id].Load(player_id)
	if ok {
		game_player := player_any.(*gserv.Player)
		if game_player != nil {
			return nil
		}
	}

	player, err := repository.FirstPlayer(player_id)
	if err != nil {
		return err
	}

	online_players[game_id].Store(player_id, gserv.NewPlayer(player.Email, player.Nickname))

	return nil
}

func GetOnlinePlayer(game_id uint, player_id uint) *gserv.Player {
	player_any, ok := online_players[game_id].Load(player_id)
	if ok {
		game_player := player_any.(*gserv.Player)
		if game_player != nil {
			return game_player
		}
	}
	return nil
}

func PlayerOffline(game_id uint, player_id uint) {
	online_players[game_id].Delete(player_id)
}

type banPlayer struct {
	GameID uint
	RoomID uint64
	Player *gserv.Player
}

func BanPlayer(player_id uint) error {
	var ban_players []banPlayer
	for game_id, online_player := range online_players {
		online_player_instance_any, ok := online_player.LoadAndDelete(player_id)
		if !ok {
			continue
		}
		online_player_instance := online_player_instance_any.(*gserv.Player)
		ban_players = append(ban_players, banPlayer{
			GameID: game_id,
			RoomID: online_player_instance.GetCurrentRoomID(),
		})
	}

	for _, ban_player := range ban_players {
		LeaveRoom(ban_player.GameID, ban_player.RoomID, player_id)
		PlayerOffline(ban_player.GameID, player_id)
	}

	return repository.DeletePlayer(player_id)
}

// 解封玩家
func UnbanPlayer(player_id uint) error {
	return repository.RestorePlayer(player_id)
}

package gameserv

import (
	"encoding/json"
	"gServ/core/repository"
	"gServ/pkg/gserv"
)

func PlayerOnline(game_id uint, player_id uint) error {
	if online_players[game_id][player_id] != nil {
		return nil
	}

	player_data_archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(game_id, player_id)
	if err != nil {
		return err
	}

	archive_data := map[string]any{}
	err = json.Unmarshal(player_data_archive.Data, &archive_data)
	if err != nil {
		return err
	}

	if online_players[game_id][player_id] != nil { // already online
		return nil
	}

	online_players[game_id][player_id] = &gserv.Player{
		Email:    player_data_archive.Player.Email,
		Nickname: player_data_archive.Player.Nickname,
		Data:     archive_data,
	}

	return nil
}

func GetOnlinePlayer(game_id uint, player_id uint) *gserv.Player {
	return online_players[game_id][player_id]
}

func PlayerOffline(game_id uint, player_id uint) error {
	delete(online_players[game_id], player_id)
	return nil
}

type banPlayer struct {
	GameID uint
	RoomID uint64
	Player *gserv.Player
}

func BanPlayer(player_id uint) error {
	var ban_players []banPlayer
	for game_id, online_player := range online_players {
		for online_player_id, online_player_instance := range online_player {
			if online_player_id == player_id {
				ban_players = append(ban_players, banPlayer{
					GameID: game_id,
					RoomID: online_player_instance.CurrentRoomID,
				})
			}
		}
	}

	for _, ban_player := range ban_players {
		game_rooms[ban_player.GameID][ban_player.RoomID].PlayerLeave(player_id)
		PlayerOffline(ban_player.GameID, player_id)
	}

	return repository.DeletePlayer(player_id)
}

// 解封玩家
func UnbanPlayer(player_id uint) error {
	return repository.RestorePlayer(player_id)
}

package gameserv

import (
	"errors"
	"gServ/pkg/gserv"
	"time"
)

func CreateRoom(name string, game_id uint, homeowner_id uint, max_player uint) (uint64, error) {
	if game_rooms[game_id] == nil {
		return 0, errors.New("game not exists")
	}

	room_id := room_id_generator.Generate()

	if game_rooms[game_id][room_id] != nil {
		return 0, errors.New("room already exists")
	}

	online_player_any, ok := online_players[game_id].Load(homeowner_id)
	if !ok {
		return 0, errors.New("player not exists")
	}
	game_player := online_player_any.(*gserv.Player)
	if game_player.GetCurrentRoomID() != 0 {
		return 0, errors.New("player already in room")
	}

	// 创建房间
	room := gserv.NewRoom(name, game_id, room_id, max_player)

	// 添加房间
	game_rooms[game_id][room_id] = room

	return room_id, nil
}

func GetRoomCount(game_id uint) uint {
	return uint(len(game_rooms[game_id]))
}

func JoinRoom(game_id uint, room_id uint64, player_id uint) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	online_player_any, ok := online_players[game_id].Load(player_id)
	if !ok {
		return errors.New("player not exists")
	}
	game_player := online_player_any.(*gserv.Player)
	if game_player.GetCurrentRoomID() != 0 {
		return errors.New("player already in room")
	}

	// 加入房间
	err := game_rooms[game_id][room_id].PlayerJoin(player_id)
	if err != nil {
		return err
	}

	// 修改玩家当前房间位置
	game_player.SetCurrentRoomID(room_id)

	return nil
}

func LeaveRoom(game_id uint, room_id uint64, player_id uint) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	// 离开房间
	err := game_rooms[game_id][room_id].PlayerLeave(player_id)
	if err != nil {
		return err
	}

	online_player_any, ok := online_players[game_id].Load(player_id)
	if !ok {
		return errors.New("player not exists")
	}
	game_player := online_player_any.(*gserv.Player)

	// 删除玩家当前房间位置
	game_player.SetCurrentRoomID(0)

	// 如果房间人数为0，则删除房间
	if game_rooms[game_id][room_id].GetPlayerCount() == 0 {
		// 删除房间
		delete(game_rooms[game_id], room_id)
	}

	return nil
}

func LockRoom(game_id uint, room_id uint64, player_id uint) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	if game_rooms[game_id][room_id].GetHomeownerID() != player_id {
		return errors.New("player not homeowner")
	}

	// 锁定房间
	game_rooms[game_id][room_id].RoomLock()

	return nil
}

func UnlockRoom(game_id uint, room_id uint64, player_id uint) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	if game_rooms[game_id][room_id].GetHomeownerID() != player_id {
		return errors.New("player not homeowner")
	}

	// 解锁房间
	game_rooms[game_id][room_id].RoomUnlock()

	return nil
}

func ForceDeleteRoom(game_id uint, room_id uint64) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	// 删除房间
	delete(game_rooms[game_id], room_id)

	return nil
}

func DeleteRoom(game_id uint, room_id uint64, player_id uint) error {
	if game_rooms[game_id] == nil {
		return errors.New("game not exists")
	}

	if game_rooms[game_id][room_id] == nil {
		return errors.New("room not exists")
	}

	if game_rooms[game_id][room_id].GetHomeownerID() != player_id {
		return errors.New("player not homeowner")
	}

	// 删除房间
	delete(game_rooms[game_id], room_id)

	return nil
}

func CleanRooms() {
	for game_id, rooms := range game_rooms {
		for room_id, room := range rooms {
			if room.GetPlayerCount() == 0 {
				if time.Since(room.GetUsedAt()) > time.Minute*5 {
					// 关闭房间
					game_rooms[game_id][room_id].RoomClose()
					// 删除房间
					delete(game_rooms[game_id], room_id)
				}
			}
		}
	}
}

func GetRoomPlayers(game_id uint, room_id uint64) []uint {
	if game_rooms[game_id] == nil || game_rooms[game_id][room_id] == nil {
		return []uint{}
	}

	room := game_rooms[game_id][room_id]
	return room.GetPlayerIDs()
}

func GetRooms(game_id uint) []*gserv.Room {
	if game_rooms[game_id] == nil {
		return []*gserv.Room{}
	}

	rooms := make([]*gserv.Room, len(game_rooms[game_id]))
	for index, room := range game_rooms[game_id] {
		rooms[index] = room
	}

	return rooms
}

func GetRoom(game_id uint, room_id uint64) *gserv.Room {
	if game_rooms[game_id] == nil {
		return nil
	}
	return game_rooms[game_id][room_id]
}

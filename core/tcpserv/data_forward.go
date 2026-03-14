package tcpserv

import (
	"errors"
	"gServ/core/gameserv"
	"net"
)

// broadcast 广播数据给房间内所有玩家
func broadcast(game_id uint, room_id uint64, source_player_id uint, data []byte) error {
	// 获取房间内所有玩家
	player_ids := gameserv.GetRoomPlayers(game_id, room_id)
	if len(player_ids) == 0 {
		return errors.New("no players in room")
	}

	// 获取所有玩家连接
	player_conns := make([]net.Conn, len(player_ids))
	for index, player_id := range player_ids {
		player_conns[index] = player_manager.Get(game_id, player_id).Conn
	}

	// 创建广播响应
	resp := NewGSERVProtocolRoomBroadcastDataResponse(
		uint32(game_id),
		room_id,
		uint32(source_player_id),
		data,
	)

	// 发送数据给所有玩家
	var handle_errs error
	for _, conn := range player_conns {
		if conn != nil {
			err := handleProtocol(conn, resp)
			if err != nil {
				handle_errs = errors.Join(handle_errs, err)
			}
		}
	}

	return handle_errs
}

// unicast 单播数据给指定玩家
func unicast(game_id uint, room_id uint64, source_player_id uint, target_player_id uint, data []byte) error {
	// 获取房间内所有玩家
	player_ids := gameserv.GetRoomPlayers(game_id, room_id)
	if len(player_ids) == 0 {
		return errors.New("no players in room")
	}

	// 获取目标玩家
	var player_id uint
	for _, id := range player_ids {
		if id == target_player_id {
			player_id = id
			break
		}
	}
	if player_id == 0 {
		return errors.New("target player not found")
	}

	// 创建单播响应
	response := &GSERVProtocolRoomUnicastDataResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_ROOM_UNICAST_DATA,
		},
		GameID:         uint32(game_id),
		RoomID:         room_id,
		SourcePlayerID: uint32(source_player_id),
		DataLength:     uint32(len(data)),
		Data:           data,
	}

	// 获取玩家连接
	player_conn := player_manager.Get(game_id, player_id).Conn
	if player_conn == nil {
		return errors.New("player not found")
	}

	// 发送数据
	err := handleProtocol(player_conn, response)
	if err != nil {
		return err
	}

	return nil
}

// multicast 组播数据给指定玩家列表
func multicast(game_id uint, room_id uint64, source_player_id uint, target_player_ids []uint, data []byte) error {
	if len(target_player_ids) == 0 {
		return errors.New("no target players specified")
	}

	// 获取房间内所有玩家
	player_ids := gameserv.GetRoomPlayers(game_id, room_id)
	if len(player_ids) == 0 {
		return errors.New("no players in room")
	}

	player_conns := make([]net.Conn, len(target_player_ids))
	for index, player_id := range player_ids {
		for _, target_player_id := range target_player_ids {
			if player_id == target_player_id {
				player_conns[index] = player_manager.Get(game_id, player_id).Conn
				break
			}
		}
	}

	// 创建组播响应
	response := &GSERVProtocolRoomMulticastDataResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_ROOM_MULTICAST_DATA,
		},
		GameID:         uint32(game_id),
		RoomID:         room_id,
		SourcePlayerID: uint32(source_player_id),
		DataLength:     uint32(len(data)),
		Data:           data,
	}

	// 发送给指定玩家组
	var handle_errs error
	for _, conn := range player_conns {
		if conn != nil {
			err := handleProtocol(conn, response)
			if err != nil {
				handle_errs = errors.Join(handle_errs, err)
			}
		}
	}

	return nil
}

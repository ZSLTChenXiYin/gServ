package gserv

import "time"

type AuthPlayer struct {
	ID        uint      // 玩家ID
	ExpiredAt time.Time // 玩家Token过期时间
}

type Player struct {
	Email    string         // 玩家邮箱
	Nickname string         // 玩家昵称
	Data     map[string]any // 玩家在线数据

	CurrentRoomID uint64
}

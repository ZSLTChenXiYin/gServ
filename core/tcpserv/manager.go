package tcpserv

import (
	"net"
	"sync"
	"time"
)

type PlayerConnect struct {
	Lock   sync.RWMutex
	Conn   net.Conn
	UsedAt time.Time
}

func NewPlayerConnect(conn net.Conn) *PlayerConnect {
	return &PlayerConnect{Conn: conn, UsedAt: time.Now()}
}

type PlayerManager struct {
	conns map[uint]map[uint]*PlayerConnect
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{conns: make(map[uint]map[uint]*PlayerConnect)}
}

// 创建玩家连接，直接赋值表示允许挤占连接
func (pm *PlayerManager) Create(game_id uint, player_id uint, conn net.Conn) error {
	if pm.conns[game_id][player_id] == nil {
		err := pm.conns[game_id][player_id].Conn.Close()
		if err != nil {
			return err
		}
	}

	pm.conns[game_id][player_id] = NewPlayerConnect(conn)

	return nil
}

func (pm *PlayerManager) Get(game_id uint, player_id uint) *PlayerConnect {
	pm.conns[game_id][player_id].Lock.RLock()
	defer pm.conns[game_id][player_id].Lock.RUnlock()

	if pm.conns[game_id][player_id] == nil {
		return nil
	}

	pm.conns[game_id][player_id].UsedAt = time.Now()

	return pm.conns[game_id][player_id]
}

func (pm *PlayerManager) Delete(game_id uint, player_id uint) {
	pm.conns[game_id][player_id].Lock.Lock()
	defer pm.conns[game_id][player_id].Lock.Unlock()

	pm.conns[game_id][player_id].Conn.Close()

	delete(pm.conns[game_id], player_id)
}

package gserv

import (
	"errors"
	"math/rand"
	"slices"
	"sync"
	"time"
)

const (
	ROOM_MAX_PLAYER = 4
)

type Room struct {
	name    string // 房间名称
	game_id uint   // 游戏ID
	room_id uint64 // 房间ID

	locked bool

	room_lock    sync.RWMutex
	homeowner_id uint
	max_player   uint   // 房间最大人数
	player_count uint   // 房间当前人数，同时也是当前空余可加入位置
	player_ids   []uint // 房间内玩家ID列表，当玩家退出时，将玩家ID置为0，并将最后一个位置的玩家ID移动到该位置，保证玩家ID列表的连续性
	used_at      time.Time
	created_at   time.Time
}

func NewRoom(name string, game_id uint, room_id uint64, homeowner_id uint, max_player uint) *Room {
	new_room := &Room{
		name:         name,
		game_id:      game_id,
		room_id:      room_id,
		locked:       false,
		homeowner_id: homeowner_id,
		max_player:   max_player,
		player_ids:   make([]uint, max_player),
		used_at:      time.Now(),
		created_at:   time.Now(),
	}
	new_room.player_ids[0] = homeowner_id
	new_room.player_count = 1
	return new_room
}

func (r *Room) GetID() uint64 {
	r.room_lock.RLock()
	defer r.room_lock.RUnlock()

	r.used_at = time.Now()
	return r.room_id
}

func (r *Room) GetName() string {
	r.room_lock.RLock()
	defer r.room_lock.RUnlock()

	r.used_at = time.Now()
	return r.name
}

func (r *Room) GetHomeownerID() uint {
	r.room_lock.RLock()
	defer r.room_lock.RUnlock()

	r.used_at = time.Now()
	return r.homeowner_id
}

func (r *Room) GetMaxPlayer() uint {
	r.room_lock.RLock()

	r.used_at = time.Now()

	return r.max_player
}

func (r *Room) GetPlayerCount() uint {
	r.room_lock.RLock()

	r.used_at = time.Now()

	return r.player_count
}

func (r *Room) GetUsedAt() time.Time {
	r.room_lock.RLock()

	r.used_at = time.Now()

	return r.used_at
}

func (r *Room) GetCreatedAt() time.Time {
	r.room_lock.RLock()

	r.used_at = time.Now()

	return r.created_at
}

// 更新房间使用时间，用于房间在进行广播、单播等操作时刷新房间使用时间，防止房间被删除
func (r *Room) UpdateUsedAt() {
	r.room_lock.Lock()
	defer r.room_lock.Unlock()

	r.used_at = time.Now()
}

func (r *Room) ExistsPlayer(player_id uint) bool {
	r.room_lock.RLock()

	r.used_at = time.Now()

	return slices.Contains(r.player_ids, player_id)
}

func (r *Room) PlayerJoin(player_id uint) error {
	r.room_lock.Lock()
	defer r.room_lock.Unlock()

	r.used_at = time.Now()

	if r.locked {
		return errors.New("room locked")
	}

	if r.player_count >= r.max_player {
		return errors.New("room full")
	}

	if r.ExistsPlayer(player_id) {
		return errors.New("player already in room")
	}

	r.player_ids[r.player_count] = player_id
	r.player_count++

	return nil
}

func (r *Room) PlayerLeave(player_id uint) error {
	r.room_lock.Lock()
	defer r.room_lock.Unlock()

	r.used_at = time.Now()

	for index := 0; index < len(r.player_ids); index++ {
		// 找到玩家ID
		if r.player_ids[index] == player_id {
			// 移动最后一个玩家ID到该位置，并减少房间人数
			if index < len(r.player_ids)-1 {
				if index == len(r.player_ids)-1 {
					r.homeowner_id = 0
					r.player_ids[index] = 0
					r.player_count--
					return nil
				}

				// 如果是房主，则将最后一个玩家ID设为房主
				if r.homeowner_id == player_id {
					r.homeowner_id = r.player_ids[len(r.player_ids)-1]
				}
				// 移动最后一个玩家ID到该位置
				r.player_ids[index] = r.player_ids[len(r.player_ids)-1]
			} else {
				// 如果玩家ID是房主，则将房主ID置为当前位置的前一个玩家ID
				if r.homeowner_id == player_id {
					r.homeowner_id = r.player_ids[index-1]
				}
				// 如果玩家ID是最后一个，则将玩家ID置为0
				r.player_ids[index] = 0
			}

			// 房间人数减1
			r.player_count--

			return nil
		}
	}

	return errors.New("player not in room")
}

func (r *Room) RoomLock() {
	r.room_lock.Lock()
	defer r.room_lock.Unlock()

	r.locked = true

	r.used_at = time.Now()
}

func (r *Room) RoomUnlock() {
	r.room_lock.Lock()
	defer r.room_lock.Unlock()

	r.locked = false

	r.used_at = time.Now()
}

// GetPlayerIDs 获取房间内所有玩家ID
func (r *Room) GetPlayerIDs() []uint {
	r.room_lock.RLock()
	defer r.room_lock.RUnlock()

	r.used_at = time.Now()

	// 返回非零的玩家ID
	result := make([]uint, 0, r.player_count)
	for _, playerID := range r.player_ids {
		if playerID != 0 {
			result = append(result, playerID)
		}
	}
	return result
}

type RoomIDGenerator struct {
	rng       *rand.Rand
	rng_lock  sync.Mutex
	last_time int64
	counter   uint16
}

func NewRoomIDGenerator() *RoomIDGenerator {
	return &RoomIDGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate 生成唯一房间ID
// 格式: 时间戳(秒) + 计数器 + 随机数
func (g *RoomIDGenerator) Generate() uint64 {
	g.rng_lock.Lock()
	defer g.rng_lock.Unlock()

	now := time.Now().Unix()

	// 如果在同一秒内，增加计数器
	if now == g.last_time {
		g.counter++
	} else {
		g.counter = 0
		g.last_time = now
	}

	// 组合ID: 时间戳(高32位) + 计数器(16位) + 随机数(16位)
	// 这样可以保证在分布式系统中也基本唯一
	timestamp := uint64(now) << 32
	counter := uint64(g.counter) << 16
	random := uint64(g.rng.Intn(65536)) // 16位随机数

	return timestamp | counter | random
}

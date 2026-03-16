package gserv

import (
	"sync"
	"time"
)

type AuthPlayer struct {
	ID        uint      // 玩家ID
	ExpiredAt time.Time // 玩家Token过期时间
}

type Player struct {
	lock sync.RWMutex

	email    string // 玩家邮箱
	nickname string // 玩家昵称

	current_room_id uint64
}

func NewPlayer(email string, nickname string) *Player {
	return &Player{
		email:    email,
		nickname: nickname,
	}
}

func (p *Player) SetNickname(nickname string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.nickname = nickname
}

func (p *Player) GetNickname() string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.nickname
}

func (p *Player) SetEmail(email string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.email = email
}

func (p *Player) GetEmail() string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.email
}

func (p *Player) SetCurrentRoomID(room_id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.current_room_id = room_id
}

func (p *Player) GetCurrentRoomID() uint64 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.current_room_id
}

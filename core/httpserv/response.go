package httpserv

import (
	"time"

	"gorm.io/datatypes"
)

type post_Api_Game_Response struct {
	GameID uint `json:"game_id"`
}

type get_Api_Games_Response struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	RoomCount uint      `json:"room_count"`
	CreatedAt time.Time `json:"created_at"`
}

type post_Api_Player_Register_Response struct {
	PlayerID uint `json:"player_id"`
}

type post_Api_Player_Login_Response struct {
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	TCPPort  uint   `json:"tcp_port"`
}

type get_Api_Player_Response struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

type post_Api_Room_Response struct {
	RoomID uint64 `json:"room_id"`
}

type get_Api_Room_Response struct {
	RoomID      uint64    `json:"room_id"`
	Name        string    `json:"name"`
	HomeownerID uint      `json:"homeowner_id"`
	MaxPlayer   uint      `json:"max_player"`
	PlayerCount uint      `json:"player_count"`
	PlayerIDs   []uint    `json:"player_ids"`
	CreatedAt   time.Time `json:"created_at"`
}

type get_Api_Rooms_Response struct {
	RoomID      uint64    `json:"room_id"`
	Name        string    `json:"name"`
	HomeownerID uint      `json:"homeowner_id"`
	MaxPlayer   uint      `json:"max_player"`
	PlayerCount uint      `json:"player_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// JSON数据操作响应结构
type get_Api_Data_Exists_Response struct {
	Exists bool `json:"exists"`
}

type post_Api_Data_Response struct {
	ID uint `json:"id"`
}

type get_Api_Data_Response struct {
	ID        uint           `json:"id"`
	GameID    uint           `json:"game_id"`
	PlayerID  uint           `json:"player_id"`
	Data      datatypes.JSON `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

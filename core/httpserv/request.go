package httpserv

import (
	"gServ/pkg/model"

	"gorm.io/datatypes"
)

type post_Api_Game_Request struct {
	Name string `json:"name"`
}

type get_Api_Games_Request struct {
	Index int `json:"index"`
	Limit int `json:"limit"`
}

type post_Api_Captcha_Email_Request struct {
	Email     string            `json:"email" validate:"required,email"`
	EmailType model.CaptchaType `json:"email_type" validate:"captcha_type"`
}

type post_Api_Player_Register_Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"required"`
	Captcha  string `json:"captcha" validate:"required"`
}

type post_Api_Player_Login_Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type put_Api_Player_Request struct {
	Nickname string `json:"nickname"`
}

type put_Api_Player_Password_Request struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type post_Api_Room_Request struct {
	GameID    uint   `json:"game_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	MaxPlayer uint   `json:"max_player" validate:"required,lte=64"`
}

// JSON数据操作请求结构
type post_Api_Data_Request struct {
	GameID   uint           `json:"game_id" binding:"required"`
	PlayerID uint           `json:"player_id" binding:"required"`
	Data     datatypes.JSON `json:"data" binding:"required"`
}

type put_Api_Data_Request struct {
	Data datatypes.JSON `json:"data" binding:"required"`
}

type get_Api_Data_Request struct {
	GameID   uint `form:"game_id" binding:"required"`
	PlayerID uint `form:"player_id" binding:"required"`
}

type delete_Api_Data_Request struct {
	GameID   uint `form:"game_id" binding:"required"`
	PlayerID uint `form:"player_id" binding:"required"`
}

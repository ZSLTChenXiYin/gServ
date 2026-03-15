package httpserv

import (
	"gServ/core/config"
	"gServ/core/repository"
	"gServ/core/validate"
	"gServ/pkg/hash"
	"gServ/pkg/jwt"
	"gServ/pkg/middleware"
	"gServ/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 玩家注册
func post_Api_Player_Register(c *gin.Context) {
	req := &post_Api_Player_Register_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	if err := validate.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求验证失败"})
		return
	}

	// 验证验证码
	captchas, err := repository.FindEmailCaptchasByEmailAndCaptchaType(req.Email, model.CAPTCHA_TYPE_REGISTER)
	if err != nil || len(captchas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	validated := false
	for _, captcha := range captchas {
		if captcha.Code == req.Captcha {
			if !captcha.UsedAt.IsZero() {
				c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已使用"})
				return
			}
			if err := repository.UpdateEmailCaptchaUsedAt(captcha.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新验证码使用状态失败"})
				return
			}
			validated = true
			break
		}
	}
	if !validated {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 生成密码哈希
	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希生成失败"})
		return
	}

	player, err := repository.CreatePlayer(req.Email, passwordHash, req.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建玩家数据失败"})
		return
	}

	c.JSON(http.StatusOK, &post_Api_Player_Register_Response{
		PlayerID: player.ID,
	})
}

// 玩家登录
func post_Api_Player_Login(c *gin.Context) {
	req := &post_Api_Player_Login_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	if err := validate.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求验证失败"})
		return
	}

	player, err := repository.FirstPlayerByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不存在"})
		return
	}

	if err := hash.ComparePassword(player.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, &post_Api_Player_Login_Response{
		Token:    token,
		Nickname: player.Nickname,
		TCPPort:  config.GetConfig().Server.TCPPort,
	})
}

// 获取玩家信息
func get_Api_Player(c *gin.Context) {
	player_id := c.Param("player_id")
	if player_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家ID不能为空"})
		return
	}

	var player_id_uint uint64
	player_id_uint, err := strconv.ParseUint(player_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家ID格式错误"})
		return
	}

	player, err := repository.FirstPlayer(uint(player_id_uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不存在"})
		return
	}

	c.JSON(http.StatusOK, &get_Api_Player_Response{
		ID:        player.ID,
		Email:     player.Email,
		Nickname:  player.Nickname,
		CreatedAt: player.CreatedAt,
	})
}

// 更新玩家信息
func put_Api_Player(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &put_Api_Player_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	err := repository.UpdatePlayer(auth_player.ID, req.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新玩家数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "玩家信息更新成功"})
}

// 更新玩家密码
func put_Api_Player_Password(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &put_Api_Player_Password_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	player, err := repository.FirstPlayer(auth_player.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不存在"})
		return
	}

	if err := hash.ComparePassword(player.PasswordHash, req.OldPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
		return
	}

	password_hash, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希生成失败"})
		return
	}

	err = repository.UpdatePlayerPasswordHash(auth_player.ID, password_hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

func delete_Api_Player(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	player_id := c.Param("id")
	player_id_uint, err := strconv.ParseUint(player_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家ID格式错误"})
		return
	}

	if auth_player.ID != uint(player_id_uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	err = repository.DeletePlayer(auth_player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除玩家数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player deleted successfully"})
}

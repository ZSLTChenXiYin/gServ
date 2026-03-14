package httpserv

import (
	"gServ/core/gameserv"
	"gServ/core/validate"
	"gServ/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建房间
func post_Api_Room(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &post_Api_Room_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析错误"})
		return
	}

	if err := validate.Validate(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求验证错误"})
		return
	}

	// 默认最大玩家数为8
	if req.MaxPlayer == 0 {
		req.MaxPlayer = 8
	}

	room_id, err := gameserv.CreateRoom(req.Name, req.GameID, auth_player.ID, req.MaxPlayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &post_Api_Room_Response{
		RoomID: room_id,
	})
}

// 获取房间信息
func get_Api_Room(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	game_id := c.Param("game_id")
	room_id := c.Param("room_id")
	if game_id == "" || room_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint, room_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}
	room_id_uint, err = strconv.ParseUint(room_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间ID格式错误"})
		return
	}

	room := gameserv.GetRoom(uint(game_id_uint), room_id_uint)
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	c.JSON(http.StatusOK, &get_Api_Room_Response{
		RoomID:      room.GetID(),
		Name:        room.GetName(),
		HomeownerID: room.GetHomeownerID(),
		MaxPlayer:   room.GetMaxPlayer(),
		PlayerCount: room.GetPlayerCount(),
		PlayerIDs:   room.GetPlayerIDs(),
		CreatedAt:   room.GetCreatedAt(),
	})
}

// 锁定房间
func put_Api_Room_Lock(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	game_id := c.Param("game_id")
	room_id := c.Param("room_id")
	if game_id == "" || room_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint, room_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}
	room_id_uint, err = strconv.ParseUint(room_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间ID格式错误"})
		return
	}

	err = gameserv.LockRoom(uint(game_id_uint), room_id_uint, auth_player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "锁定房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "房间锁定成功"})
}

func put_Api_Room_Unlock(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	game_id := c.Param("game_id")
	room_id := c.Param("room_id")
	if game_id == "" || room_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint, room_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}
	room_id_uint, err = strconv.ParseUint(room_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间ID格式错误"})
		return
	}

	err = gameserv.UnlockRoom(uint(game_id_uint), room_id_uint, auth_player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解锁房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "房间解锁成功"})
}

// 删除房间
func delete_Api_Room(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	game_id := c.Param("game_id")
	room_id := c.Param("room_id")
	if game_id == "" || room_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint, room_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}
	room_id_uint, err = strconv.ParseUint(room_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "房间ID格式错误"})
		return
	}

	err = gameserv.DeleteRoom(uint(game_id_uint), room_id_uint, auth_player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "房间删除成功"})
}

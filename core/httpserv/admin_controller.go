package httpserv

import (
	"net/http"
	"strconv"

	"gServ/core/gameserv"
	"gServ/core/log"
	"gServ/core/repository"

	"github.com/gin-gonic/gin"
)

// 根据index和limit获取游戏列表分页
func get_Admin_Games(c *gin.Context) {
	req := &get_Admin_Games_Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.StdErrorf("HTTP服务绑定获取游戏列表请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, err := repository.FindGamesWithIndexAndLimit(req.Index, req.Limit)
	if err != nil {
		log.StdErrorf("HTTP服务获取游戏列表请求失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取游戏列表失败"})
		return
	}

	resp := make([]get_Api_Games_Response, len(games))
	for index, game := range games {
		resp[index] = get_Api_Games_Response{
			ID:        game.ID,
			Name:      game.Name,
			RoomCount: gameserv.GetRoomCount(game.ID),
			CreatedAt: game.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// 对比config.GetConfig().Server.AuthCode和post_Api_Game_Request中的auth_code
func post_Admin_Game(c *gin.Context) {
	req := &post_Admin_Game_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := gameserv.CreateGame(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建游戏数据失败"})
		return
	}

	resp := &post_Api_Game_Response{
		GameID: id,
	}

	c.JSON(http.StatusOK, resp)
}

func delete_Admin_Game(c *gin.Context) {
	game_id := c.Param("game_id")
	game_id_uint, err := strconv.ParseUint(game_id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效游戏ID"})
		return
	}

	err = gameserv.DeleteGame(uint(game_id_uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除游戏数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "游戏删除成功"})
}

func delete_Admin_Room(c *gin.Context) {
	game_id := c.Param("game_id")
	room_id := c.Param("room_id")

	game_id_uint, err := strconv.ParseUint(game_id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效游戏ID"})
		return
	}

	room_id_uint, err := strconv.ParseUint(room_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效房间ID"})
		return
	}

	err = gameserv.ForceDeleteRoom(uint(game_id_uint), room_id_uint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "房间删除成功"})
}

func delete_Admin_Ban_Player(c *gin.Context) {
	player_id := c.Param("player_id")
	player_id_uint, err := strconv.ParseUint(player_id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效玩家ID"})
		return
	}

	err = gameserv.BanPlayer(uint(player_id_uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "封禁玩家失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "玩家封禁成功"})
}

func put_Admin_Ban_Player(c *gin.Context) {
	playerID := c.Param("player_id")
	_, err := strconv.ParseUint(playerID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效玩家ID"})
		return
	}

	// 这里应该实现解封玩家的逻辑
	// 例如：更新玩家的状态为正常，清除封禁时间等
	// 暂时返回成功
	c.JSON(http.StatusOK, gin.H{"message": "玩家解封成功"})
}

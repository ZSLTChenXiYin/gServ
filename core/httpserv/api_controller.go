package httpserv

import (
	"gServ/core/gameserv"
	"gServ/core/log"
	"gServ/core/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func get_Api_Health(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// 根据index和limit获取游戏列表分页
func get_Api_Games(c *gin.Context) {
	req := &get_Api_Games_Request{}
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

// 获取房间列表
func get_Api_Rooms(c *gin.Context) {
	game_id := c.Param("game_id")
	if game_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID不能为空"})
		return
	}

	var game_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}

	rooms := gameserv.GetRooms(uint(game_id_uint))

	resp := make([]get_Api_Rooms_Response, len(rooms))
	for index, room := range rooms {
		resp[index].RoomID = room.GetID()
		resp[index].Name = room.GetName()
		resp[index].HomeownerID = room.GetHomeownerID()
		resp[index].MaxPlayer = room.GetMaxPlayer()
		resp[index].PlayerCount = room.GetPlayerCount()
		resp[index].CreatedAt = room.GetCreatedAt()
	}

	c.JSON(http.StatusOK, resp)
}

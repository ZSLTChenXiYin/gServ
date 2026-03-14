package httpserv

import (
	"fmt"
	"gServ/core/log"
	"gServ/core/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建或更新JSON数据
func post_Api_Data(c *gin.Context) {
	req := &post_Api_Data_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.StdErrorf("HTTP服务绑定创建数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查游戏和玩家是否存在
	if exists, err := repository.ExistsData(repository.TABLENAME_GAME, req.GameID); err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏不存在"})
		return
	}
	if exists, err := repository.ExistsData(repository.TABLENAME_PLAYER, req.PlayerID); err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家不存在"})
		return
	}

	id, err := repository.CreateOrUpdatePlayerData(req.GameID, req.PlayerID, req.Data)
	if err != nil {
		log.StdErrorf("HTTP服务创建或更新玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建或更新数据失败"})
		return
	}

	resp := &post_Api_Data_Response{
		ID: id,
	}

	c.JSON(http.StatusOK, resp)
}

// 获取JSON数据
func get_Api_Data(c *gin.Context) {
	req := &get_Api_Data_Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.StdErrorf("HTTP服务绑定获取数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := repository.GetPlayerData(req.GameID, req.PlayerID)
	if err != nil {
		log.StdErrorf("HTTP服务获取玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据不存在"})
		return
	}

	resp := &get_Api_Data_Response{
		ID:        data.ID,
		GameID:    data.GameID,
		PlayerID:  data.PlayerID,
		Data:      data.Data,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, resp)
}

// 更新JSON数据
func put_Api_Data(c *gin.Context) {
	req := &put_Api_Data_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.StdErrorf("HTTP服务绑定更新数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从URL参数获取game_id和player_id
	gameID := c.Param("game_id")
	playerID := c.Param("player_id")

	var gid, pid uint
	if _, err := fmt.Sscanf(gameID, "%d", &gid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的游戏ID"})
		return
	}
	if _, err := fmt.Sscanf(playerID, "%d", &pid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的玩家ID"})
		return
	}

	err := repository.UpdatePlayerData(gid, pid, req.Data)
	if err != nil {
		log.StdErrorf("HTTP服务更新玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据更新成功"})
}

// 删除JSON数据
func delete_Api_Data(c *gin.Context) {
	req := &delete_Api_Data_Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.StdErrorf("HTTP服务绑定删除数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.DeletePlayerData(req.GameID, req.PlayerID)
	if err != nil {
		log.StdErrorf("HTTP服务删除玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据删除成功"})
}

// 批量获取游戏下的所有玩家数据
func get_Api_Game_Data(c *gin.Context) {
	gameID := c.Param("game_id")
	var gid uint
	if _, err := fmt.Sscanf(gameID, "%d", &gid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的游戏ID"})
		return
	}

	dataList, err := repository.GetAllPlayerDataByGame(gid)
	if err != nil {
		log.StdErrorf("HTTP服务获取游戏所有玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	resp := make([]get_Api_Data_Response, len(dataList))
	for i, data := range dataList {
		resp[i] = get_Api_Data_Response{
			ID:        data.ID,
			GameID:    data.GameID,
			PlayerID:  data.PlayerID,
			Data:      data.Data,
			CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, resp)
}

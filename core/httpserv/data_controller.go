package httpserv

import (
	"gServ/core/log"
	"gServ/core/repository"
	"gServ/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func get_Api_Data_Exists(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	game_id := c.Param("game_id")
	if game_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}

	exists, err := repository.ExistsPlayerDataArchive(uint(game_id_uint), auth_player.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查玩家数据失败"})
		return
	}

	c.JSON(http.StatusOK, &get_Api_Data_Exists_Response{
		Exists: exists,
	})
}

// 创建或更新JSON数据
func post_Api_Data(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &post_Api_Data_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.StdErrorf("HTTP服务绑定创建数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	// 检查游戏和玩家是否存在
	if exists, err := repository.ExistsData(repository.TABLENAME_GAME, req.GameID); err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏不存在"})
		return
	}
	if exists, err := repository.ExistsData(repository.TABLENAME_PLAYER, auth_player.ID); err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "玩家不存在"})
		return
	}

	id, err := repository.CreatePlayerDataArchive(req.GameID, auth_player.ID, req.Data)
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
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &get_Api_Data_Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		log.StdErrorf("HTTP服务绑定获取数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(req.GameID, auth_player.ID)
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
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// 更新JSON数据
func put_Api_Data(c *gin.Context) {
	auth_player := middleware.GetAuthPlayerFromGinContext(c)
	if auth_player.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}

	req := &put_Api_Data_Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.StdErrorf("HTTP服务绑定更新数据请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game_id := c.Param("game_id")
	if game_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID和房间ID不能为空"})
		return
	}

	var game_id_uint uint64
	game_id_uint, err := strconv.ParseUint(game_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "游戏ID格式错误"})
		return
	}

	err = repository.UpdatePlayerDataArchive(uint(game_id_uint), auth_player.ID, req.Data)
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

	err := repository.DeletePlayerDataArchive(req.GameID, req.PlayerID)
	if err != nil {
		log.StdErrorf("HTTP服务删除玩家数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据删除成功"})
}

package repository

import (
	"gServ/pkg/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateOrUpdatePlayerData 创建或更新玩家数据
func CreateOrUpdatePlayerData(gameID, playerID uint, data datatypes.JSON) (uint, error) {
	var archive model.PlayerDataArchive
	
	// 先尝试查找是否已存在
	err := database.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&archive).Error
	
	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		archive = model.PlayerDataArchive{
			GameID:   gameID,
			PlayerID: playerID,
			Data:     data,
		}
		err = database.Create(&archive).Error
		if err != nil {
			return 0, err
		}
		return archive.ID, nil
	} else if err != nil {
		// 其他错误
		return 0, err
	}
	
	// 已存在，更新数据
	archive.Data = data
	err = database.Save(&archive).Error
	if err != nil {
		return 0, err
	}
	return archive.ID, nil
}

// GetPlayerData 获取玩家数据
func GetPlayerData(gameID, playerID uint) (*model.PlayerDataArchive, error) {
	var archive model.PlayerDataArchive
	err := database.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&archive).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &archive, nil
}

// UpdatePlayerData 更新玩家数据
func UpdatePlayerData(gameID, playerID uint, data datatypes.JSON) error {
	var archive model.PlayerDataArchive
	err := database.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&archive).Error
	if err == gorm.ErrRecordNotFound {
		return gorm.ErrRecordNotFound
	}
	if err != nil {
		return err
	}
	
	archive.Data = data
	return database.Save(&archive).Error
}

// DeletePlayerData 删除玩家数据
func DeletePlayerData(gameID, playerID uint) error {
	return database.Where("game_id = ? AND player_id = ?", gameID, playerID).Delete(&model.PlayerDataArchive{}).Error
}

// GetAllPlayerDataByGame 获取游戏下的所有玩家数据
func GetAllPlayerDataByGame(gameID uint) ([]model.PlayerDataArchive, error) {
	var archives []model.PlayerDataArchive
	err := database.Where("game_id = ?", gameID).Find(&archives).Error
	if err != nil {
		return nil, err
	}
	return archives, nil
}

// GetPlayerDataByPlayer 获取玩家的所有游戏数据
func GetPlayerDataByPlayer(playerID uint) ([]model.PlayerDataArchive, error) {
	var archives []model.PlayerDataArchive
	err := database.Where("player_id = ?", playerID).Find(&archives).Error
	if err != nil {
		return nil, err
	}
	return archives, nil
}
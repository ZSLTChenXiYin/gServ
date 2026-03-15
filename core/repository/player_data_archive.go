package repository

import (
	"errors"
	"gServ/pkg/model"

	"gorm.io/datatypes"
)

// CreatePlayerDataArchive 创建玩家数据
func CreatePlayerDataArchive(game_id uint, player_id uint, data datatypes.JSON) (uint, error) {
	var archive model.PlayerDataArchive

	// 先尝试查找是否已存在
	exists, err := ExistsPlayerDataArchive(game_id, player_id)
	if err != nil {
		return 0, err
	}
	if exists {
		// 已存在
		return 0, errors.New("玩家数据已存在")
	}

	// 不存在，创建新记录
	archive = model.PlayerDataArchive{
		GameID:   game_id,
		PlayerID: player_id,
		Data:     data,
	}
	err = database.Create(&archive).Error
	if err != nil {
		return 0, err
	}

	return archive.ID, nil
}

// FirstPlayerDataArchiveByGameIDAndPlayerID 获取玩家数据
func FirstPlayerDataArchiveByGameIDAndPlayerID(game_id uint, player_id uint) (*model.PlayerDataArchive, error) {
	archive := &model.PlayerDataArchive{}
	err := database.Where("game_id = ? AND player_id = ?", game_id, player_id).First(archive).Error
	if err != nil {
		return nil, err
	}
	return archive, nil
}

// UpdatePlayerDataArchive 更新玩家数据
func UpdatePlayerDataArchive(game_id uint, player_id uint, data datatypes.JSON) error {
	archive := &model.PlayerDataArchive{}
	err := database.Where("game_id = ? AND player_id = ?", game_id, player_id).First(archive).Error
	if err != nil {
		return err
	}

	archive.Data = data
	return database.Save(archive).Error
}

// DeletePlayerDataArchive 删除玩家数据
func DeletePlayerDataArchive(game_id uint, player_id uint) error {
	return database.Where("game_id = ? AND player_id = ?", game_id, player_id).Delete(&model.PlayerDataArchive{}).Error
}

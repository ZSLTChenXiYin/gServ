package repository

import (
	"fmt"
	"gServ/pkg/model"
)

const (
	TABLENAME_GAME                = "games"
	TABLENAME_PLAYER              = "players"
	TABLENAME_PLAYER_DATA_ARCHIVE = "player_data_archives"
)

func ExistsTable(tablename string) bool {
	switch tablename {
	case TABLENAME_GAME, TABLENAME_PLAYER, TABLENAME_PLAYER_DATA_ARCHIVE:
		return true
	default:
		return false
	}
}

func ExistsData(tablename string, id uint) (bool, error) {
	if !ExistsTable(tablename) {
		return false, fmt.Errorf("table %s does not exist", tablename)
	}

	var count int64
	err := database.Table(tablename).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func ExistsPlayerDataArchive(game_id uint, player_id uint) (bool, error) {
	var count int64
	err := database.Model(&model.PlayerDataArchive{}).Where("game_id = ? AND player_id = ?", game_id, player_id).Count(&count).Error
	return count > 0, err
}

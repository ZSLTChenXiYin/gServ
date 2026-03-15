package repository

import (
	"gServ/pkg/model"
	"time"
)

func FirstPlayer(id uint) (*model.Player, error) {
	player := &model.Player{
		ModelHeader: model.ModelHeader{ID: id},
	}
	err := database.First(player).Error
	return player, err
}

func ScanPlayerExpiredAt(id uint) (*time.Time, error) {
	player_expired_at := &time.Time{}
	err := database.Model(&model.Player{}).Where("id = ?", id).Select("expired_at").Scan(player_expired_at).Error
	return player_expired_at, err
}

func CreatePlayer(email, passwordHash, nickname string) (*model.Player, error) {
	player := &model.Player{
		Email:        email,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		ExpiredAt:    time.Now().AddDate(1, 0, 0), // 默认1年有效期
	}
	err := database.Create(player).Error
	return player, err
}

// FirstPlayerByEmail 验证玩家凭据
func FirstPlayerByEmail(email string) (*model.Player, error) {
	player := &model.Player{}
	err := database.Where("email = ?", email).First(player).Error
	return player, err
}

// UpdatePlayer 更新玩家信息
func UpdatePlayer(playerID uint, nickname string) error {
	return database.Model(&model.Player{}).
		Where("id = ?", playerID).
		Update("nickname", nickname).Error
}

// UpdatePlayerPasswordHash 更新玩家密码
func UpdatePlayerPasswordHash(playerID uint, password_hash string) error {
	return database.Model(&model.Player{}).
		Where("id = ?", playerID).
		Update("password_hash", password_hash).Error
}

func DeletePlayer(player_id uint) error {
	return database.Delete(&model.Player{
		ModelHeader: model.ModelHeader{ID: player_id},
	}).Error
}

func RestorePlayer(player_id uint) error {
	return database.Unscoped().Model(&model.Player{}).
		Where("id = ?", player_id).
		Update("deleted_at", nil).Error
}

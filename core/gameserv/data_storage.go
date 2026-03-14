package gameserv

import (
	"encoding/json"
	"errors"
	"gServ/core/repository"
	"gServ/pkg/model"
	"time"
)

// DataStorage 数据存储服务
type DataStorage struct {
	gameID   uint
	playerID uint
}

// NewDataStorage 创建数据存储服务实例
func NewDataStorage(gameID, playerID uint) *DataStorage {
	return &DataStorage{
		gameID:   gameID,
		playerID: playerID,
	}
}

// SaveData 保存玩家数据
func (ds *DataStorage) SaveData(key string, value interface{}) error {
	// 获取现有数据
	archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(ds.gameID, ds.playerID)
	if err != nil {
		// 如果不存在，创建新的数据存档
		data := map[string]interface{}{key: value}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		archive = &model.PlayerDataArchive{
			GameID:   ds.gameID,
			PlayerID: ds.playerID,
			Data:     jsonData,
		}
		return repository.CreatePlayerDataArchive(archive)
	}

	// 解析现有数据
	var data map[string]interface{}
	if err := json.Unmarshal(archive.Data, &data); err != nil {
		return err
	}

	// 更新数据
	data[key] = value

	// 重新序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	archive.Data = jsonData
	return repository.UpdatePlayerDataArchive(archive)
}

// LoadData 加载玩家数据
func (ds *DataStorage) LoadData(key string) (interface{}, error) {
	archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(ds.gameID, ds.playerID)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(archive.Data, &data); err != nil {
		return nil, err
	}

	value, exists := data[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// DeleteData 删除玩家数据
func (ds *DataStorage) DeleteData(key string) error {
	archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(ds.gameID, ds.playerID)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(archive.Data, &data); err != nil {
		return err
	}

	delete(data, key)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	archive.Data = jsonData
	return repository.UpdatePlayerDataArchive(archive)
}

// GetAllData 获取所有玩家数据
func (ds *DataStorage) GetAllData() (map[string]interface{}, error) {
	archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(ds.gameID, ds.playerID)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(archive.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// SaveBatchData 批量保存数据
func (ds *DataStorage) SaveBatchData(data map[string]interface{}) error {
	archive, err := repository.FirstPlayerDataArchiveByGameIDAndPlayerID(ds.gameID, ds.playerID)
	if err != nil {
		// 如果不存在，创建新的数据存档
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		archive = &model.PlayerDataArchive{
			GameID:   ds.gameID,
			PlayerID: ds.playerID,
			Data:     jsonData,
		}
		return repository.CreatePlayerDataArchive(archive)
	}

	// 解析现有数据
	var existingData map[string]interface{}
	if err := json.Unmarshal(archive.Data, &existingData); err != nil {
		return err
	}

	// 合并数据
	for key, value := range data {
		existingData[key] = value
	}

	// 重新序列化
	jsonData, err := json.Marshal(existingData)
	if err != nil {
		return err
	}

	archive.Data = jsonData
	return repository.UpdatePlayerDataArchive(archive)
}

// DataStorageManager 数据存储管理器
type DataStorageManager struct {
	storages map[uint]map[uint]*DataStorage // gameID -> playerID -> DataStorage
}

var dataStorageManager *DataStorageManager

func init() {
	dataStorageManager = &DataStorageManager{
		storages: make(map[uint]map[uint]*DataStorage),
	}
}

// GetDataStorage 获取数据存储服务
func GetDataStorage(gameID, playerID uint) *DataStorage {
	if dataStorageManager.storages[gameID] == nil {
		dataStorageManager.storages[gameID] = make(map[uint]*DataStorage)
	}

	if dataStorageManager.storages[gameID][playerID] == nil {
		dataStorageManager.storages[gameID][playerID] = NewDataStorage(gameID, playerID)
	}

	return dataStorageManager.storages[gameID][playerID]
}

// CleanupExpiredData 清理过期数据
func CleanupExpiredData() {
	// 这里可以实现定期清理过期数据的逻辑
	// 例如：清理超过30天未更新的数据
}

// BackupData 备份玩家数据
func (ds *DataStorage) BackupData() (string, error) {
	data, err := ds.GetAllData()
	if err != nil {
		return "", err
	}

	backup := map[string]interface{}{
		"game_id":   ds.gameID,
		"player_id": ds.playerID,
		"data":      data,
		"backup_at": time.Now().Format(time.RFC3339),
		"backup_id": time.Now().UnixNano(),
	}

	jsonData, err := json.Marshal(backup)
	if err != nil {
		return "", err
	}

	// 这里可以将备份数据保存到专门的备份表或文件中
	// 目前先返回JSON字符串
	return string(jsonData), nil
}

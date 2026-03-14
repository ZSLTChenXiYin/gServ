package repository

import (
	"fmt"
	"gServ/core/config"
	"gServ/core/log"
	"gServ/pkg/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	ONLINE_PLAYER_START_ID = 100000
)

var (
	database *gorm.DB
)

func Init() error {
	// 创建 Gorm 配置
	gorm_config := &gorm.Config{
		Logger: log.GetZapGormLogger(),
	}

	database_driver := config.GetConfig().Database.Driver
	database_dsn := config.GetConfig().Database.DSN

	// 创建数据库连接
	var err error
	switch database_driver {
	case config.DATABASE_DRIVER_SQLITE:
		database, err = gorm.Open(sqlite.Open(database_dsn), gorm_config)
	case config.DATABASE_DRIVER_MYSQL:
		database, err = gorm.Open(mysql.Open(database_dsn), gorm_config)
	default:
		return fmt.Errorf("无效数据库驱动: %s", database_driver)
	}
	if err != nil {
		return fmt.Errorf("数据库连接错误: %v", err)
	}

	// 创建数据库表
	tables := []any{
		&model.Game{},
		&model.Player{},
		&model.PlayerDataArchive{},
	}

	err = database.AutoMigrate(tables...)
	if err != nil {
		return fmt.Errorf("自动迁移错误: %v", err)
	}

	var count int64
	database.Model(&model.Player{}).Count(&count)
	if count == 0 && database_driver == config.DATABASE_DRIVER_MYSQL {
		tablename := model.Player{}.TableName()
		err = database.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d", tablename, ONLINE_PLAYER_START_ID)).Error
	}

	return nil
}

package config

import (
	"fmt"
	"sync"
)

var (
	viper_config_lock sync.Mutex
	viper_config      *ViperConfig
)

func Init(name string) error {
	err := initViperConfig(name)
	if err != nil {
		return fmt.Errorf("初始化viper配置错误: %v", err)
	}

	// 加载配置
	viper_config, err = loadViperConfig()
	if err != nil {
		return fmt.Errorf("加载viper配置错误: %v", err)
	}

	return err
}

func GetConfig() *ViperConfig {
	return viper_config
}

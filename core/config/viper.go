package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func initViperConfig(name string) error {
	// 设置配置文件名（不含扩展名）
	viper.SetConfigName(fmt.Sprintf("%s.conf", name))
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 添加配置文件搜索路径
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper读取配置错误: %v", err)
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		var err error
		// 加载配置
		viper_config, err = loadViperConfig()
		if err != nil {
			panic(fmt.Errorf("viper加载配置错误: %v", err))
		}
	})

	return nil
}

func loadViperConfig() (*ViperConfig, error) {
	// 并发安全
	viper_config_lock.Lock()
	defer viper_config_lock.Unlock()

	viper_config := &ViperConfig{}
	err := viper.Unmarshal(viper_config)
	return viper_config, err
}

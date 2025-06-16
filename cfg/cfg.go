package cfg

import (
	"fmt"
	"github.com/before80/lot/lg"
	"github.com/spf13/viper"
)

// DefaultConfig 定义整体 JSON 文件的结构
type DefaultConfig struct {
	Db         DBConfig `mapstructure:"db"`
	EnvIsLocal int      `mapstructure:"env_is_local"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

var Default DefaultConfig

func init() {
	var err error
	Default, err = getDefaultConfigInfo()
	if err != nil {
		panic(fmt.Sprintf("获取默认配置信息出现错误：%v", err))
	}
}

// GetDefaultConfigInfo 获取默认配置信息
func getDefaultConfigInfo() (defaultConfig DefaultConfig, err error) {
	viper.SetConfigName("Default")  // 配置文件名称（不包含扩展名）
	viper.SetConfigType("toml")     // 配置文件类型
	viper.AddConfigPath("./config") // 配置文件所在目录

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			lg.ErrorToFile(fmt.Sprintf("配置文件未找到: %v\n", err))
		} else {
			lg.ErrorToFile(fmt.Sprintf("读取配置文件时出错: %v\n", err))
		}
		return
	}
	if err = viper.Unmarshal(&defaultConfig); err != nil {
		lg.ErrorToFile(fmt.Sprintf("解析配置信息到对象时出错: %v\n", err))
		return
	}

	return defaultConfig, err
}

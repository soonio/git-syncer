package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config config

type remote struct {
	Origin string `mapstructure:"origin" json:"origin" yaml:"origin"` // 原始origin模版
	New    string `mapstructure:"new" json:"new" yaml:"new"`          // 新的远程地址模版
}

type dict struct {
	Key   string `mapstructure:"key" json:"key" yaml:"key"`       // 新的远程地址模版
	Value string `mapstructure:"value" json:"value" yaml:"value"` // 新的远程地址模版
}

type config struct {
	Storage string   `mapstructure:"storage" json:"storage" yaml:"storage"` // 存储库
	Remote  remote   `mapstructure:"remote" json:"remote" yaml:"remote"`    // 远程配置
	Repo    []string `mapstructure:"repo" json:"repo" yaml:"repo"`          // 仓库信息
	Maps    []dict   `mapstructure:"maps" json:"maps" yaml:"maps"`          // 别名映射
}

func initializeViper(config string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("Can't json decode error: %s \n", err))
	}

	return v
}

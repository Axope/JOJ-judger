package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configs struct {
	Log      LogConfig
	Rabbitmq RabbitmqConfig
	Mongo    MongoConfig
}

var cfg Configs

func InitConfigs() {
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic("unmarshal error")
	}
}

func GetLogConfig() LogConfig {
	return cfg.Log
}
func GetRBTConfig() RabbitmqConfig {
	return cfg.Rabbitmq
}
func GetMongoConfig() MongoConfig {
	return cfg.Mongo
}

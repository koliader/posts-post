package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Environment   string `mapstructure:"ENVIRONMENT"`
	RbmUrl        string `mapstructure:"RBM_URL"`
	RedisUrl      string `mapstructure:"REDIS_URL"`
	RedisDBNumber int    `mapstructure:"REDIS_DB_NUMBER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("example")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

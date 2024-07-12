package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Token    TokenConfig    `mapstructure:"token"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"database_name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type TokenConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `mapstructure:"refreshTokenDuration"`
	SecretKey            string        `mapstructure:"secretKey"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

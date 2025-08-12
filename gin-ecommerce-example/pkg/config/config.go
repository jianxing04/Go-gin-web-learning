package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	MySQLDSN      string `mapstructure:"MYSQL_DSN"`
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	ESAddresses   string `mapstructure:"ES_ADDRESSES"`
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}
	return nil
}

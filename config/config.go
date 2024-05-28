package config

import (
	"github.com/nkarakotova/lim-repo/flags"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres      flags.PostgresFlags `mapstructure:"postgres"`
	Address       string              `mapstructure:"address"`
	Port          string              `mapstructure:"port"`
	LogLevel      string              `mapstructure:"loglevel"`
	LogFile       string              `mapstructure:"logfile"`
	Mode          string              `mapstructure:"mode"`
	AdminLogin    string              `mapstructure:"admin_login"`
	AdminPassword string              `mapstructure:"admin_password"`

	FirstTrainingTime int `mapstructure:"first_training_time"`
	LastTrainingTime int `mapstructure:"last_training_time"`
}

func (c *Config) ParseConfig(configFileName, pathToConfig string) error {
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType("json")
	v.AddConfigPath(pathToConfig)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}

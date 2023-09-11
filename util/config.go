package util

import (
	"time"

	"github.com/spf13/viper"
)

//Config stores all configuration of the application.
//The values are read by viper from  a config file  or env variables
type Config struct {
	DBDriver            string        `mapStructure:"DB_DRIVER"`
	DBSource            string        `mapStructure:"DB_SOURCE"`
	ServerAddress       string        `mapStructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapStructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapStructure:"ACCCESS_TOKEN_DURATION"`
}

//LoadConfig reads configuration  from  file or env file
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //can read  json xml  and another

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	config.DBDriver = viper.GetString("DB_DRIVER")
	config.DBSource = viper.GetString("DB_SOURCE")
	config.ServerAddress = viper.GetString("SERVER_ADDRESS")
	config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
	config.AccessTokenDuration = viper.GetDuration("ACCCESS_TOKEN_DURATION")

	err = viper.Unmarshal(&config)
	return config, err
}

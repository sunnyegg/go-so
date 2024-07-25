package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config is a struct for holding configuration values
// The values are read by Viper from the environment variables
type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TwitchClientID       string        `mapstructure:"TWITCH_CLIENT_ID"`
	TwitchClientSecret   string        `mapstructure:"TWITCH_SECRET_ID"`
	FeAddress            string        `mapstructure:"FE_ADDRESS"`
	RedirectURI          string        `mapstructure:"REDIRECT_URI"`
}

// LoadConfig loads the configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can be yaml, json, etc
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

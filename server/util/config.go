package util

import (
	"time"

	"github.com/spf13/viper"
)

/*
Config stores all configuration of the application.
The values are read by viper from a config file or environment variable.
*/
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DB_URL               string        `mapstructure:"DB_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYM_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MongoURI             string        `mapstructure:"MONGO_URI"`
}

/* LoadConfig reads configuration from file or environment variables. */
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

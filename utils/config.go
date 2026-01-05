package utils

import "github.com/spf13/viper"

// Config stores all configurations of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"` // this name must match the key in the config file or environment variable
	DBSource string `mapstructure:"DB_SOURCE"` // this name must match the key in the config file or environment variable
	ServerAddress string `mapstructure:"SERVER_ADDRESS"` // this name must match the key in the config file or environment variable
}

// LoadConfig reads configuration from file or environment variables
// it returns the config or an error if the configuration is invalid
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // tell viper to look for config file with this name "app"
	viper.SetConfigType("env") // tell viper to look for config file with this extension "env"

	viper.AutomaticEnv() // allow viper to read environment variables

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config) // unmarshal the config into the Config struct
	return
}
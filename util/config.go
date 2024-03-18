package util

import "github.com/spf13/viper"

// Congif configures application level variables
type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	DBName            string `mapstructure:"DB_NAME"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	TokenSimmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

// LoadConfig reads configuration from file or environment variables and returns a capy of Config struct
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return config, nil
}

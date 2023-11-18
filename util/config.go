package util

import (
	"time"

	"github.com/spf13/viper"
)

// konfig semua konfigurasi
// nilai dibaca oleh viper dengan .env
type Config struct {
	DBDriver                  string        `mapstructure:"DB_DRIVER"`
	DBSource                  string        `mapstructure:"DB_SOURCE"`
	ServerAddress             string        `mapstructure:"SERVER_ADDRESS"`
	PasetoSymmetricKey        string        `mapstructure:"PASETO_SYMMETRIC_KEY"`
	PasetoAccessTokenDuration time.Duration `mapstructure:"PASETO_ACCESS_TOKEN_DURATION"`
}

// fungsi untuk load konfigurasi
func LoadConfig(path string) (config Config, err error) {
	// menambah path
	viper.AddConfigPath(path)
	// Config berdasarkan nama file
	viper.SetConfigName("dev")
	// ekstensi dari file yang akan diambil
	viper.SetConfigType("env")
	// Viper will check for an environment variable any time a viper.Get request is made.
	viper.AutomaticEnv()
	// baca nilai dalam config
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return
	}

	return
}

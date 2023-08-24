package config

import "os"

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

func NewConfig() *Config {
	return &conf
}

type AppConfig struct {
	Port            string
	EnableSonicJson bool
}

type DatabaseConfig struct {
	Url string
}
type CacheConfig struct {
	Url string
}

var conf = Config{
	App: AppConfig{
		Port:            os.Getenv("PORT"),
		EnableSonicJson: os.Getenv("ENABLE_SONIC_JSON") == "1",
	},
	Database: DatabaseConfig{
		Url: os.Getenv("DATABASE_URL"),
	},
	Cache: CacheConfig{
		Url: os.Getenv("CACHE_URL"),
	},
}

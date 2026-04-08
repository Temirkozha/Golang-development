package config

import (
	"os"
)

type Config struct {
	AppPort string
	DBURL   string
	JWTKey  string
}

func NewConfig() *Config {
	return &Config{
		AppPort: "8090",
		DBURL:   os.Getenv("DB_URL"),
		JWTKey:  os.Getenv("JWT_SECRET"),
	}
}
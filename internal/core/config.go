package core

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	cfg.Port 		= os.Getenv("PORT")
	cfg.DBHost 		= os.Getenv("POSTGRES_HOST")
	cfg.DBPort 		= os.Getenv("POSTGRES_PORT")
	cfg.DBUser 		= os.Getenv("POSTGRES_USER")
	cfg.DBPassword 	= os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName 		= os.Getenv("POSTGRES_DB")
	cfg.DBSSLMode 	= os.Getenv("POSTGRES_SSLMODE")

	return cfg, nil
}

func (c *Config) DBConnString() string {
	ssl := c.DBSSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		ssl,
	)
}

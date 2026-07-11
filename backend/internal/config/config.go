package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	JWTSecret        string
	JWTExpHours      int
	ServerPort       string
	FrontendURL      string
	OpenAIAPIKey     string
}

func Load() *Config {
	expHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	return &Config{
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "postgres"),
		DatabaseName:     getEnv("DATABASE_NAME", "concesionaria"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret"),
		JWTExpHours:      expHours,
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:5173"),
		OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DatabaseHost +
		" port=" + c.DatabasePort +
		" user=" + c.DatabaseUser +
		" password=" + c.DatabasePassword +
		" dbname=" + c.DatabaseName +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

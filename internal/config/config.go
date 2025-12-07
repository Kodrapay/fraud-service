package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServiceName   string
	Port          string
	APIKey        string
	PostgresDSN   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load(serviceName, defaultPort string) Config {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0 // Default to DB 0 if conversion fails
	}

	return Config{
		ServiceName:   serviceName,
		Port:          getEnv("PORT", defaultPort),
		APIKey:        getEnv("API_KEY", ""),
		PostgresDSN:   getEnv("POSTGRES_DSN", ""),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

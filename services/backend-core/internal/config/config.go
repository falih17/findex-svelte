package config

import (
	"os"
)

type Config struct {
	Environment  string
	Port         string
	DatabaseURL  string
	RedisHost    string
	RedisPort    string
	RedisPass    string
	AIServiceURL string
	JWTSecret    string
}

func Load() *Config {
	return &Config{
		Environment:  getEnv("ENVIRONMENT", "development"),
		Port:         getEnv("BACKEND_PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://findex_user:findex_secret_pwd@postgres:5432/findex_db?sslmode=disable"),
		RedisHost:    getEnv("REDIS_HOST", "redis"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		RedisPass:    getEnv("REDIS_PASSWORD", "redis_secret_pwd"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://ai-face-service:8000"),
		JWTSecret:    getEnv("JWT_SECRET", "super_secret_jwt_key_at_least_32_characters_long!"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

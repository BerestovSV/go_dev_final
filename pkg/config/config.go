package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TokenDuration time.Duration
	DBFile        string
	Password      string
	JWTSecret     string
	Port          string
}

func Load() *Config {
	port := getEnv("TODO_PORT", "7540")
	if port[0] != ':' {
		port = ":" + port
	}

	cfg := &Config{
		TokenDuration: parseDuration("TOKEN_DURATION", 8*time.Hour),
		DBFile:        getEnv("TODO_DBFILE", "scheduler.db"),
		Password:      getEnv("TODO_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		Port:          port,
	}

	// Fallback для JWTSecret
	if cfg.JWTSecret == "" && cfg.Password != "" {
		cfg.JWTSecret = getPasswordHash(cfg.Password)
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	return defaultValue
}

func getPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

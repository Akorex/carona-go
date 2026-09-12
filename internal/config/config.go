package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port              string `validate:"required,numeric"`
	DatabaseURL       string `validate:"required"`
	JWTSecret         string `validate:"required,min=16"`
	JWTExpiresInHours int    `validate:"required,min=1"`
}

// Load reads .env, loads environment variables, and validates them (Yup/Zod style).
func Load() (*Config, error) {
	// Load .env file if present
	_ = godotenv.Load()

	cfg := &Config{
		Port:              getEnv("PORT", "8000"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/caronago?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-key-that-is-at-least-16-chars"),
		JWTExpiresInHours: getEnvAsInt("JWT_EXPIRES_IN_HOURS", 24),
	}

	// Validate the configuration struct
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid environment configuration: %w", err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

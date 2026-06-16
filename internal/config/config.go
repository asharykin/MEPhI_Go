package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	Port         string
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://bank_user:bank_password@localhost:5433/bank_db?sslmode=disable"),
		Port:         getEnv("PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", "your-email@gmail.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "your-password"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

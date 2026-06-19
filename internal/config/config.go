package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	Security SecurityConfig
}

type DatabaseConfig struct {
	connectionString string
}

type JWTConfig struct {
	Secret string
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type SecurityConfig struct {
	HMACSecret    string
	PGPPublicKey  string
	PGPPrivateKey string
}

func LoadConfig() *Config {
	return &Config{
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
		SMTP:     loadSMTPConfig(),
		Security: loadSecurityConfig(),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		connectionString: getEnv("DB_URL", "postgres://bank_user:bank_password@localhost:5433/bank_db?sslmode=disable"),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
	}
}

func loadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host: getEnv("SMTP_HOST", "smtp.example.com"),
		Port: getIntEnv("SMTP_PORT", 587),
		User: getEnv("SMTP_USER", "noreply@example.com"),
		Pass: getEnv("SMTP_PASS", "password"),
	}
}

func loadSecurityConfig() SecurityConfig {
	public_data, err := os.ReadFile("./pgp_public.key")
	if err != nil {
		panic("Failed to read PGP key file")
	}
	private_data, err := os.ReadFile("./pgp_private.key")
	if err != nil {
		panic("Failed to read PGP private key file")
	}
	return SecurityConfig{
		HMACSecret:    getEnv("HMAC_SECRET", "hmac-secret-key-change-in-production"),
		PGPPublicKey:  string(public_data),
		PGPPrivateKey: string(private_data),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *DatabaseConfig) GetConnectionString() string {
	return c.connectionString
}

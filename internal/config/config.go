package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Port       string
	AppEnv     string
	JWTSecret  string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	AdminName string
	AdminPass string

	TLSCert string
	TLSKey  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "production"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "mydb"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		AdminName:  getEnv("ADMIN_NAME", "admin"),
		AdminPass:  getEnv("ADMIN_PASSWORD", "password"),
		TLSCert:    getEnv("TLS_CERT", ""),
		TLSKey:     getEnv("TLS_KEY", ""),
	}

	zap.L().Info("Configuration loaded")

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

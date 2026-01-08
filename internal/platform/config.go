// Package platform provides platform-specific configurations and utilities.
package platform

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	Environment          string
	SecretsEncryptionKey string
	ComputeBackend       string
	DefaultVPCCIDR       string
	NetworkPoolStart     string
	NetworkPoolEnd       string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://cloud:cloud@localhost:5433/thecloud"),
		Environment:          getEnv("APP_ENV", "development"),
		SecretsEncryptionKey: os.Getenv("SECRETS_ENCRYPTION_KEY"),
		ComputeBackend:       getEnv("COMPUTE_BACKEND", "docker"),
		DefaultVPCCIDR:       getEnv("DEFAULT_VPC_CIDR", "10.0.0.0/16"),
		NetworkPoolStart:     getEnv("NETWORK_POOL_START", "192.168.100.0"),
		NetworkPoolEnd:       getEnv("NETWORK_POOL_END", "192.168.200.255"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

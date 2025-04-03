package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TokenTTL    time.Duration
	StoragePath string
	GRPC        GRPSConfig
}

type GRPSConfig struct {
	ApiPort int
	Timeout time.Duration

	Host     string
	Db       string
	User     string
	Password string
}

func MustLoad() *Config {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	var cfg Config

	// Получение значений из переменных окружения
	cfg.TokenTTL = getEnvAsDuration("TOKEN_TTL", 1*time.Hour)
	cfg.StoragePath = os.Getenv("STORAGE_PATH")
	cfg.GRPC.ApiPort = getEnvAsInt("GRPC_API", 50051)
	cfg.GRPC.Timeout = getEnvAsDuration("GRPC_TIMEOUT", 5*time.Second)
	cfg.GRPC.Host = os.Getenv("DB_HOST")
	cfg.GRPC.Db = os.Getenv("GPRC_DB")
	cfg.GRPC.User = os.Getenv("DB_USER")
	cfg.GRPC.Password = os.Getenv("DB_PASSWORD")

	return &cfg
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}

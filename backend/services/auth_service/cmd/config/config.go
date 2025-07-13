package config

import (
	"os"
	"strconv"
	"time"
)

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type ServiceConfig struct {
	Profile_service_utl string
	Jwt_secret          string
}

func LoadServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		Profile_service_utl: os.Getenv("PROFILE_SERVICE_URL"),
		Jwt_secret:          os.Getenv("JWT_SECRET"),
	}
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Name:            os.Getenv("DB_NAME"),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

package utils

import (
	"time"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port     string
	CertFile string
	KeyFile  string
}

// LoadDBConfig loads database configuration from environment
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Host:            GetEnv("DB_HOST", "localhost"),
		Port:            GetEnv("DB_PORT", "3306"),
		User:            GetEnv("DB_USER", "api_user"),
		Password:        GetEnv("DB_PASSWORD", "api_password"),
		Database:        GetEnv("DB_NAME", "SCHOOL_DB"),
		MaxOpenConns:    GetEnvAsInt("DB_MAX_OPEN_CONNS", 10),
		MaxIdleConns:    GetEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: GetEnvAsDuration("DB_CONN_MAX_LIFETIME", 3*time.Minute),
	}
}

// LoadServerConfig loads server configuration from environment
func LoadServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:     GetEnv("SERVER_PORT", "8080"),
		CertFile: GetEnv("TLS_CERT_FILE", "cert.pem"),
		KeyFile:  GetEnv("TLS_KEY_FILE", "key.pem"),
	}
}

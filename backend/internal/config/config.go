package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Session     SessionConfig
	RateLimit   RateLimitConfig
	CORS        CORSConfig
	WebSocket   WebSocketConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string
}

// SessionConfig holds session-related configuration
type SessionConfig struct {
	Duration time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	Enabled           bool
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	PingPeriod      time.Duration
	PongWait        time.Duration
	WriteWait       time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:         getEnv("PORT", "8000"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/forum.db"),
		},
		Session: SessionConfig{
			Duration: getDurationEnv("SESSION_DURATION", 24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_RPM", 100),
			Enabled:           getBoolEnv("RATE_LIMIT_ENABLED", true),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "*")},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  getIntEnv("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getIntEnv("WS_WRITE_BUFFER_SIZE", 1024),
			PingPeriod:      getDurationEnv("WS_PING_PERIOD", 54*time.Second),
			PongWait:        getDurationEnv("WS_PONG_WAIT", 60*time.Second),
			WriteWait:       getDurationEnv("WS_WRITE_WAIT", 10*time.Second),
		},
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}
	if c.Session.Duration < time.Minute {
		return fmt.Errorf("session duration must be at least 1 minute")
	}
	if c.RateLimit.RequestsPerMinute < 1 {
		return fmt.Errorf("rate limit must be at least 1 request per minute")
	}
	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}


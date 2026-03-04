package config

import (
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Session     SessionConfig
	RateLimit   RateLimitConfig
	CORS        CORSConfig
	Websocket   WebSocketConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/database/forum.db"),
		},
		Session: SessionConfig{
			Duration: getEnvDuration("SESSION_DURATION", 24*time.Hour),
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
		Websocket: WebSocketConfig{
			ReadBufferSize:  getIntEnv("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getIntEnv("WS_WRITE_BUFFER_SIZE", 1024),
			PingPeriod:      getEnvDuration("WS_PING_PERIOD", 54*time.Second),
			PongWait:        getEnvDuration("WS_PONG_WAIT", 60*time.Second),
			WriteWait:       getEnvDuration("WS_WRITE_WAIT", 10*time.Second),
		},
	}

	return cfg, nil
}

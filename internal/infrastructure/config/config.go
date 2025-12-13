package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const defaultJWTSecret = "supersecretjwtkeythatshouldbechangedinproduction"

// Config aggregates runtime configuration values for the application.
type Config struct {
	DatabaseURL     string
	APIPort         string
	SchedulerUserID uint
	JWTSecret       string
}

// Load reads configuration from environment variables and applies sane defaults.
func Load() (*Config, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	apiPort := strings.TrimSpace(os.Getenv("API_PORT"))
	if apiPort == "" {
		apiPort = "8080"
	}
	if _, err := strconv.Atoi(apiPort); err != nil {
		return nil, fmt.Errorf("invalid API_PORT: %w", err)
	}

	schedulerID, err := parseSchedulerUserID()
	if err != nil {
		return nil, err
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET_KEY"))
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
		log.Println("WARNING: JWT_SECRET_KEY not set, using default development key")
	}

	return &Config{
		DatabaseURL:     databaseURL,
		APIPort:         apiPort,
		SchedulerUserID: schedulerID,
		JWTSecret:       jwtSecret,
	}, nil
}

func parseSchedulerUserID() (uint, error) {
	raw := strings.TrimSpace(os.Getenv("SCHEDULER_USER_ID"))
	if raw == "" {
		return 1, nil
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid SCHEDULER_USER_ID: %w", err)
	}

	return uint(value), nil
}

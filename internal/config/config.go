package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultCheckTimeout  = 3 * time.Second
	defaultCheckInterval = 30 * time.Second
)

type Config struct {
	CheckTimeout  time.Duration
	CheckInterval time.Duration
}

func Load() (Config, error) {
	checkTimeout, err := loadDuration(
		"NETWATCH_CHECK_TIMEOUT",
		defaultCheckTimeout,
	)
	if err != nil {
		return Config{}, err
	}

	checkInterval, err := loadDuration(
		"NETWATCH_CHECK_INTERVAL",
		defaultCheckInterval,
	)
	if err != nil {
		return Config{}, err
	}

	return Config{
		CheckTimeout:  checkTimeout,
		CheckInterval: checkInterval,
	}, nil
}

func loadDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return duration, nil
}

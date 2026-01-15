package config

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

// Test profiles for different scenarios
type Profile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Requests    int    `json:"requests"`
	Concurrency int    `json:"concurrency"`
	RPS         int    `json:"rps"`
	Duration    int    `json:"duration_minutes"`
	ThinkTimeMs int    `json:"think_time_ms"`
}

type Config struct {
	URL          string
	Requests     int
	RPS          int
	Concurrency  int
	Timeout      time.Duration
	ThinkTime    time.Duration
	Profile      string
	Endpoints    []string
	Authorized   bool
	ShowProgress bool
}

type Stats struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	StatusCodes     sync.Map
	MinLatency      int64
	MaxLatency      int64
	TotalLatency    int64
	StartTime       time.Time
	EndTime         time.Time
}

func (cfg *Config) ValidateConfig() error {
	if cfg.URL == "" {
		return fmt.Errorf("URL is required (use -url flag)")
	}

	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if cfg.Requests < 1 {
		return fmt.Errorf("requests must be at least 1")
	}

	if cfg.RPS < 1 {
		return fmt.Errorf("RPS must be at least 1")
	}

	if cfg.Concurrency < 1 {
		return fmt.Errorf("concurrency must be at least 1")
	}

	return nil
}

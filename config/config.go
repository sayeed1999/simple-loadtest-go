package config

import (
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

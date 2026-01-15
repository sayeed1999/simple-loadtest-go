package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sayeed1999/simple-loadtest-go/config"
)

var DefaultProfiles = map[string]config.Profile{
	"normal": {
		Name:        "Normal Traffic",
		Description: "Regular daily traffic pattern",
		Requests:    10000,
		Concurrency: 50,
		RPS:         100,
		Duration:    10,
		ThinkTimeMs: 500,
	},
	"peak": {
		Name:        "Peak Hours",
		Description: "Evening rush hour traffic",
		Requests:    100000,
		Concurrency: 200,
		RPS:         500,
		Duration:    30,
		ThinkTimeMs: 300,
	},
	"flash-sale": {
		Name:        "Flash Sale",
		Description: "High-intensity flash sale event",
		Requests:    500000,
		Concurrency: 1000,
		RPS:         2000,
		Duration:    60,
		ThinkTimeMs: 100,
	},
	"stress": {
		Name:        "Stress Test",
		Description: "Maximum load to find breaking point",
		Requests:    1000000,
		Concurrency: 2000,
		RPS:         5000,
		Duration:    120,
		ThinkTimeMs: 50,
	},
}

func main() {
	cfg := parseFlags()

	if cfg.URL == "" && cfg.Profile == "" {
		showProfiles()
		fmt.Println("Use ./loadtest -help for more configuration information.\n")
		os.Exit(1)
	}

	// Load profile if specified
	if cfg.Profile != "" {
		if err := loadProfile(cfg); err != nil {
			fmt.Printf("❌ Profile error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := cfg.ValidateConfig(); err != nil {
		fmt.Printf("❌ config.Configuration error: %v\n", err)
		os.Exit(1)
	}

	showWarning()
	if !cfg.Authorized && !getConfirmation() {
		fmt.Println("Load test cancelled.")
		os.Exit(0)
	}

	displayTestInfo(cfg)
	stats := runLoadTest(cfg)
	printResults(stats, cfg)
}

func parseFlags() *config.Config {
	cfg := &config.Config{}

	flag.StringVar(&cfg.URL, "url", "", "Target URL to test (required)")
	flag.StringVar(&cfg.Profile, "profile", "", "Use predefined profile: normal, peak, flash-sale, stress")
	flag.IntVar(&cfg.Requests, "requests", 100, "Total number of requests")
	flag.IntVar(&cfg.RPS, "rps", 10, "Requests per second")
	flag.IntVar(&cfg.Concurrency, "concurrency", 5, "Number of concurrent workers")
	timeoutSec := flag.Int("timeout", 30, "Request timeout in seconds")
	thinkTimeMs := flag.Int("think-time", 100, "Delay between requests in milliseconds")
	flag.BoolVar(&cfg.Authorized, "authorized", false, "Skip authorization prompt (use only for your own systems)")
	flag.BoolVar(&cfg.ShowProgress, "progress", true, "Show progress during test")

	listProfiles := flag.Bool("list-profiles", false, "List available test profiles")

	flag.Parse()

	if *listProfiles {
		showProfiles()
		os.Exit(0)
	}

	cfg.Timeout = time.Duration(*timeoutSec) * time.Second
	cfg.ThinkTime = time.Duration(*thinkTimeMs) * time.Millisecond

	return cfg
}

func showProfiles() {
	fmt.Println("\n📋 Available Test Profiles for Ecommerce Load Testing")
	fmt.Println(strings.Repeat("=", 80))

	profiles := []string{"normal", "peak", "flash-sale", "stress"}
	for _, key := range profiles {
		p := DefaultProfiles[key]
		fmt.Printf("\n🏷️  %s (%s)\n", p.Name, key)
		fmt.Printf("   %s\n", p.Description)
		fmt.Printf("   • Requests: %d\n", p.Requests)
		fmt.Printf("   • Concurrency: %d users\n", p.Concurrency)
		fmt.Printf("   • RPS: %d req/sec\n", p.RPS)
		fmt.Printf("   • Duration: ~%d minutes\n", p.Duration)
		fmt.Printf("   • Think Time: %dms\n", p.ThinkTimeMs)
	}

	fmt.Printf("\n💡 Usage: ./loadtest -url http://your-site.com -profile [normal,peak,flash-sale,stress]\n\n")
}

func loadProfile(cfg *config.Config) error {
	profile, exists := DefaultProfiles[cfg.Profile]
	if !exists {
		return fmt.Errorf("unknown profile '%s'. Use -list-profiles to see available profiles", cfg.Profile)
	}

	// Only override if not explicitly set by user
	if !isFlagPassed("requests") {
		cfg.Requests = profile.Requests
	}
	if !isFlagPassed("concurrency") {
		cfg.Concurrency = profile.Concurrency
	}
	if !isFlagPassed("rps") {
		cfg.RPS = profile.RPS
	}
	if !isFlagPassed("think-time") {
		cfg.ThinkTime = time.Duration(profile.ThinkTimeMs) * time.Millisecond
	}

	return nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func showWarning() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚠️  LEGAL WARNING - READ CAREFULLY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("This tool is for testing YOUR OWN infrastructure only.")
	fmt.Println("Using this against systems you don't own or operate is ILLEGAL and a CYBERCRIME.")
	fmt.Println("It may result in IMPRISONMENT or other serious punishment if reported & tracked via IP.")
	fmt.Println("You must have explicit written authorization before testing any system.")
	fmt.Println(strings.Repeat("=", 80))
}

func getConfirmation() bool {
	fmt.Print("\n✋ Do you own this system and have authorization to test it? (yes/no): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "yes"
}

func displayTestInfo(cfg *config.Config) {
	fmt.Printf("\n🚀 Starting Load Test\n")
	fmt.Println(strings.Repeat("-", 80))

	if cfg.Profile != "" {
		profile := DefaultProfiles[cfg.Profile]
		fmt.Printf("Profile:            %s - %s\n", profile.Name, profile.Description)
	}

	fmt.Printf("Target:             %s\n", cfg.URL)
	fmt.Printf("Total Requests:     %s\n", formatNumber(cfg.Requests))
	fmt.Printf("Concurrent Users:   %d\n", cfg.Concurrency)
	fmt.Printf("Max RPS:            %d req/sec\n", cfg.RPS)
	fmt.Printf("Think Time:         %v\n", cfg.ThinkTime)
	fmt.Printf("Timeout:            %v\n", cfg.Timeout)

	estimatedDuration := float64(cfg.Requests) / float64(cfg.RPS)
	fmt.Printf("Estimated Duration: %.1f seconds (%.1f minutes)\n", estimatedDuration, estimatedDuration/60)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
}

func runLoadTest(cfg *config.Config) *config.Stats {
	stats := &config.Stats{
		StartTime:  time.Now(),
		MinLatency: int64(^uint64(0) >> 1), // Max int64
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.Concurrency * 2,
			MaxIdleConnsPerHost: cfg.Concurrency * 2,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Rate limiter
	interval := time.Second / time.Duration(cfg.RPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	workChan := make(chan struct{}, cfg.Requests)

	// Progress reporter
	if cfg.ShowProgress {
		go progressReporter(ctx, stats, cfg.Requests)
	}

	// Start workers
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, &wg, workChan, cfg, client, stats)
	}

	// Send work with rate limiting
	for i := 0; i < cfg.Requests; i++ {
		<-ticker.C
		workChan <- struct{}{}
	}

	close(workChan)
	wg.Wait()

	stats.EndTime = time.Now()
	return stats
}

func worker(ctx context.Context, wg *sync.WaitGroup, work <-chan struct{}, cfg *config.Config, client *http.Client, stats *config.Stats) {
	defer wg.Done()

	for range work {
		select {
		case <-ctx.Done():
			return
		default:
			makeRequest(client, cfg.URL, stats)

			// Add random think time variation (±20%)
			thinkTime := cfg.ThinkTime
			if thinkTime > 0 {
				variation := time.Duration(rand.Intn(int(thinkTime.Milliseconds())/5)) * time.Millisecond
				if rand.Intn(2) == 0 {
					thinkTime += variation
				} else {
					thinkTime -= variation
				}
				time.Sleep(thinkTime)
			}
		}
	}
}

func makeRequest(client *http.Client, targetURL string, stats *config.Stats) {
	atomic.AddInt64(&stats.TotalRequests, 1)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		atomic.AddInt64(&stats.FailedRequests, 1)
		incrementStatusCode(stats, 0)
		return
	}

	req.Header.Set("User-Agent", "go-loadtest/2.0 (Professional Load Testing)")
	req.Header.Set("Accept", "text/html,application/json")

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		atomic.AddInt64(&stats.FailedRequests, 1)
		incrementStatusCode(stats, 0)
		return
	}
	defer resp.Body.Close()

	// Read and discard body to reuse connection
	ioutil.ReadAll(resp.Body)

	// Update latency stats
	updateLatency(stats, latency)

	// Track status codes
	incrementStatusCode(stats, resp.StatusCode)

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		atomic.AddInt64(&stats.SuccessRequests, 1)
	} else {
		atomic.AddInt64(&stats.FailedRequests, 1)

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("\r⚠️  Server rate limiting (429) - consider reducing RPS\n")
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			fmt.Printf("\r⚠️  Service unavailable (503) - server may be overloaded\n")
		}
	}
}

func updateLatency(stats *config.Stats, latency int64) {
	atomic.AddInt64(&stats.TotalLatency, latency)

	for {
		min := atomic.LoadInt64(&stats.MinLatency)
		if latency >= min || atomic.CompareAndSwapInt64(&stats.MinLatency, min, latency) {
			break
		}
	}

	for {
		max := atomic.LoadInt64(&stats.MaxLatency)
		if latency <= max || atomic.CompareAndSwapInt64(&stats.MaxLatency, max, latency) {
			break
		}
	}
}

func incrementStatusCode(stats *config.Stats, code int) {
	key := fmt.Sprintf("%d", code)
	val, _ := stats.StatusCodes.LoadOrStore(key, new(int64))
	atomic.AddInt64(val.(*int64), 1)
}

func progressReporter(ctx context.Context, stats *config.Stats, total int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			completed := atomic.LoadInt64(&stats.TotalRequests)
			success := atomic.LoadInt64(&stats.SuccessRequests)
			failed := atomic.LoadInt64(&stats.FailedRequests)

			percentage := float64(completed) / float64(total) * 100
			elapsed := time.Since(stats.StartTime).Seconds()
			currentRPS := float64(completed) / elapsed

			fmt.Printf("\r📊 Progress: %d/%d (%.1f%%) | Success: %d | Failed: %d | RPS: %.1f   ",
				completed, total, percentage, success, failed, currentRPS)
		}
	}
}

func printResults(stats *config.Stats, cfg *config.Config) {
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Printf("\n\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 LOAD TEST RESULTS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("\n⏱️  Timing:\n")
	fmt.Printf("   Total Duration:       %v\n", duration.Round(time.Millisecond))
	fmt.Printf("   Started:              %s\n", stats.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Ended:                %s\n", stats.EndTime.Format("2006-01-02 15:04:05"))

	fmt.Printf("\n📈 Requests:\n")
	fmt.Printf("   Total Requests:       %s\n", formatNumber(int(stats.TotalRequests)))
	fmt.Printf("   Successful:           %s (%.1f%%)\n",
		formatNumber(int(stats.SuccessRequests)),
		float64(stats.SuccessRequests)/float64(stats.TotalRequests)*100)
	fmt.Printf("   Failed:               %s (%.1f%%)\n",
		formatNumber(int(stats.FailedRequests)),
		float64(stats.FailedRequests)/float64(stats.TotalRequests)*100)

	fmt.Printf("\n⚡ Performance:\n")
	avgRPS := float64(stats.TotalRequests) / duration.Seconds()
	fmt.Printf("   Average RPS:          %.2f requests/sec\n", avgRPS)
	fmt.Printf("   Target RPS:           %d requests/sec\n", cfg.RPS)
	fmt.Printf("   Concurrent Users:     %d\n", cfg.Concurrency)

	if stats.TotalRequests > 0 {
		avgLatency := stats.TotalLatency / stats.TotalRequests
		fmt.Printf("\n⏲️  Latency:\n")
		fmt.Printf("   Average:              %dms\n", avgLatency)
		fmt.Printf("   Min:                  %dms\n", stats.MinLatency)
		fmt.Printf("   Max:                  %dms\n", stats.MaxLatency)
	}

	fmt.Printf("\n📋 HTTP Status Codes:\n")
	stats.StatusCodes.Range(func(key, value interface{}) bool {
		code := key.(string)
		count := atomic.LoadInt64(value.(*int64))
		percentage := float64(count) / float64(stats.TotalRequests) * 100
		fmt.Printf("   %s: %s (%.1f%%)\n", code, formatNumber(int(count)), percentage)
		return true
	})

	fmt.Println(strings.Repeat("=", 80))

	// Performance assessment
	assessPerformance(stats, cfg, avgRPS)
}

func assessPerformance(stats *config.Stats, cfg *config.Config, avgRPS float64) {
	fmt.Printf("\n💡 Performance Assessment:\n")

	successRate := float64(stats.SuccessRequests) / float64(stats.TotalRequests) * 100

	if successRate >= 99.5 {
		fmt.Println("   ✅ Excellent - System handled load very well")
	} else if successRate >= 95 {
		fmt.Println("   ✓ Good - System performed adequately")
	} else if successRate >= 90 {
		fmt.Println("   ⚠️  Fair - Some issues detected, investigate errors")
	} else {
		fmt.Println("   ❌ Poor - Significant issues, system may be overloaded")
	}

	if avgRPS < float64(cfg.RPS)*0.8 {
		fmt.Println("   ⚠️  Could not achieve target RPS - system may be bottlenecked")
	}

	if stats.TotalRequests > 0 {
		avgLatency := stats.TotalLatency / stats.TotalRequests
		if avgLatency > 2000 {
			fmt.Println("   ⚠️  High average latency - check server performance")
		}
	}

	fmt.Println()
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 100000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	if n < 10000000 {
		return fmt.Sprintf("%d,%02d,%03d", n/100000, (n/1000)%100, n%1000)
	}
	return fmt.Sprintf("%d", n)
}

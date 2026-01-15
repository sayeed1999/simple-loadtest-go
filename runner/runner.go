package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sayeed1999/simple-loadtest-go/config"
)

func Run(cfg *config.Config) *config.Stats {
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
			// Stop progress bar when test is done
			completed := atomic.LoadInt64(&stats.TotalRequests)
			success := atomic.LoadInt64(&stats.SuccessRequests)
			failed := atomic.LoadInt64(&stats.FailedRequests)
			elapsed := time.Since(stats.StartTime).Seconds()
			currentRPS := float64(completed) / elapsed
			fmt.Printf("\r📊 Progress: %d/%d (100.0%%) | Success: %d | Failed: %d | RPS: %.1f   \n",
				total, total, success, failed, currentRPS)
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

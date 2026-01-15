# simple-loadtest-go

An educational HTTP load testing tool written in Go for ecommerce and web applications. Includes predefined profiles for different traffic scenarios.

## ⚠️ LEGAL WARNING

**THIS TOOL IS FOR TESTING YOUR OWN INFRASTRUCTURE ONLY.**

- Using this tool against systems you don't own or operate is **ILLEGAL and a CYBERCRIME**
- It may result in **IMPRISONMENT or other serious punishment** if reported & tracked via IP
- You must have **explicit written authorization** before testing any system
- The authors are not responsible for misuse of this tool

## Features

✅ **Predefined Test Profiles** for ecommerce scenarios (normal, peak, flash-sale, stress)
✅ Realistic traffic simulation with variable think times
✅ Detailed performance metrics and latency tracking
✅ HTTP status code distribution
✅ Real-time progress reporting
✅ Connection pooling and keep-alive
✅ Comprehensive result analysis
✅ Authorization confirmation for safety

## Test Profiles

The tool includes 4 predefined profiles for common ecommerce scenarios:

### 🛒 Normal Traffic
**Use case:** Regular daily traffic pattern
- **Requests:** 10,000
- **Concurrent Users:** 50
- **RPS:** 100 req/sec
- **Duration:** ~10 minutes
- **Think Time:** 500ms

### 🌆 Peak Hours
**Use case:** Evening rush hour traffic
- **Requests:** 100,000
- **Concurrent Users:** 200
- **RPS:** 500 req/sec
- **Duration:** ~30 minutes
- **Think Time:** 300ms

### 🔥 Flash Sale
**Use case:** High-intensity flash sale events
- **Requests:** 500,000
- **Concurrent Users:** 1,000
- **RPS:** 2,000 req/sec
- **Duration:** ~60 minutes
- **Think Time:** 100ms

### 💪 Stress Test
**Use case:** Maximum load to find breaking point
- **Requests:** 1,000,000
- **Concurrent Users:** 2,000
- **RPS:** 5,000 req/sec
- **Duration:** ~120 minutes
- **Think Time:** 50ms

## Installation

```bash
# Clone the repository
git clone https://github.com/sayeed1999/simple-loadtest-go.git
cd simple-loadtest-go

# Build (For Windows)
GOOS=windows go build -o ./output/loadtest.exe main.go

# Build (In Linux/Mac)
GOOS=linux go build -o ./output/loadtest main.go

## The GOOS param is not necessary when building for same OS from same OS.
```

## Usage

### List Available Profiles

```bash
./loadtest -list-profiles
```

### Quick Start with Profiles

```bash
# Normal traffic test
./loadtest -url http://localhost:8080 -profile normal

# Peak hour simulation
./loadtest -url http://your-site.com -profile peak

# Flash sale stress test
./loadtest -url http://your-site.com -profile flash-sale

# Maximum stress test
./loadtest -url http://your-site.com -profile stress
```

### Custom Configuration

Override profile settings:

```bash
# Use peak profile but with custom RPS
./loadtest -url http://localhost:8080 -profile peak -rps 1000

# Custom test without profile
./loadtest \
  -url http://localhost:8080/api/products \
  -requests 50000 \
  -rps 500 \
  -concurrency 100 \
  -timeout 30 \
  -think-time 200
```

### Skip Authorization Prompt (for CI/CD)

```bash
./loadtest -url http://localhost:8080 -profile normal -authorized
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-url` | Target URL to test (required) | - |
| `-profile` | Use predefined profile: normal, peak, flash-sale, stress | - |
| `-requests` | Total number of requests | 100 |
| `-rps` | Requests per second | 10 |
| `-concurrency` | Number of concurrent workers | 5 |
| `-timeout` | Request timeout in seconds | 30 |
| `-think-time` | Delay between requests (ms) | 100 |
| `-authorized` | Skip authorization prompt | false |
| `-progress` | Show progress during test | true |
| `-list-profiles` | List available test profiles | - |

## Example Output

```
=================================================================================
⚠️  LEGAL WARNING - READ CAREFULLY
=================================================================================
This tool is for testing YOUR OWN infrastructure only.
Using this against systems you don't own or operate is ILLEGAL and a CYBERCRIME.
It may result in IMPRISONMENT or other serious punishment if reported & tracked via IP.
You must have explicit written authorization before testing any system.
=================================================================================

✋ Do you own this system and have authorization to test it? (yes/no): yes

🚀 Starting Load Test
--------------------------------------------------------------------------------
Profile:            Peak Hours - Evening rush hour traffic
Target:             http://localhost:8080
Total Requests:     100,000
Concurrent Users:   200
Max RPS:            500 req/sec
Think Time:         300ms
Timeout:            30s
Estimated Duration: 200.0 seconds (3.3 minutes)
--------------------------------------------------------------------------------

📊 Progress: 45832/100000 (45.8%) | Success: 45750 | Failed: 82 | RPS: 458.3

=================================================================================
📊 LOAD TEST RESULTS
=================================================================================

⏱️  Timing:
   Total Duration:       3m28.5s
   Started:              2026-01-14 15:30:00
   Ended:                2026-01-14 15:33:28

📈 Requests:
   Total Requests:       100,000
   Successful:           99,234 (99.2%)
   Failed:               766 (0.8%)

⚡ Performance:
   Average RPS:          479.62 requests/sec
   Target RPS:           500 requests/sec
   Concurrent Users:     200

⏲️  Latency:
   Average:              156ms
   Min:                  12ms
   Max:                  2847ms

📋 HTTP Status Codes:
   200: 99,234 (99.2%)
   429: 542 (0.5%)
   500: 224 (0.2%)

=================================================================================

💡 Performance Assessment:
   ✅ Excellent - System handled load very well
```

## What Makes This Tool Professional

1. **Realistic Traffic Patterns**: Variable think times simulate real user behavior
2. **Connection Pooling**: Efficient HTTP client with connection reuse
3. **Detailed Metrics**: Latency tracking, status codes, success rates
4. **Real-time Monitoring**: Progress updates during long tests
5. **Performance Assessment**: Automatic analysis of results
6. **Predefined Profiles**: Industry-standard test scenarios
7. **Safe Defaults**: Built-in protections and warnings

## Load Testing Best Practices

### Start Small and Scale Up
```bash
# 1. Test with minimal load first
./loadtest -url http://your-site.com -requests 100 -rps 10

# 2. Gradually increase
./loadtest -url http://your-site.com -profile normal

# 3. Move to peak scenarios
./loadtest -url http://your-site.com -profile peak

# 4. Finally, stress test
./loadtest -url http://your-site.com -profile stress
```

### Test Different Endpoints
```bash
# Homepage
./loadtest -url http://your-site.com -profile normal

# Product listing
./loadtest -url http://your-site.com/products -profile peak

# API endpoints
./loadtest -url http://your-site.com/api/search -profile flash-sale

# Checkout flow (most critical)
./loadtest -url http://your-site.com/checkout -profile stress
```

### Monitor Your Server
While running tests, monitor:
- CPU usage
- Memory consumption
- Database connections
- Response times
- Error logs

### When to Use Each Profile

- **Normal**: Daily capacity verification, smoke tests
- **Peak**: Preparing for expected high-traffic periods (evenings, weekends)
- **Flash Sale**: Testing promotional events, limited-time offers
- **Stress**: Finding system limits, capacity planning

## Use Cases (Authorized Only)

### ✅ Legal Use Cases
- Testing your own ecommerce platform
- Load testing staging/test environments
- Capacity planning before major sales events
- Performance regression testing in CI/CD
- Finding bottlenecks in your infrastructure
- Validating auto-scaling configurations

### ❌ Illegal - DO NOT DO
- Testing third-party websites without permission
- "Testing" competitors' sites
- Any unauthorized load testing
- Attacking or disrupting services
- Testing production systems without approval

## Performance Metrics Explained

### Success Rate
- **99.5%+**: Excellent - Production ready
- **95-99%**: Good - Minor issues to investigate
- **90-95%**: Fair - Significant issues found
- **<90%**: Poor - Major problems, not production ready

### Latency Guidelines (for Ecommerce)
- **<100ms**: Excellent user experience
- **100-300ms**: Good, acceptable
- **300-1000ms**: Slow, users will notice
- **>1000ms**: Poor, unacceptable for most cases

### RPS Achievement
If actual RPS is significantly lower than target:
- Server is bottlenecked
- Database queries too slow
- Network bandwidth limited
- Consider scaling horizontally

## Testing Locally

Set up a simple test server:

```go
// testserver.go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(50 * time.Millisecond) // Simulate processing
    fmt.Fprintf(w, "OK")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Test server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

Run it: `go run testserver.go`

Then test: `./loadtest -url http://localhost:8080 -requests 50`

## Learning Resources

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [HTTP Client Best Practices](https://go.dev/doc/tutorial/web-service-gin)
- [Rate Limiting in Go](https://gobyexample.com/rate-limiting)

## Contributing

Contributions that improve educational value or safety features are welcome. Pull requests that remove safety limitations will be rejected.

## License

GNU GENERAL PUBLIC LICENSE - See [LICENSE file](./LICENSE) for details

## Disclaimer

This tool is provided for educational purposes only. The authors are not responsible for any misuse or damage caused by this tool. Always ensure you have proper authorization before conducting any load testing.

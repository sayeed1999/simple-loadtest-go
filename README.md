# go-loadtest

A professional HTTP load testing tool written in Go for ecommerce and web applications. Includes predefined profiles for different traffic scenarios.

## ⚠️ LEGAL WARNING

**THIS TOOL IS FOR TESTING YOUR OWN INFRASTRUCTURE ONLY.**

- Using this tool against systems you don't own or operate is **ILLEGAL and a CYBERCRIME**
- It may result in **IMPRISONMENT or other serious punishment** if reported & tracked via IP
- You must have **explicit written authorization** before testing any system
- The authors are not responsible for misuse of this tool

## Educational Purpose

This tool is designed to teach:
- HTTP client programming in Go
- Concurrency patterns using goroutines and channels
- Rate limiting and throttling techniques
- Performance metrics collection
- Responsible load testing practices

## Features

✅ Built-in rate limiting (max 50 RPS)
✅ Request limits (max 1000 requests)
✅ Configurable concurrency (max 10 workers)
✅ Mandatory authorization confirmation
✅ Respect for server responses (429, 503)
✅ Think time between requests (realistic traffic)
✅ Proper timeout handling
✅ Clear warning messages

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-loadtest.git
cd go-loadtest

# Build
go build -o loadtest main.go
```

## Usage

### Basic Example (Testing localhost)

```bash
./loadtest -url http://localhost:8080 -requests 50 -rps 5
```

### All Options

```bash
./loadtest \
  -url http://localhost:8080/api/test \
  -requests 100 \
  -rps 10 \
  -concurrency 5 \
  -timeout 10 \
  -think-time 100
```

### Command Line Flags

| Flag | Description | Default | Max Limit |
|------|-------------|---------|-----------|
| `-url` | Target URL to test (required) | - | - |
| `-requests` | Total number of requests | 100 | 1000 |
| `-rps` | Requests per second | 10 | 50 |
| `-concurrency` | Number of concurrent workers | 5 | 10 |
| `-timeout` | Request timeout in seconds | 10 | - |
| `-think-time` | Delay between requests (ms) | 100 | min 100 |

## Example Output

```
==========================================================
⚠️  LEGAL WARNING - READ CAREFULLY
==========================================================
This tool is for testing YOUR OWN infrastructure only.
Using this against systems you don't own or operate is ILLEGAL.
You must have explicit written authorization before testing.
==========================================================

Do you have authorization to test this endpoint? (yes/no): yes

🚀 Starting load test...
Target: http://localhost:8080
Total Requests: 100
Concurrency: 5
Max RPS: 10
--------------------------------------------------

==========================================================
📊 LOAD TEST RESULTS
==========================================================
Total Requests:     100
Successful:         98 (98.0%)
Failed:             2 (2.0%)
Total Duration:     10.5s
Average RPS:        9.52
==========================================================
```

## Safety Features

1. **Hard Limits**: Tool enforces maximum values for requests, RPS, and concurrency
2. **Authorization Check**: Requires explicit confirmation before running
3. **Rate Limiting**: Built-in request rate limiting to prevent overwhelming targets
4. **Think Time**: Mandatory delays between requests to simulate realistic traffic
5. **Server Respect**: Monitors for 429 (rate limit) and 503 (unavailable) responses
6. **Timeouts**: Prevents hanging requests

## Use Cases (Authorized Only)

✅ Testing your own web applications during development
✅ Load testing your staging/test environments
✅ Capacity planning for your infrastructure
✅ Learning Go concurrency patterns
✅ Understanding HTTP client behavior

❌ Testing production systems without approval
❌ Testing third-party websites
❌ Attacking or disrupting services
❌ Any unauthorized testing

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

MIT License - See LICENSE file for details

## Disclaimer

This tool is provided for educational purposes only. The authors are not responsible for any misuse or damage caused by this tool. Always ensure you have proper authorization before conducting any load testing.

# gurl

[![CI](https://github.com/antlabs/gurl/actions/workflows/ci.yml/badge.svg)](https://github.com/antlabs/gurl/actions/workflows/ci.yml)
[![CodeQL](https://github.com/antlabs/gurl/actions/workflows/codeql.yml/badge.svg)](https://github.com/antlabs/gurl/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/antlabs/gurl)](https://goreportcard.com/report/github.com/antlabs/gurl)
[![License](https://img.shields.io/github/license/antlabs/gurl)](LICENSE)
[![Release](https://img.shields.io/github/v/release/antlabs/gurl)](https://github.com/antlabs/gurl/releases)

**gurl** = **g**o + **curl** + wr**k**

A modern, high-performance HTTP benchmarking tool written in Go, inspired by **wrk**, with native support for parsing **curl** commands.  

Turn your everyday `curl` into a scalable load test in seconds — no configuration needed. Just copy, paste, and go.

## Features

- 🚀 High-performance HTTP load testing
- 🔧 Parse curl commands with `--parse-curl` option
- 📁 **Multiple curl commands** from file with load distribution strategies
- 📊 Detailed statistics similar to wrk
- 📈 **Per-endpoint statistics** - TPS, latency, status codes for each endpoint
- 🎨 **Live Terminal UI** with real-time charts and statistics
- 🎭 **Mock HTTP Server** - Built-in test server for benchmarking validation
- 🌐 **RESTful API** - Submit and manage benchmark tasks via HTTP API
- ⚡ Async I/O for maximum performance
- 🎯 Configurable connections, threads, and duration
- 📉 Latency distribution analysis
- ⌨️ Interactive controls (press 'q' to stop early)
- 🔀 Load strategies: random, round-robin

## Installation

```bash
go install github.com/antlabs/gurl/cmd/gurl@latest
```

Or build from source:

```bash
git clone https://github.com/antlabs/gurl.git
cd gurl
go build -o gurl
```

## Usage

### Basic Usage

```bash
# Simple GET request
gurl -c 100 -d 30s http://example.com

# POST request with custom headers
gurl -c 50 -d 10s -X POST -H "Content-Type: application/json" -d '{"key":"value"}' http://api.example.com
```

### Parse Curl Commands

```bash
# Parse a curl command and use it for benchmarking
gurl --parse-curl "curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"test\"}' http://api.example.com/users" -c 100 -d 30s
```

### Options

- `-c, --connections`: Number of HTTP connections to keep open (default: 10)
- `-d, --duration`: Duration of test (default: 10s)
- `-t, --threads`: Number of threads to use (default: 2)
- `-R, --rate`: Work rate (requests/sec) 0=unlimited (default: 0)
- `--timeout`: Socket/request timeout (default: 30s)
- `--parse-curl`: Parse curl command and use it for benchmarking
- `--parse-curl-file`: Parse multiple curl commands from file (one per line)
- `--load-strategy`: Load distribution strategy: random, round-robin (default: random)
- `-X, --method`: HTTP method (default: GET)
- `-H, --header`: HTTP header to add to request
- `-d, --data`: HTTP request body
- `--content-type`: Content-Type header
- `-v, --verbose`: Verbose output
- `--latency`: Print latency statistics
- `--live-ui`: Enable live terminal UI with real-time stats (interactive mode)
- `--ui-theme`: UI color theme: dark, light, or auto (default: auto)
- `--use-nethttp`: Force use standard library net/http instead of pulse

## Examples

### Basic Load Test

```bash
gurl -c 100 -d 30s -t 4 http://example.com
```

Output:
```
Running 30s test @ http://example.com
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.34ms   5.67ms  123.45ms    75.23%
    Req/Sec     2.34k     0.12k    3.45k     89.12%
  280000 requests in 30.00s, 45.67MB read
Requests/sec:   9333.33
Transfer/sec:     1.52MB
```

### Load Test with Curl Command

```bash
gurl --parse-curl "curl -X POST -H 'Authorization: Bearer token123' -H 'Content-Type: application/json' -d '{\"user\":\"test\"}' https://api.example.com/login" -c 50 -d 10s --latency
```

### Rate Limited Test

```bash
# Limit to 1000 requests per second
gurl -c 10 -d 60s -R 1000 http://example.com
```

### Multiple Curl Commands

Test multiple endpoints simultaneously with different load distribution strategies:

```bash
# Create a file with multiple curl commands (one per line)
cat > endpoints.txt << EOF
curl https://api.example.com/users
curl -X POST https://api.example.com/orders -H "Content-Type: application/json" -d '{"item":"book"}'
curl https://api.example.com/products
curl https://api.example.com/search?q=test
EOF

# Random distribution (default) - randomly select from endpoints
gurl --parse-curl-file endpoints.txt -c 100 -d 60s --use-nethttp

# Round-robin distribution - evenly distribute across endpoints
gurl --parse-curl-file endpoints.txt --load-strategy round-robin -c 100 -d 60s --use-nethttp

# With rate limiting
gurl --parse-curl-file endpoints.txt --load-strategy random -c 50 -d 30s -R 1000 --use-nethttp
```

**File Format**:
- One curl command per line
- Empty lines are ignored
- Lines starting with `#` are treated as comments
- Supports all curl options (headers, methods, data, etc.)

**Load Strategies**:
- `random`: Randomly select an endpoint for each request (default)
- `round-robin`: Evenly distribute requests across all endpoints

**Per-Endpoint Statistics**:

When testing multiple endpoints, gurl automatically provides detailed statistics for each endpoint:

```
=== Per-Endpoint Statistics ===

[https://api.example.com/users]
  Requests:     250
  Requests/sec: 50.00
  Latency:      avg=120.50ms, min=45.20ms, max=350.80ms
  Status codes: [200] 248 (99.2%), [500] 2 (0.8%)
  Data:         125.5KB total, 25.1KB/sec

[https://api.example.com/orders]
  Requests:     245
  Errors:       5 (2.0%)
  Requests/sec: 49.00
  Latency:      avg=145.30ms, min=60.10ms, max=420.50ms
  Status codes: [200] 230 (93.9%), [400] 10 (4.1%)
  Data:         98.2KB total, 19.6KB/sec
```

Each endpoint shows:
- **Requests**: Total number of requests sent to this endpoint
- **Errors**: Number of failed requests (connection errors, timeouts)
- **Requests/sec (TPS)**: Throughput for this specific endpoint
- **Latency**: Average, minimum, and maximum response times
- **Status codes**: HTTP status code distribution with percentages
- **Data**: Total data transferred and transfer rate

### Live Terminal UI

Enable real-time interactive UI with live statistics:

```bash
# Basic live UI
gurl --live-ui -c 100 -d 60s --use-nethttp http://example.com

# Live UI with rate limiting
gurl --live-ui -c 1000 -d 300s -R 5000 --use-nethttp http://api.example.com

# Live UI with multiple endpoints
gurl --parse-curl-file endpoints.txt --live-ui -c 100 -d 60s --use-nethttp
```

The live UI displays:
- **Progress Bar**: Visual progress indicator with elapsed/total time
- **Real-time Stats**: Requests, slowest/fastest/average latency, requests per second
- **Status Code Distribution**: HTTP status codes with color coding (2xx=green, 4xx=yellow, 5xx=red)
- **Error Statistics**: Connection errors and error rate
- **Request Chart**: Bar chart showing requests per second over time
- **Response Time Histogram**: Latency distribution (p50, p75, p90, p95, p99)
- **Per-Endpoint Table** (multi-endpoint mode): Real-time statistics for each endpoint

**Multi-Endpoint Live UI**:

When testing multiple endpoints with `--parse-curl-file`, the Live UI automatically displays an additional table showing per-endpoint statistics:

```
┌─Per-Endpoint Statistics (live)────────────────────────────┐
│ Endpoint              │ Req/s │ Avg    │ Min   │ Max     │
├───────────────────────┼───────┼────────┼───────┼─────────┤
│ /api/fast             │ 1000  │ 10ms   │ 5ms   │ 50ms    │
│ /api/slow             │  950  │ 105ms  │ 95ms  │ 150ms   │
│ /api/error            │  957  │ 15ms   │ 10ms  │ 30ms    │
└───────────────────────┴───────┴────────┴───────┴─────────┘
```

Each row shows:
- Endpoint URL (truncated if too long)
- Requests per second (TPS)
- Average, minimum, and maximum latency
- Error count and percentage

**Interactive Controls**:
- Press `q` or `Ctrl+C` to stop the test early
- UI updates every second with live data

**Color Themes**:

gurl supports adaptive color themes for optimal visibility on different terminal backgrounds:

```bash
# Auto-detect terminal background (default)
gurl --live-ui -c 100 -d 60s http://example.com

# Force dark theme (for dark terminals)
gurl --live-ui --ui-theme dark -c 100 -d 60s http://example.com

# Force light theme (for light terminals)
gurl --live-ui --ui-theme light -c 100 -d 60s http://example.com

# Set theme via environment variable
export GURL_THEME=light
gurl --live-ui -c 100 -d 60s http://example.com
```

The system automatically detects your terminal's background color and adjusts the UI colors accordingly:
- **Dark Theme**: Optimized for black/dark backgrounds (white text, bright colors)
- **Light Theme V4.0 (Ultimate)**: Revolutionarily designed for white/light backgrounds:
  - **主色调**: Deep Blue (60%) - 深沉有力，统一全局
  - **强调色**: Magenta (30%) - 醒目突出，创造焦点
  - **功能色**: Red/Green (10%) - 经典语义
  - **彻底摒弃**: Cyan和Yellow（对比度不够）
  
The ultimate light theme uses only 4 colors with strong contrast ratios, creating a powerful, clean, and highly visible experience on white backgrounds. Based on the 60-30-10 golden ratio for optimal visual hierarchy.

**V4.0 Key Features**:
- 🔵 Blue as primary color (strong contrast)
- 💜 Magenta for emphasis (balanced usage)
- ✅ Only 4 colors (simple & powerful)
- ❌ No Cyan/Yellow (removed for clarity)

For more details, see:
- [Ultimate Light Theme V4.0](docs/ULTIMATE_LIGHT_THEME.md) ⭐ NEW
- [UI Theme Documentation](docs/UI_THEME.md)
- [Theme Evolution History](docs/THEME_EVOLUTION.md)

**Note**: Live UI currently requires `--use-nethttp` flag.

### Mock HTTP Server

Start a built-in mock HTTP server for testing and benchmarking:

```bash
# Start a simple echo server (returns request details)
gurl --mock-server --mock-port 8080

# Custom response with delay
gurl --mock-server --mock-port 8080 --mock-delay 100ms --mock-response '{"status":"ok"}'

# Custom status code
gurl --mock-server --mock-port 8080 --mock-status 500 --mock-response '{"error":"server error"}'

# Use configuration file for multiple routes
gurl --mock-server --mock-config examples/mock-server.yaml
```

**Mock Server Configuration File** (`mock-server.yaml`):

```yaml
port: 8080

routes:
  # Echo endpoint - returns request details
  - path: /echo
    method: GET
    echo: true

  # Fast endpoint - no delay
  - path: /fast
    method: GET
    status_code: 200
    response: '{"message": "Fast response"}'

  # Slow endpoint - 100ms delay
  - path: /slow
    method: GET
    status_code: 200
    delay: 100ms
    response: '{"message": "Slow response"}'

  # Error endpoint
  - path: /error
    method: GET
    status_code: 500
    response: '{"error": "Internal server error"}'

  # POST endpoint with echo
  - path: /api/users
    method: POST
    status_code: 201
    echo: true
```

**Testing the Mock Server**:

```bash
# Terminal 1: Start mock server
gurl --mock-server --mock-config examples/mock-server.yaml

# Terminal 2: Run benchmark against it
gurl -c 100 -d 30s --use-nethttp http://localhost:8080/fast
gurl -c 100 -d 30s --use-nethttp http://localhost:8080/slow

# Test multiple endpoints
cat > endpoints.txt << EOF
curl http://localhost:8080/fast
curl http://localhost:8080/slow
curl http://localhost:8080/error
EOF

gurl --parse-curl-file endpoints.txt -c 50 -d 30s --use-nethttp
```

**Mock Server Features**:
- **Echo mode**: Returns request details (method, headers, body)
- **Custom responses**: Define JSON or text responses
- **Configurable delays**: Simulate slow endpoints
- **Status codes**: Test error handling
- **Multiple routes**: Define different endpoints with different behaviors
- **Request logging**: See all incoming requests in real-time

## RESTful API Server

gurl provides a RESTful API server that allows you to submit and manage benchmark tasks via HTTP requests.

### Starting the API Server

```bash
# Start API server on default port 8080
gurl --api-server

# Start API server on custom port
gurl --api-server --api-port 9090
```

### API Endpoints

- **POST** `/api/v1/benchmark` - Submit a benchmark task
- **POST** `/api/v1/batch` - Submit a batch test task
- **GET** `/api/v1/status/:id` - Get task status
- **GET** `/api/v1/results/:id` - Get task results
- **GET** `/health` - Health check endpoint

### Quick Example

```bash
# Start API server
gurl --api-server --api-port 8080

# In another terminal, submit a benchmark task
curl -X POST http://localhost:8080/api/v1/benchmark \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/users",
    "connections": 100,
    "duration": "30s",
    "threads": 4,
    "use_nethttp": true
  }'

# Response:
# {
#   "task_id": "task_1234567890123456789",
#   "status": "accepted",
#   "message": "Benchmark task created and started"
# }

# Check task status
curl http://localhost:8080/api/v1/status/task_1234567890123456789

# Get results when completed
curl http://localhost:8080/api/v1/results/task_1234567890123456789
```

For detailed API documentation, see [API.md](docs/API.md).

## Batch Testing with Configuration Files

gurl supports batch testing through YAML or JSON configuration files, allowing you to run multiple tests with different parameters in a single command.

### Basic Batch Testing

```bash
# Run batch tests from YAML configuration
gurl --batch-config examples/batch-config.yaml

# Run batch tests from JSON configuration
gurl --batch-config examples/batch-config.json

# Run tests sequentially instead of concurrently
gurl --batch-config batch-tests.yaml --batch-sequential

# Limit concurrent batch tests (default: 3)
gurl --batch-config batch-tests.yaml --batch-concurrency 5

# Generate different report formats
gurl --batch-config batch-tests.yaml --batch-report csv
gurl --batch-config batch-tests.yaml --batch-report json
```

### Configuration File Format

#### YAML Format (`batch-config.yaml`)

```yaml
version: "1.0"
tests:
  - name: "用户登录API"
    curl: 'curl -X POST https://api.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 100
    duration: "30s"
    threads: 4
    
  - name: "获取用户信息"
    curl: 'curl -H "Authorization: Bearer token123" https://api.example.com/user/profile'
    connections: 50
    duration: "60s"
    threads: 2
    
  - name: "创建订单API"
    curl: 'curl -X POST https://api.example.com/orders -H "Content-Type: application/json" -d "{\"product_id\":1,\"quantity\":2}"'
    connections: 80
    duration: "45s"
    rate: 100
    timeout: "10s"
    verbose: true
```

#### JSON Format (`batch-config.json`)

```json
{
  "version": "1.0",
  "tests": [
    {
      "name": "API健康检查",
      "curl": "curl https://api.example.com/health",
      "connections": 10,
      "duration": "10s"
    },
    {
      "name": "用户注册接口",
      "curl": "curl -X POST https://api.example.com/register -H \"Content-Type: application/json\" -d \"{\\\"email\\\":\\\"test@example.com\\\",\\\"password\\\":\\\"password123\\\"}\"",
      "connections": 50,
      "duration": "30s",
      "threads": 3,
      "rate": 200
    }
  ]
}
```

### Configuration Parameters

Each test in the batch configuration supports the following parameters:

| Parameter | Type | Description | Default |
|-----------|------|-------------|----------|
| `name` | string | Test name (required) | - |
| `curl` | string | Curl command to parse (required) | - |
| `connections` | int | Number of HTTP connections | 10 |
| `duration` | string | Test duration (e.g., "30s", "5m") | 10s |
| `threads` | int | Number of threads | 2 |
| `rate` | int | Requests per second limit (0=unlimited) | 0 |
| `timeout` | string | Request timeout (e.g., "5s") | 30s |
| `verbose` | bool | Enable verbose output | false |
| `use_nethttp` | bool | Force use standard net/http | false |

### Batch Testing Options

| Option | Description | Default |
|--------|-------------|----------|
| `--batch-config` | Path to batch configuration file (YAML/JSON) | - |
| `--batch-concurrency` | Maximum concurrent batch tests | 3 |
| `--batch-sequential` | Run tests sequentially instead of concurrently | false |
| `--batch-report` | Report format: text, csv, json | text |

### Batch Test Output

#### Text Report (Default)
```
=== Batch Test Report ===

Total Tests: 3
Success Rate: 100.00%
Total Time: 1m30s
Start Time: 2024-01-15 10:30:00
End Time: 2024-01-15 10:31:30

=== Test Results ===

1. 用户登录API
   Duration: 30.2s
   Status: SUCCESS
   Requests: 15000
   RPS: 496.67
   Avg Latency: 201ms

2. 获取用户信息
   Duration: 60.1s
   Status: SUCCESS
   Requests: 30000
   RPS: 499.17
   Avg Latency: 100ms

=== Performance Summary ===

Total Requests: 45000
Combined RPS: 995.84
Latency Stats:
  Average: 150ms
  Median:  120ms
  Min:     50ms
  Max:     500ms
```

#### CSV Report
```bash
gurl --batch-config batch-tests.yaml --batch-report csv > results.csv
```

#### JSON Report
```bash
gurl --batch-config batch-tests.yaml --batch-report json > results.json
```

## URL Template Variables

gurl supports dynamic URL template variables that allow you to generate different values for each request, making it perfect for realistic load testing scenarios.

### Built-in Template Functions

| Function | Description | Usage | Example |
|----------|-------------|-------|----------|
| `random` | Random number in range | `{{random:min-max}}` | `{{random:1-1000}}` |
| `uuid` | Generate UUID | `{{uuid}}` | `550e8400-e29b-41d4-a716-446655440000` |
| `timestamp` | Current timestamp | `{{timestamp:format}}` | `{{timestamp:unix}}` |
| `now` | Alias for timestamp | `{{now:format}}` | `{{now:rfc3339}}` |
| `sequence` | Incrementing numbers | `{{sequence:start}}` | `{{sequence:1}}` |
| `choice` | Random selection | `{{choice:a,b,c}}` | `{{choice:GET,POST,PUT}}` |

### Template Formats

#### Timestamp Formats
- `unix` - Unix timestamp (default): `1640995200`
- `unix_ms` - Unix timestamp in milliseconds: `1640995200000`
- `rfc3339` - RFC3339 format: `2022-01-01T00:00:00Z`
- `iso8601` - ISO8601 format
- `date` - Date only: `2022-01-01`
- `time` - Time only: `15:04:05`

### Basic Template Usage

#### Simple Random User ID
```bash
# Test with random user IDs from 1 to 1000
gurl -c 50 -d 30s 'https://api.example.com/users/{{random:1-1000}}'
```

#### UUID Session Testing
```bash
# Each request gets a unique session ID
gurl -c 20 -d 60s 'https://api.example.com/data?session={{uuid}}'
```

#### Timestamp-based Requests
```bash
# Include current timestamp in requests
gurl -c 10 -d 30s 'https://api.example.com/events?timestamp={{timestamp:unix}}'
```

#### Sequential Page Testing
```bash
# Test pagination with incrementing page numbers
gurl -c 5 -d 60s 'https://api.example.com/items?page={{sequence:1}}&limit=20'
```

### Custom Variable Definitions

Define your own variables using the `--var` option:

```bash
# Define custom variables
gurl --var user_id=random:1-10000 \
     --var method=choice:GET,POST,PUT \
     --var session=uuid \
     -c 30 -d 45s \
     'https://api.example.com/{{method}}/users/{{user_id}}?session={{session}}'
```

### Advanced Template Examples

#### E-commerce API Simulation
```bash
gurl --var user_id=random:1-10000 \
     --var product_id=random:100-999 \
     --var quantity=choice:1,2,3,4,5 \
     --var payment=choice:credit_card,paypal,apple_pay \
     -c 50 -d 60s \
     --parse-curl 'curl -X POST https://shop.example.com/api/orders \
                   -H "Content-Type: application/json" \
                   -d "{\"user_id\":{{user_id}},\"product_id\":{{product_id}},\"quantity\":{{quantity}},\"payment_method\":\"{{payment}}\",\"timestamp\":\"{{timestamp:rfc3339}}\"}"}'
```

#### Multi-endpoint Testing
```bash
gurl --var endpoint=choice:users,orders,products,reviews \
     --var id=random:1-1000 \
     --var action=choice:view,edit,delete \
     -c 25 -d 30s \
     'https://api.example.com/{{endpoint}}/{{id}}/{{action}}'
```

### Template Variables in Batch Configuration

Template variables work seamlessly with batch configuration files:

```yaml
version: "1.0"
tests:
  - name: "Dynamic User API Test"
    curl: 'curl https://api.example.com/users/{{random:1-10000}}'
    connections: 50
    duration: "30s"
    
  - name: "Session-based Requests"
    curl: 'curl -H "X-Session-ID: {{uuid}}" -H "X-Timestamp: {{timestamp:unix}}" https://api.example.com/data'
    connections: 30
    duration: "45s"
    
  - name: "Sequential Pagination"
    curl: 'curl "https://api.example.com/posts?page={{sequence:1}}&limit=10"'
    connections: 10
    duration: "60s"
```

Run with custom variables:
```bash
gurl --batch-config template-test.yaml \
     --var api_key=uuid \
     --var region=choice:us-east,us-west,eu-central
```

### Template Variable Help

Get detailed help and examples for template variables:

```bash
# Show all available template functions and examples
gurl --help-templates
```

### Template Variable Best Practices

1. **Realistic Data Generation**: Use appropriate ranges and choices that match your real-world data
2. **Performance Considerations**: Template parsing adds minimal overhead but consider it for very high-rate tests
3. **Debugging**: Use `-v` (verbose) flag to see original and processed URLs/commands
4. **Variable Reuse**: Define commonly used variables once with `--var` instead of inline functions
5. **Batch Testing**: Combine template variables with batch configuration for comprehensive test suites

### Advanced Batch Testing Examples

#### API Load Testing Suite
```yaml
version: "1.0"
tests:
  - name: "Health Check"
    curl: 'curl https://api.example.com/health'
    connections: 5
    duration: "10s"
    
  - name: "Authentication"
    curl: 'curl -X POST https://api.example.com/auth -d "username=test&password=123"'
    connections: 20
    duration: "30s"
    
  - name: "Data Retrieval"
    curl: 'curl -H "Authorization: Bearer token" https://api.example.com/data'
    connections: 100
    duration: "60s"
    rate: 500
    
  - name: "Heavy Processing"
    curl: 'curl -X POST https://api.example.com/process -d @large-payload.json'
    connections: 10
    duration: "120s"
    timeout: "30s"
```

#### E-commerce API Testing
```yaml
version: "1.0"
tests:
  - name: "Product Search"
    curl: 'curl "https://shop.example.com/api/search?q=laptop"'
    connections: 50
    duration: "45s"
    
  - name: "Add to Cart"
    curl: 'curl -X POST https://shop.example.com/api/cart -H "Content-Type: application/json" -d "{\"product_id\":123,\"quantity\":1}"'
    connections: 30
    duration: "30s"
    
  - name: "Checkout Process"
    curl: 'curl -X POST https://shop.example.com/api/checkout -H "Authorization: Bearer token" -d @checkout-data.json'
    connections: 20
    duration: "60s"
    rate: 50
```

## Output Format

The output format is similar to wrk:

- **Thread Stats**: Average, standard deviation, and maximum latency
- **Req/Sec**: Requests per second statistics
- **Total Summary**: Total requests, duration, and data transferred
- **Error Summary**: Connection, read, write, and timeout errors (if any)
- **Status Code Distribution**: HTTP status code breakdown
- **Latency Distribution**: Percentile breakdown (with --latency flag)

## Performance

gurl is designed for high performance:

- Async I/O using custom pulse library
- Efficient connection pooling
- Minimal memory allocation during testing
- Multi-threaded architecture

## License

MIT License - see LICENSE file for details.

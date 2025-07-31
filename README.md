# gurl

**gurl** = **m**odern + c**url** + wr**k**

A modern HTTP benchmarking tool inspired by wrk, with the ability to parse curl commands for load testing.

## Features

- 🚀 High-performance HTTP load testing
- 🔧 Parse curl commands with `--parse-curl` option
- 📊 Detailed statistics similar to wrk
- ⚡ Async I/O for maximum performance
- 🎯 Configurable connections, threads, and duration
- 📈 Latency distribution analysis

## Installation

```bash
go install github.com/antlabs/gurl@latest
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
- `-X, --method`: HTTP method (default: GET)
- `-H, --header`: HTTP header to add to request
- `-d, --data`: HTTP request body
- `--content-type`: Content-Type header
- `-v, --verbose`: Verbose output
- `--latency`: Print latency statistics

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

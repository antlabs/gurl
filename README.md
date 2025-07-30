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

# JSONDiff Feature for Batch Testing

## Overview

The JSONDiff feature allows you to compare JSON responses between different API versions (A/B testing) during batch performance testing. This is useful for:

- Verifying that new API versions return the same data as old versions
- Detecting unintended changes in API responses
- Validating data consistency across different environments

## Basic Usage

### Simple Example: Compare Entire Response

```yaml
version: "1.0"
tests:
  - name: "A版本-接口"
    curl: 'curl https://api-v1.example.com/endpoint'
    connections: 100
    duration: "30s"
    
  - name: "B版本-接口"
    curl: 'curl https://api-v2.example.com/endpoint'
    connections: 100
    duration: "30s"

jsondiff:
  pairs:
    - base: "A版本-接口"
      target: "B版本-接口"
      compare_field: "data"  # Only compare the 'data' field
```

### Example: Compare code-msg-data Structure

```yaml
version: "1.0"
tests:
  - name: "A版本-用户查询"
    curl: 'curl https://api-v1.example.com/user/123'
    connections: 50
    duration: "30s"
    
  - name: "B版本-用户查询"
    curl: 'curl https://api-v2.example.com/user/123'
    connections: 50
    duration: "30s"

jsondiff:
  pairs:
    - name: "用户查询对比"
      base: "A版本-用户查询"
      target: "B版本-用户查询"
      compare_field: "data"  # Only compare the 'data' field
      ignore_fields:
        - "timestamp"
        - "request_id"
```

## Configuration Options

### JSONDiff Pair Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | No | Name of the comparison (auto-generated if not provided) |
| `base` | string | Yes | Name of the base test (from tests array) |
| `target` | string | Yes | Name of the target test (from tests array) |
| `compare_field` | string | No | JSON path to compare (e.g., "data", "result.user"). Applied to both base and target. If not specified, compares entire response |
| `base_field` | string | No | JSON path for base test only. Use when old and new versions have different field names |
| `target_field` | string | No | JSON path for target test only. Use when old and new versions have different field names |
| `ignore_fields` | []string | No | List of field paths to ignore during comparison |

**Field Priority**:
- If `base_field` and `target_field` are specified, they are used (supports different field paths)
- If only `compare_field` is specified, it applies to both base and target
- If none are specified, compares entire response

### Field Path Syntax

Field paths use dot notation for nested objects:

- `"data"` - Compare only the data field
- `"data.user"` - Compare only the user object inside data
- `"result.items.0"` - Compare the first item in the items array
- `"timestamp"` - Ignore any field named "timestamp" at any level

## Use Cases

### 1. Compare Only Data Field

Most common use case - ignore code/msg and only compare data:

```yaml
jsondiff:
  pairs:
    - base: "A版本-接口"
      target: "B版本-接口"
      compare_field: "data"
```

### 2. Compare Entire Response with Ignored Fields

Compare everything but ignore dynamic fields:

```yaml
jsondiff:
  pairs:
    - base: "A版本-接口"
      target: "B版本-接口"
      ignore_fields:
        - "timestamp"
        - "request_id"
        - "trace_id"
        - "server_time"
```

### 3. Compare Nested Fields

Compare specific nested structures:

```yaml
jsondiff:
  pairs:
    - base: "A版本-接口"
      target: "B版本-接口"
      compare_field: "data.user.profile"
      ignore_fields:
        - "last_login_time"
        - "update_time"
```

### 4. Multiple Comparisons

Compare multiple API pairs in one batch:

```yaml
version: "1.0"
tests:
  - name: "A-登录"
    curl: 'curl -X POST https://api-v1.example.com/login -d "..."'
    connections: 100
    duration: "30s"
    
  - name: "B-登录"
    curl: 'curl -X POST https://api-v2.example.com/login -d "..."'
    connections: 100
    duration: "30s"
    
  - name: "A-用户信息"
    curl: 'curl https://api-v1.example.com/user/info'
    connections: 50
    duration: "30s"
    
  - name: "B-用户信息"
    curl: 'curl https://api-v2.example.com/user/info'
    connections: 50
    duration: "30s"

jsondiff:
  pairs:
    - name: "登录接口对比"
      base: "A-登录"
      target: "B-登录"
      compare_field: "data"
      
    - name: "用户信息对比"
      base: "A-用户信息"
      target: "B-用户信息"
      compare_field: "data"
      ignore_fields:
        - "last_login_time"
```

### 5. Compare Different Field Paths (New Feature)

When old and new API versions use different field names:

```yaml
version: "1.0"
tests:
  - name: "老版本-获取地区156"
    curl: 'curl https://api-v1.example.com/region/156'
    connections: 1
    threads: 1
    requests: 1
    
  - name: "新版本-获取地区156"
    curl: 'curl https://api-v2.example.com/region/156'
    connections: 1
    threads: 1
    requests: 1

jsondiff:
  pairs:
    - name: "获取地址接口A/B对比"
      base: "老版本-获取地区156"
      target: "新版本-获取地区156"
      base_field: "result"      # Old version uses "result"
      target_field: "data"      # New version uses "data"
      ignore_fields:
        - "timestamp"
        - "request_id"
```

**Old version response:**
```json
{
  "code": 200,
  "msg": "success",
  "result": {
    "region": "Beijing",
    "code": "156"
  }
}
```

**New version response:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "region": "Beijing",
    "code": "156"
  }
}
```

With `base_field: "result"` and `target_field: "data"`, the comparison will pass even though the field names are different, as long as the content is identical.

## Running Tests

```bash
# Run batch test with jsondiff
gurl --batch-config examples/batch-jsondiff.yaml

# Run with verbose output to see comparison details
gurl --batch-config examples/batch-jsondiff.yaml --verbose

# Generate JSON report
gurl --batch-config examples/batch-jsondiff.yaml --batch-report json > results.json
```

## Output Example

```
=== Batch Test Report ===

Total Tests: 2
Success Rate: 100.00%
Total Time: 10.5s
Start Time: 2024-12-27 19:30:00
End Time: 2024-12-27 19:30:10

=== Test Results ===

1. A版本-接口
   Duration: 5.2s
   Status: SUCCESS
   Requests: 500
   RPS: 96.15
   Avg Latency: 104ms

2. B版本-接口
   Duration: 5.3s
   Status: SUCCESS
   Requests: 530
   RPS: 100.00
   Avg Latency: 100ms

=== JSON Diff Results ===

Comparison: A版本-接口 vs B版本-接口
  Base: A版本-接口
  Target: B版本-接口
  Status: PASSED
  Message: field 'data' is identical in both responses

JSONDiff Summary: 1 passed, 0 failed
```

## Example Files

- `batch-jsondiff-simple.yaml` - Simple comparison example
- `batch-jsondiff.yaml` - Full featured example with multiple comparisons
- `batch-jsondiff-codemsg.yaml` - Example for code-msg-data structure

## Tips

1. **Start with compare_field**: Use `compare_field: "data"` to focus on the actual data and ignore metadata
2. **Ignore dynamic fields**: Always ignore timestamps, request IDs, and other dynamic fields
3. **Test incrementally**: Start with one comparison pair, verify it works, then add more
4. **Use verbose mode**: Run with `--verbose` to see detailed comparison results during testing
5. **Capture sample early**: The tool captures a sample response before running the benchmark, so ensure your endpoint is accessible

## Limitations

- Comparisons are based on a single sample response captured before the benchmark
- Array order matters - elements must be in the same order
- Numeric precision is exact - no tolerance for floating point differences (yet)
- Large responses may impact memory usage

## Future Enhancements

Planned features:
- Array comparison with order-independence
- Numeric tolerance for floating point comparisons
- Regular expression matching for field values
- Custom comparison functions
- Sampling multiple responses for comparison

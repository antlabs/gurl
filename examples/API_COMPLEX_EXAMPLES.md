# gurl API 复杂场景使用示例

本文档展示如何使用 gurl RESTful API 进行各种复杂场景的压测，包括自定义 headers、body、query 参数等。

## 目录

- [基础概念](#基础概念)
- [示例 1: POST + JSON Body](#示例-1-post--json-body)
- [示例 2: GET + Query 参数](#示例-2-get--query-参数)
- [示例 3: PUT + Form Data](#示例-3-put--form-data)
- [示例 4: 使用完整 curl 命令](#示例-4-使用完整-curl-命令)
- [示例 5: 批量测试](#示例-5-批量测试)
- [参数说明](#参数说明)

## 基础概念

gurl API 支持两种方式提交压测任务：

1. **结构化方式**: 分别指定 URL、method、headers、body 等参数
2. **curl 命令方式**: 直接提供完整的 curl 命令字符串

## 示例 1: POST + JSON Body

测试场景：用户注册 API，包含 JSON 格式的请求体和多个自定义 headers。

### 使用 Shell 脚本

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/users",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json",
      "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
      "X-Request-ID": "req-001",
      "User-Agent": "MyApp/1.0"
    },
    "body": "{\"username\":\"john\",\"email\":\"john@example.com\",\"password\":\"secret123\"}",
    "content_type": "application/json",
    "connections": 50,
    "duration": "30s",
    "threads": 10,
    "rate": 200,
    "timeout": "5s"
  }'
```

### 使用 JSON 配置文件

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-post-json.json
```

### Python 示例

```python
import requests
import json

api_url = "http://localhost:3434/api/v1/benchmark"

config = {
    "url": "https://api.example.com/users",
    "method": "POST",
    "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer YOUR_TOKEN",
        "X-Request-ID": "req-001"
    },
    "body": json.dumps({
        "username": "john",
        "email": "john@example.com",
        "age": 30
    }),
    "content_type": "application/json",
    "connections": 50,
    "duration": "30s",
    "threads": 10,
    "rate": 200
}

response = requests.post(api_url, json=config)
task = response.json()
print(f"Task ID: {task['task_id']}")

# 等待任务完成
import time
time.sleep(35)

# 获取结果
result_url = f"http://localhost:3434/api/v1/results/{task['task_id']}"
result = requests.get(result_url).json()
print(json.dumps(result['results'], indent=2))
```

## 示例 2: GET + Query 参数

测试场景：商品搜索 API，包含多个 query 参数和认证 headers。

### URL 中包含 Query 参数

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/products?page=1&limit=50&sort=price&order=asc&category=electronics",
    "method": "GET",
    "headers": {
      "Accept": "application/json",
      "X-API-Key": "your-api-key",
      "Cache-Control": "no-cache"
    },
    "connections": 30,
    "duration": "20s",
    "threads": 5,
    "rate": 100
  }'
```

### 复杂的 Query 参数（数组和嵌套）

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/search?q=laptop&filter[price][min]=500&filter[price][max]=2000&filter[brands][]=Dell&filter[brands][]=HP&fields=id,name,price",
    "method": "GET",
    "headers": {
      "Accept": "application/json",
      "Accept-Language": "zh-CN",
      "X-Session-Token": "session_xyz"
    },
    "connections": 20,
    "duration": "15s",
    "threads": 4
  }'
```

## 示例 3: PUT + Form Data

测试场景：更新产品信息，使用 form-urlencoded 格式。

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/products/12345",
    "method": "PUT",
    "headers": {
      "Content-Type": "application/x-www-form-urlencoded",
      "Authorization": "Basic dXNlcjpwYXNzd29yZA==",
      "X-CSRF-Token": "csrf-token-here"
    },
    "body": "name=Updated+Product&price=99.99&stock=100&category=electronics&status=active",
    "content_type": "application/x-www-form-urlencoded",
    "connections": 20,
    "duration": "15s",
    "threads": 4,
    "rate": 50
  }'
```

## 示例 4: 使用完整 curl 命令

当你已经有一个现成的 curl 命令时，可以直接使用。

### 从浏览器开发者工具复制的 curl 命令

```bash
# 1. 在浏览器中打开开发者工具 (F12)
# 2. 在 Network 标签页找到请求
# 3. 右键 -> Copy -> Copy as cURL
# 4. 将复制的命令传递给 gurl API

curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "curl": "curl \"https://api.example.com/orders\" -X POST -H \"accept: application/json\" -H \"authorization: Bearer YOUR_TOKEN\" -H \"content-type: application/json\" --data-raw \"{\\\"items\\\":[{\\\"id\\\":1,\\\"qty\\\":2}],\\\"total\\\":199.98}\"",
    "connections": 50,
    "duration": "30s",
    "threads": 10,
    "rate": 200
  }'
```

### 复杂的 curl 命令示例

```bash
CURL_CMD='curl -X POST "https://api.example.com/orders?source=web&campaign=summer2024" \
  -H "Accept: application/json, text/plain, */*" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" \
  -H "Content-Type: application/json; charset=utf-8" \
  -H "Origin: https://www.example.com" \
  -H "Referer: https://www.example.com/checkout" \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)" \
  -H "X-Requested-With: XMLHttpRequest" \
  -H "X-Client-Version: 2.5.0" \
  -H "X-Session-ID: sess_abc123" \
  --data-raw '"'"'{"order_id":"ORD-2024-001","items":[{"product_id":"PROD-001","quantity":1,"price":1299.99}],"payment":{"method":"credit_card","amount":1299.99},"shipping":{"address":"123 Main St","city":"SF","state":"CA","zip":"94102"},"metadata":{"source":"web","campaign":"summer2024"}}'"'"

curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d "{
    \"curl\": $(echo "$CURL_CMD" | jq -Rs .),
    \"connections\": 100,
    \"duration\": \"60s\",
    \"threads\": 20,
    \"rate\": 500,
    \"timeout\": \"10s\"
  }"
```

## 示例 5: 批量测试

测试多个不同的 API 端点。

```bash
curl -X POST "http://localhost:3434/api/v1/batch" \
  -H "Content-Type: application/json" \
  -d '{
    "tests": [
      {
        "name": "用户登录",
        "url": "https://api.example.com/auth/login",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{\"username\":\"test\",\"password\":\"test123\"}",
        "connections": 20,
        "duration": "10s",
        "threads": 4
      },
      {
        "name": "获取用户信息",
        "curl": "curl -H \"Authorization: Bearer TOKEN\" https://api.example.com/users/me",
        "connections": 30,
        "duration": "10s",
        "threads": 5
      },
      {
        "name": "搜索商品",
        "url": "https://api.example.com/products?q=laptop&limit=20",
        "method": "GET",
        "headers": {
          "Accept": "application/json"
        },
        "connections": 40,
        "duration": "10s",
        "threads": 8
      }
    ]
  }'
```

## 参数说明

### 通用参数

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| `url` | string | 是* | 目标 URL | - |
| `curl` | string | 是* | curl 命令（二选一） | - |
| `method` | string | 否 | HTTP 方法 | GET |
| `headers` | object | 否 | 自定义 headers | {} |
| `body` | string | 否 | 请求体内容 | - |
| `content_type` | string | 否 | Content-Type | - |
| `connections` | int | 否 | 并发连接数 | 10 |
| `duration` | string | 否 | 持续时间 | 10s |
| `threads` | int | 否 | 线程数 | 2 |
| `rate` | int | 否 | 每秒请求数限制 | 0(不限) |
| `timeout` | string | 否 | 请求超时时间 | 30s |
| `use_nethttp` | bool | 否 | 使用标准库 | false |
| `extra` | object | 否 | 额外的元数据 | {} |

\* `url` 和 `curl` 必须提供其中一个

### Headers 示例

常用的 headers：

```json
{
  "headers": {
    // 认证相关
    "Authorization": "Bearer YOUR_TOKEN",
    "X-API-Key": "your-api-key",
    "Cookie": "session_id=abc123",
    
    // 内容类型
    "Content-Type": "application/json",
    "Accept": "application/json",
    
    // 语言和编码
    "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
    "Accept-Encoding": "gzip, deflate, br",
    
    // 请求来源
    "Origin": "https://www.example.com",
    "Referer": "https://www.example.com/page",
    
    // 自定义头
    "X-Request-ID": "req-12345",
    "X-Client-Version": "1.0.0",
    "X-Device-ID": "device-xyz",
    
    // 缓存控制
    "Cache-Control": "no-cache",
    "If-None-Match": "etag-value"
  }
}
```

### Body 格式示例

#### JSON 格式

```json
{
  "body": "{\"key1\":\"value1\",\"key2\":{\"nested\":\"value\"},\"array\":[1,2,3]}",
  "content_type": "application/json"
}
```

#### Form Data 格式

```json
{
  "body": "key1=value1&key2=value2&key3=value+with+spaces",
  "content_type": "application/x-www-form-urlencoded"
}
```

#### XML 格式

```json
{
  "body": "<?xml version=\"1.0\"?><root><item>value</item></root>",
  "content_type": "application/xml"
}
```

## 获取结果

### 检查任务状态

```bash
curl "http://localhost:3434/api/v1/status/{task_id}"
```

响应：

```json
{
  "id": "task_xxx",
  "status": "completed",
  "created_at": "2024-01-01T00:00:00Z",
  "started_at": "2024-01-01T00:00:01Z",
  "completed_at": "2024-01-01T00:00:31Z"
}
```

### 获取详细结果

```bash
curl "http://localhost:3434/api/v1/results/{task_id}"
```

响应：

```json
{
  "id": "task_xxx",
  "status": "completed",
  "results": {
    "total_requests": 5000,
    "total_errors": 50,
    "duration": "30.00s",
    "average_latency": "125.45ms",
    "min_latency": "45.12ms",
    "max_latency": "890.23ms",
    "latency_stddev": "67.89ms",
    "requests_per_sec": 166.67,
    "status_codes": {
      "200": 4950,
      "500": 50
    },
    "latency_percentiles": {
      "p50": "120.34ms",
      "p75": "145.67ms",
      "p90": "178.90ms",
      "p95": "234.56ms",
      "p99": "567.89ms"
    },
    "total_bytes": 12500000
  }
}
```

## 完整工作流示例

### Shell 脚本

参考 `api-complex-example.sh` 获取完整的自动化测试脚本。

```bash
chmod +x examples/api-complex-example.sh
./examples/api-complex-example.sh
```

### 使用配置文件

```bash
# 1. 提交任务
TASK_ID=$(curl -s -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-post-json.json | jq -r '.task_id')

echo "Task ID: $TASK_ID"

# 2. 等待完成
sleep 35

# 3. 获取结果
curl -s "http://localhost:3434/api/v1/results/$TASK_ID" | jq .
```

## 注意事项

1. **请求体转义**: 在 JSON 中传递 JSON 字符串时，需要正确转义引号
2. **URL 编码**: Query 参数中的特殊字符需要进行 URL 编码
3. **Headers 大小写**: HTTP headers 是不区分大小写的，但建议使用标准的大小写格式
4. **Rate 限制**: 设置 `rate` 参数可以控制每秒请求数，避免压垮服务器
5. **超时设置**: 根据实际 API 的响应时间设置合理的 `timeout` 值

## 故障排查

### 问题 1: 请求体格式错误

**错误**: `Invalid request body`

**解决**: 确保 JSON 格式正确，特别是嵌套的 JSON 字符串需要转义

### 问题 2: Headers 未生效

**检查**: 使用 httpbin.org 的 `/headers` 端点验证 headers 是否正确发送

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://httpbin.org/headers",
    "headers": {"X-Custom": "test"},
    "connections": 1,
    "duration": "1s"
  }'
```

### 问题 3: 结果为空

**解决**: 等待足够的时间让压测完成，建议等待时间 = duration + 5秒

## 更多示例

查看 `examples/` 目录获取更多示例配置文件：

- `api-post-json.json` - POST 请求 + JSON
- `api-get-query.json` - GET 请求 + Query 参数
- `api-put-form.json` - PUT 请求 + Form Data
- `api-curl-command.json` - 完整 curl 命令

## 相关文档

- [API 文档](../docs/API.md)
- [命令行使用指南](../README.md)


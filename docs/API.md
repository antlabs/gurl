# gurl RESTful API 文档

gurl 提供了 RESTful API 接口，允许通过 HTTP 请求提交和管理压测任务。

## 启动 API 服务器

```bash
# 使用默认端口 8080
gurl --api-server

# 指定端口
gurl --api-server --api-port 9090
```

启动后，API 服务器将在指定端口上监听请求。

## API 端点

### 1. 提交压测任务

**POST** `/api/v1/benchmark`

提交一个新的压测任务。任务将异步执行。

#### 请求体示例

```json
{
  "url": "https://api.example.com/users",
  "connections": 100,
  "duration": "30s",
  "threads": 4,
  "rate": 0,
  "timeout": "30s",
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token123"
  },
  "use_nethttp": true
}
```

#### 使用 curl 命令

```json
{
  "curl": "curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"test\"}' https://api.example.com/users",
  "connections": 100,
  "duration": "30s",
  "threads": 4,
  "use_nethttp": true
}
```

#### 请求参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `url` | string | 是* | - | 目标 URL（如果不使用 curl） |
| `curl` | string | 是* | - | curl 命令（如果不使用 url） |
| `connections` | int | 否 | 10 | HTTP 连接数 |
| `duration` | string | 否 | "10s" | 测试持续时间（如 "30s", "1m"） |
| `threads` | int | 否 | 2 | 线程数 |
| `rate` | int | 否 | 0 | 请求速率限制（0=无限制） |
| `timeout` | string | 否 | "30s" | 请求超时时间 |
| `method` | string | 否 | "GET" | HTTP 方法 |
| `headers` | object | 否 | {} | HTTP 请求头 |
| `body` | string | 否 | "" | 请求体 |
| `content_type` | string | 否 | "" | Content-Type 头 |
| `use_nethttp` | bool | 否 | false | 强制使用标准库 net/http |

*注：`url` 和 `curl` 至少需要提供一个。

#### 响应示例

```json
{
  "task_id": "task_1234567890123456789",
  "status": "accepted",
  "message": "Benchmark task created and started",
  "config": {
    "url": "https://api.example.com/users",
    "connections": 100,
    "duration": "30s",
    "threads": 4,
    "rate": 0,
    "timeout": "30s",
    "method": "GET",
    "headers": {
      "Authorization": "Bearer token123"
    },
    "use_nethttp": true
  }
}
```

### 2. 查询任务状态

**GET** `/api/v1/status/:task_id`

查询指定任务的状态。

#### 响应示例

```json
{
  "id": "task_1234567890123456789",
  "status": "running",
  "created_at": "2024-01-15T10:30:00Z",
  "started_at": "2024-01-15T10:30:01Z",
  "completed_at": null,
  "error": null,
  "config": {
    "url": "https://api.example.com/users",
    "connections": 100,
    "duration": "30s"
  }
}
```

#### 任务状态

- `pending`: 任务已创建，等待执行
- `running`: 任务正在执行
- `completed`: 任务已完成
- `failed`: 任务执行失败

### 3. 获取压测结果

**GET** `/api/v1/results/:task_id`

获取指定任务的压测结果。只有当任务状态为 `completed` 时才有结果数据。

#### 响应示例

```json
{
  "id": "task_1234567890123456789",
  "status": "completed",
  "created_at": "2024-01-15T10:30:00Z",
  "started_at": "2024-01-15T10:30:01Z",
  "completed_at": "2024-01-15T10:30:31Z",
  "error": null,
  "results": {
    "total_requests": 15000,
    "total_errors": 5,
    "duration": "30.00s",
    "average_latency": "201.50ms",
    "min_latency": "45.20ms",
    "max_latency": "500.80ms",
    "latency_stddev": "67.30ms",
    "requests_per_sec": 500.00,
    "status_codes": {
      "200": 14995,
      "500": 5
    },
    "latency_percentiles": {
      "50": "180.00ms",
      "75": "220.00ms",
      "90": "280.00ms",
      "95": "350.00ms",
      "99": "450.00ms"
    },
    "total_bytes": 15728640,
    "endpoint_stats": {
      "https://api.example.com/users": {
        "requests": 15000,
        "errors": 5,
        "requests_per_sec": 500.00,
        "average_latency": "201.50ms",
        "min_latency": "45.20ms",
        "max_latency": "500.80ms",
        "status_codes": {
          "200": 14995,
          "500": 5
        },
        "total_bytes": 15728640
      }
    }
  }
}
```

### 4. 提交批量测试任务

**POST** `/api/v1/batch`

提交批量测试任务（当前版本暂未实现完整功能）。

#### 请求体示例

```json
{
  "tests": [
    {
      "name": "用户登录API",
      "curl": "curl -X POST https://api.example.com/login -H 'Content-Type: application/json' -d '{\"username\":\"test\"}'",
      "connections": 100,
      "duration": "30s",
      "threads": 4
    },
    {
      "name": "获取用户信息",
      "url": "https://api.example.com/user/profile",
      "connections": 50,
      "duration": "60s",
      "threads": 2
    }
  ]
}
```

### 5. 健康检查

**GET** `/health`

检查 API 服务器健康状态。

#### 响应示例

```json
{
  "status": "ok",
  "service": "gurl-api"
}
```

### 6. API 信息

**GET** `/`

获取 API 信息和使用说明。

## 使用示例

### 1. 提交压测任务

```bash
curl -X POST http://localhost:8080/api/v1/benchmark \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/users",
    "connections": 100,
    "duration": "30s",
    "threads": 4,
    "use_nethttp": true
  }'
```

响应：
```json
{
  "task_id": "task_1234567890123456789",
  "status": "accepted",
  "message": "Benchmark task created and started"
}
```

### 2. 查询任务状态

```bash
curl http://localhost:8080/api/v1/status/task_1234567890123456789
```

### 3. 获取压测结果

```bash
curl http://localhost:8080/api/v1/results/task_1234567890123456789
```

### 4. 使用 curl 命令提交任务

```bash
curl -X POST http://localhost:8080/api/v1/benchmark \
  -H "Content-Type: application/json" \
  -d '{
    "curl": "curl -X POST https://api.example.com/users -H '\''Content-Type: application/json'\'' -d '\''{\"name\":\"test\"}'\''",
    "connections": 100,
    "duration": "30s",
    "threads": 4,
    "use_nethttp": true
  }'
```

## 注意事项

1. **异步执行**: 所有压测任务都是异步执行的。提交任务后会立即返回 `task_id`，需要通过状态查询接口检查任务进度。

2. **任务生命周期**: 
   - 任务创建后状态为 `pending`
   - 开始执行后状态变为 `running`
   - 完成后状态变为 `completed` 或 `failed`

3. **结果获取**: 只有当任务状态为 `completed` 时，才能通过 `/api/v1/results/:id` 获取完整的压测结果。

4. **内存管理**: 任务结果会保存在内存中。建议定期清理已完成的任务，或实现持久化存储。

5. **并发限制**: 当前版本没有限制并发任务数，建议根据服务器资源合理控制并发任务数量。

## 错误处理

API 使用标准的 HTTP 状态码：

- `200 OK`: 请求成功
- `202 Accepted`: 任务已接受（提交任务时）
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 任务不存在
- `405 Method Not Allowed`: HTTP 方法不允许
- `500 Internal Server Error`: 服务器内部错误

错误响应格式：
```json
{
  "error": "错误描述信息"
}
```


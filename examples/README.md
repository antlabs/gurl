# gurl 示例文件

这个目录包含了 gurl 的各种使用示例，包括 MCP 服务配置、RESTful API 调用、批量测试等。

## 文件分类

### RESTful API 示例

1. **`api-example.sh`** - 基础 API 使用示例
2. **`api-quick-demo.sh`** - 复杂场景快速演示（推荐）
3. **`api-complex-example.sh`** - 完整的复杂场景示例
4. **`API_COMPLEX_EXAMPLES.md`** - 详细的 API 使用文档

#### JSON 配置文件

- `api-post-json.json` - POST 请求 + JSON Body 示例
- `api-get-query.json` - GET 请求 + Query 参数示例
- `api-put-form.json` - PUT 请求 + Form Data 示例
- `api-curl-command.json` - 使用完整 curl 命令的示例

### MCP 服务示例

1. `mcp_config.json` - gurl MCP 服务的工具定义配置文件
2. `claude_mcp_config.json` - Claude 应用使用的 MCP 客户端配置示例
3. `mcp_inline_batch_test.md` - 内联批量测试示例

### 批量测试示例

1. `batch-config.yaml` - 批量测试配置示例文件
2. `batch-config.json` - JSON 格式的批量测试配置

### 其他示例

1. `mock-server.yaml` - Mock 服务器配置示例
2. `custom-variables.yaml` - 自定义变量示例
3. `template-variables.yaml` - 模板变量示例
4. `multi-curls.txt` - 多个 curl 命令示例
5. `simple-test.yaml` - 简单测试配置

## 快速开始

### 1. RESTful API 使用

#### 启动 API 服务器

```bash
./gurl --api-server --api-port 3434
```

#### 运行快速演示

```bash
# 演示复杂场景（推荐）
./examples/api-quick-demo.sh

# 或使用基础示例
./examples/api-example.sh
```

#### 使用 JSON 配置文件

```bash
# POST 请求 + JSON Body
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-post-json.json

# GET 请求 + Query 参数
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-get-query.json
```

#### 复杂场景示例

```bash
# POST 请求，包含自定义 headers 和 JSON body
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/users",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json",
      "Authorization": "Bearer YOUR_TOKEN",
      "X-Request-ID": "req-001"
    },
    "body": "{\"username\":\"test\",\"email\":\"test@example.com\"}",
    "connections": 50,
    "duration": "30s",
    "threads": 10,
    "rate": 200
  }'

# 查看结果
# 获取任务 ID 后
curl "http://localhost:3434/api/v1/results/{task_id}"
```

**更多详细示例请参考**: [`API_COMPLEX_EXAMPLES.md`](./API_COMPLEX_EXAMPLES.md)

### 2. 启动 gurl MCP 服务

```bash
# 在终端中启动 gurl MCP 服务
./gurl --mcp
```

### 与 Claude 桌面应用集成

1. 打开 Claude 桌面应用
2. 进入设置 (Settings)
3. 找到 "Developer" 选项卡
4. 在 "Custom Tools" 部分，添加指向 `claude_mcp_config.json` 文件的配置

### 工具使用示例

#### 1. gurl.http_request

执行单个 HTTP 请求:

```json
{
  "name": "gurl.http_request",
  "arguments": {
    "url": "https://httpbin.org/get",
    "method": "GET"
  }
}
```

#### 2. gurl.benchmark

执行 HTTP 压测:

```json
{
  "name": "gurl.benchmark",
  "arguments": {
    "url": "https://httpbin.org/get",
    "connections": 10,
    "duration": "30s",
    "threads": 4
  }
}
```

#### 3. gurl.batch_test (使用配置文件)

执行批量测试 (传统方式):

```json
{
  "name": "gurl.batch_test",
  "arguments": {
    "config": "./examples/batch-config.yaml",
    "concurrency": 2
  }
}
```

#### 4. gurl.batch_test (内联测试定义)

执行批量测试 (无需配置文件):

```json
{
  "name": "gurl.batch_test",
  "arguments": {
    "tests": [
      {
        "name": "API 1 测试",
        "url": "https://httpbin.org/get",
        "connections": 10,
        "duration": "30s"
      },
      {
        "name": "API 2 测试",
        "curl": "curl -X POST https://httpbin.org/post -d 'data'",
        "connections": 20,
        "duration": "60s"
      }
    ],
    "concurrency": 2
  }
}
```

## 批量测试配置示例

`batch-config.yaml` 文件定义了多个测试任务:

```yaml
version: "1.0"
tests:
  - name: "用户登录API"
    curl: 'curl -X POST https://api.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 100
    duration: "30s"
    threads: 4
    
  - name: "获取用户信息"
    curl: 'curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" https://api.example.com/user/profile'
    connections: 50
    duration: "60s"
    threads: 2
```

## 常用场景

### 场景 1: 测试带认证的 API

```bash
# 使用 Bearer Token
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/protected",
    "headers": {
      "Authorization": "Bearer YOUR_JWT_TOKEN"
    },
    "connections": 20,
    "duration": "10s"
  }'

# 使用 API Key
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/data",
    "headers": {
      "X-API-Key": "your-api-key"
    },
    "connections": 30,
    "duration": "15s"
  }'
```

### 场景 2: 测试 POST 接口（JSON 数据）

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-post-json.json
```

### 场景 3: 测试带复杂 Query 参数的 GET 请求

```bash
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-get-query.json
```

### 场景 4: 使用从浏览器复制的 curl 命令

```bash
# 在浏览器开发者工具中复制 curl 命令后
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d @examples/api-curl-command.json
```

### 场景 5: 限流压测

```bash
# 限制为每秒 100 个请求
curl -X POST "http://localhost:3434/api/v1/benchmark" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://api.example.com/endpoint",
    "connections": 50,
    "duration": "30s",
    "rate": 100
  }'
```

## API 参数说明

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| `url` | string | 是* | 目标 URL | - |
| `curl` | string | 是* | curl 命令 | - |
| `method` | string | 否 | HTTP 方法 | GET |
| `headers` | object | 否 | 自定义 headers | {} |
| `body` | string | 否 | 请求体 | - |
| `content_type` | string | 否 | Content-Type | - |
| `connections` | int | 否 | 并发连接数 | 10 |
| `duration` | string | 否 | 持续时间 | 10s |
| `threads` | int | 否 | 线程数 | 2 |
| `rate` | int | 否 | 每秒请求数限制 | 0（不限） |
| `timeout` | string | 否 | 请求超时 | 30s |

\* `url` 和 `curl` 二选一

## 注意事项

- 确保 gurl 可执行文件在系统 PATH 中，或者在配置中指定完整路径
- 批量测试配置文件需要使用绝对路径或相对于 gurl 可执行文件的路径
- 内联批量测试无需配置文件，直接在工具参数中定义测试任务
- API 压测任务是异步的，需要轮询状态或等待足够时间后获取结果
- 使用 `rate` 参数可以限制每秒请求数，避免压垮服务器
- 在 JSON 中传递 JSON 字符串时，需要正确转义引号

## 相关文档

- [详细 API 使用文档](./API_COMPLEX_EXAMPLES.md)
- [主项目 README](../README.md)
- [API 文档](../docs/API.md)
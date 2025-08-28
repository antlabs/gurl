# gurl MCP 配置示例

这个目录包含了 gurl 作为 MCP (Model Context Protocol) 服务端的配置示例。

## 文件说明

1. `mcp_config.json` - gurl MCP 服务的工具定义配置文件
2. `claude_mcp_config.json` - Claude 应用使用的 MCP 客户端配置示例
3. `batch-config.yaml` - 批量测试配置示例文件
4. `mcp_inline_batch_test.md` - 内联批量测试示例

## 使用方法

### 启动 gurl MCP 服务

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

## 注意事项

- 确保 gurl 可执行文件在系统 PATH 中，或者在配置中指定完整路径
- 批量测试配置文件需要使用绝对路径或相对于 gurl 可执行文件的路径
- 内联批量测试无需配置文件，直接在工具参数中定义测试任务
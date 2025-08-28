# gurl MCP 内联批量测试示例

这个示例展示了如何使用 gurl 的 MCP 服务进行内联批量测试，无需配置文件。

## 内联批量测试示例

```json
{
  "name": "gurl.batch_test",
  "arguments": {
    "tests": [
      {
        "name": "HTTP GET 测试",
        "url": "https://httpbin.org/get",
        "connections": 10,
        "duration": "30s",
        "threads": 2
      },
      {
        "name": "HTTP POST 测试",
        "curl": "curl -X POST https://httpbin.org/post -d 'key=value'",
        "connections": 20,
        "duration": "60s",
        "threads": 4
      },
      {
        "name": "带认证的 API 测试",
        "curl": "curl -H 'Authorization: Bearer token123' https://httpbin.org/headers",
        "connections": 5,
        "duration": "15s",
        "threads": 1
      }
    ],
    "concurrency": 2,
    "verbose": true
  }
}
```

## 参数说明

- `tests`: 测试任务数组，每个任务支持以下参数：
  - `name`: 测试名称（可选）
  - `url`: 目标 URL（可与 `curl` 参数互斥）
  - `curl`: curl 命令（可与 `url` 参数互斥）
  - `connections`: 连接数（默认: 10）
  - `duration`: 测试时长（默认: "10s"）
  - `threads`: 线程数（默认: 2）
  - `rate`: 请求速率（默认: 0，无限制）
  - `timeout`: 超时时间（默认: "30s"）
  - `verbose`: 详细输出（默认: false）
  - `use_nethttp`: 使用标准库 net/http（默认: false）

- `concurrency`: 最大并发测试数（默认: 3）
- `verbose`: 详细输出（默认: false）

## 使用方法

1. 启动 gurl MCP 服务：
   ```bash
   ./gurl --mcp
   ```

2. 通过 MCP 客户端发送上述 JSON 请求

3. 等待测试完成并查看结果
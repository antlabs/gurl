# Output File 功能文档

## 概述

`output_file` 功能允许你将每个批量测试的结果保存到单独的 JSON 文件中。这对于以下场景非常有用：

- 保存测试历史记录
- 后续分析和对比
- 集成到 CI/CD 流程
- 生成测试报告
- 调试和问题排查

## 配置方式

在 `tests` 数组中的每个测试项中添加 `output_file` 字段：

```yaml
version: "1.0"
tests:
  - name: "A版本-用户查询"
    curl: 'curl -X POST https://api.example.com/user'
    connections: 10
    duration: "5s"
    output_file: "results/test-a-result.json"  # 保存到此文件
    
  - name: "B版本-用户查询"
    curl: 'curl -X POST https://api.example.com/user'
    connections: 10
    duration: "5s"
    output_file: "results/test-b-result.json"  # 保存到此文件
```

## 输出文件格式

每个输出文件都是一个 JSON 格式的文件，包含以下信息：

```json
{
  "name": "A版本-用户查询",
  "start_time": "2024-12-27T23:00:00+08:00",
  "end_time": "2024-12-27T23:00:05+08:00",
  "duration": "5s",
  "status": "success",
  "statistics": {
    "total_requests": 500,
    "total_errors": 0,
    "duration": "5s",
    "rps": "100.00",
    "avg_latency_ms": 95,
    "min_latency_ms": 50,
    "max_latency_ms": 200,
    "percentiles": {
      "p50": 90,
      "p75": 110,
      "p90": 150,
      "p95": 180,
      "p99": 195
    },
    "status_codes": {
      "200": 500
    }
  },
  "sample_response": "{\"code\":200,\"msg\":\"success\",\"data\":{...}}"
}
```

## 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 测试名称 |
| `start_time` | string | 测试开始时间（RFC3339 格式） |
| `end_time` | string | 测试结束时间（RFC3339 格式） |
| `duration` | string | 测试持续时间 |
| `status` | string | 测试状态：`success` 或 `failed` |
| `error` | string | 错误信息（仅在失败时存在） |
| `statistics` | object | 测试统计信息 |
| `sample_response` | string | 样本响应（用于 jsondiff） |

### statistics 对象

| 字段 | 类型 | 说明 |
|------|------|------|
| `total_requests` | number | 总请求数 |
| `total_errors` | number | 总错误数 |
| `duration` | string | 测试持续时间 |
| `rps` | string | 每秒请求数 |
| `avg_latency_ms` | number | 平均延迟（毫秒） |
| `min_latency_ms` | number | 最小延迟（毫秒） |
| `max_latency_ms` | number | 最大延迟（毫秒） |
| `percentiles` | object | 延迟百分位数 |
| `status_codes` | object | HTTP 状态码分布 |

## 使用场景

### 场景 1：保存 A/B 测试结果

```yaml
version: "1.0"
tests:
  - name: "A版本-接口"
    curl: 'curl https://api-v1.example.com/endpoint'
    connections: 100
    duration: "30s"
    output_file: "results/version-a.json"
    
  - name: "B版本-接口"
    curl: 'curl https://api-v2.example.com/endpoint'
    connections: 100
    duration: "30s"
    output_file: "results/version-b.json"

jsondiff:
  pairs:
    - base: "A版本-接口"
      target: "B版本-接口"
      compare_field: "data"
```

运行后，你将得到两个文件：
- `results/version-a.json` - A 版本的测试结果
- `results/version-b.json` - B 版本的测试结果

### 场景 2：多接口测试结果归档

```yaml
version: "1.0"
tests:
  - name: "登录接口"
    curl: 'curl -X POST https://api.example.com/login -d "..."'
    connections: 50
    duration: "10s"
    output_file: "results/login-test.json"
    
  - name: "用户信息接口"
    curl: 'curl https://api.example.com/user/info'
    connections: 30
    duration: "10s"
    output_file: "results/userinfo-test.json"
    
  - name: "订单列表接口"
    curl: 'curl https://api.example.com/orders'
    connections: 40
    duration: "10s"
    output_file: "results/orders-test.json"
```

### 场景 3：CI/CD 集成

在 CI/CD 流程中保存测试结果，用于后续分析：

```yaml
version: "1.0"
tests:
  - name: "性能测试-主接口"
    curl: 'curl https://api.example.com/main'
    connections: 100
    duration: "60s"
    output_file: "ci-results/build-${BUILD_ID}-main.json"
```

然后在 CI 脚本中分析结果：

```bash
# 运行测试
gurl --batch-config test.yaml

# 分析结果
python analyze_results.py ci-results/build-${BUILD_ID}-main.json

# 上传到监控系统
curl -X POST https://monitor.example.com/api/results \
  -H "Content-Type: application/json" \
  -d @ci-results/build-${BUILD_ID}-main.json
```

## 文件路径说明

### 相对路径

使用相对路径时，文件将相对于当前工作目录创建：

```yaml
output_file: "results/test.json"  # 相对于当前目录
```

### 绝对路径

也可以使用绝对路径：

```yaml
output_file: "/var/log/gurl/test.json"  # 绝对路径
```

### 目录自动创建

**注意**：目前目录不会自动创建，请确保目标目录已存在。如果目录不存在，文件写入将失败。

建议在运行测试前创建目录：

```bash
mkdir -p results
gurl --batch-config test.yaml
```

## 与 jsondiff 结合使用

`output_file` 功能与 `jsondiff` 完美配合：

```yaml
version: "1.0"
tests:
  - name: "老版本-获取地区156"
    curl: 'curl https://api-v1.example.com/region/156'
    connections: 1
    threads: 1
    requests: 1
    output_file: "results/old-version.json"
    
  - name: "新版本-获取地区156"
    curl: 'curl https://api-v2.example.com/region/156'
    connections: 1
    threads: 1
    requests: 1
    output_file: "results/new-version.json"

jsondiff:
  pairs:
    - name: "版本对比"
      base: "老版本-获取地区156"
      target: "新版本-获取地区156"
      base_field: "result"
      target_field: "data"
```

运行后：
1. 每个测试的完整结果保存到各自的文件
2. jsondiff 对比结果显示在控制台
3. 可以后续分析保存的文件

## 最佳实践

### 1. 使用有意义的文件名

```yaml
# 好的命名
output_file: "results/2024-12-27-login-api-test.json"

# 不好的命名
output_file: "test1.json"
```

### 2. 按日期或版本组织

```yaml
# 按日期
output_file: "results/2024-12-27/login-test.json"

# 按版本
output_file: "results/v1.2.3/api-test.json"

# 按环境
output_file: "results/production/api-test.json"
```

### 3. 在 verbose 模式下查看保存状态

```bash
gurl --batch-config test.yaml --verbose
```

输出将显示：
```
Test result saved to: results/test-a-result.json
Test result saved to: results/test-b-result.json
```

### 4. 错误处理

如果文件保存失败（如目录不存在、权限不足），测试仍会继续执行，但会在 verbose 模式下显示警告：

```
Warning: failed to save result to file 'results/test.json': open results/test.json: no such file or directory
```

## 示例文件

项目提供了示例配置文件：

- `examples/batch-with-output-files.yaml` - 展示如何使用 output_file 功能

## 命令行使用

```bash
# 基本使用
gurl --batch-config batch-with-output-files.yaml

# 使用 verbose 模式查看文件保存状态
gurl --batch-config batch-with-output-files.yaml --verbose

# 结合其他选项
gurl --batch-config batch-with-output-files.yaml --verbose --batch-report json
```

## 注意事项

1. **文件覆盖**：如果指定的文件已存在，将被覆盖
2. **目录权限**：确保有写入目标目录的权限
3. **磁盘空间**：大量测试可能产生大量文件，注意磁盘空间
4. **并发安全**：不同测试使用不同的文件名，避免冲突
5. **可选功能**：`output_file` 是可选的，不指定则不保存

## 后续分析

保存的 JSON 文件可以用于：

### Python 分析

```python
import json

# 读取测试结果
with open('results/test-a-result.json', 'r') as f:
    result = json.load(f)

# 分析性能
print(f"RPS: {result['statistics']['rps']}")
print(f"P95 Latency: {result['statistics']['percentiles']['p95']}ms")
print(f"Error Rate: {result['statistics']['total_errors'] / result['statistics']['total_requests'] * 100}%")
```

### jq 查询

```bash
# 查看 RPS
jq '.statistics.rps' results/test-a-result.json

# 查看 P95 延迟
jq '.statistics.percentiles.p95' results/test-a-result.json

# 对比两个测试的 RPS
diff <(jq '.statistics.rps' results/test-a-result.json) \
     <(jq '.statistics.rps' results/test-b-result.json)
```

### 监控集成

将结果发送到监控系统：

```bash
# Prometheus Pushgateway
cat results/test-a-result.json | \
  jq -r '.statistics | "rps \(.rps)\navg_latency_ms \(.avg_latency_ms)"' | \
  curl --data-binary @- http://pushgateway:9091/metrics/job/gurl_test

# InfluxDB
curl -X POST 'http://influxdb:8086/write?db=gurl' \
  --data-binary "gurl_test,name=test-a rps=$(jq '.statistics.rps' results/test-a-result.json)"
```

## 相关文档

- [Batch Testing 设计文档](./batch-design.md)
- [JSONDiff 设计文档](./jsondiff-design.md)
- [API 文档](./API.md)

# JSONDiff 功能设计文档

## 概述

JSONDiff 功能用于在批量测试中对比不同版本API的JSON响应，实现A/B测试场景下的响应一致性验证。

## 1. 基本结构

在现有的 batch 配置中，每个测试项维持原有结构：

- 使用 `curl` 字符串描述请求
- 使用 `connections`、`duration`、`threads` 等字段描述压测参数

在此基础上，在配置文件根级别增加一个可选字段 `jsondiff`：

```yaml
version: "1.0"

tests:
  - name: "A版本-用户登录"
    curl: 'curl -X POST https://api-v1.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 1
    threads: 1
    requests: 1
    
  - name: "B版本-用户登录"
    curl: 'curl -X POST https://api-v2.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 1
    threads: 1
    requests: 1

# JSON响应对比配置
jsondiff:
  pairs:
    - name: "登录接口A/B对比"
      base: "A版本-用户登录"
      target: "B版本-用户登录"
      compare_field: "data"
      ignore_fields:
        - "timestamp"
        - "request_id"
```

## 2. 配置字段说明

### 2.1 jsondiff 顶层配置

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `pairs` | []JSONDiffPair | 是 | 对比对列表，每个对比对定义一组需要对比的测试 |

### 2.2 JSONDiffPair 配置

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 否 | 对比对的名称，如不指定则自动生成为 "base vs target" |
| `base` | string | 是 | 基准测试的名称（引用 tests 数组中的 name） |
| `target` | string | 是 | 目标测试的名称（引用 tests 数组中的 name） |
| `compare_field` | string | 否 | 要对比的JSON字段路径（使用gjson语法），同时应用于base和target。如不指定则对比整个响应 |
| `base_field` | string | 否 | 基准测试要对比的字段路径，用于老版本和新版本字段名不同的场景 |
| `target_field` | string | 否 | 目标测试要对比的字段路径，用于老版本和新版本字段名不同的场景 |
| `ignore_fields` | []string | 否 | 要忽略的字段列表，支持嵌套路径 |

**字段优先级说明**：
- 如果指定了 `base_field` 和 `target_field`，则使用它们（支持不同字段路径对比）
- 如果只指定了 `compare_field`，则同时应用于base和target
- 如果都不指定，则对比整个响应

## 3. 使用场景

### 3.1 场景一：只对比data字段

**最常用场景**：只关心业务数据是否一致，忽略code、msg等元数据。

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
      compare_field: "data"  # 只对比data字段
```

### 3.2 场景二：对比code-msg-data结构

**标准响应结构**：对比整个响应，但忽略动态字段。

```yaml
jsondiff:
  pairs:
    - name: "用户查询对比"
      base: "A版本-用户查询"
      target: "B版本-用户查询"
      compare_field: "data"
      ignore_fields:
        - "timestamp"
        - "request_id"
        - "trace_id"
        - "created_at"
```

### 3.3 场景三：对比嵌套字段

**深层字段对比**：只对比响应中的特定嵌套结构。

```yaml
jsondiff:
  pairs:
    - base: "A版本-订单详情"
      target: "B版本-订单详情"
      compare_field: "data.order.items"  # 只对比订单项
      ignore_fields:
        - "update_time"
        - "last_modified"
```

### 3.4 场景四：批量对比多个接口

**多接口对比**：一次性对比多个API端点。

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

### 3.5 场景五：对比不同字段路径（新功能）

**跨版本字段名变更**：老版本使用 `result`，新版本使用 `data`。

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
      base_field: "result"      # 老版本从 result 提取
      target_field: "data"      # 新版本从 data 提取
      ignore_fields:
        - "timestamp"
        - "request_id"
```

**响应示例**：

老版本响应：
```json
{
  "code": 200,
  "msg": "success",
  "result": {
    "region": "Beijing",
    "code": "156",
    "name": "China"
  },
  "timestamp": "2024-01-01"
}
```

新版本响应：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "region": "Beijing",
    "code": "156",
    "name": "China"
  },
  "timestamp": "2024-01-02"
}
```

使用 `base_field: "result"` 和 `target_field: "data"` 后，只会对比：
- 老版本的 `result` 内容
- 新版本的 `data` 内容

即使字段名不同，只要内容一致就会通过对比。

## 4. 字段路径语法

字段路径使用点号（`.`）分隔的嵌套结构：

| 路径示例 | 说明 |
|---------|------|
| `data` | 根级别的data字段 |
| `data.user` | data对象中的user字段 |
| `data.user.profile` | 多层嵌套字段 |
| `result.items.0` | 数组中的第一个元素 |
| `timestamp` | 任何级别的timestamp字段（用于ignore_fields） |

### 4.1 ignore_fields 匹配规则

`ignore_fields` 支持多种匹配方式：

1. **精确匹配**：`"data.timestamp"` 只忽略 `data.timestamp`
2. **字段名匹配**：`"timestamp"` 忽略任何级别的 `timestamp` 字段
3. **前缀匹配**：`"data.user"` 忽略 `data.user` 及其所有子字段

示例：

```yaml
ignore_fields:
  - "timestamp"           # 忽略所有timestamp字段
  - "data.created_at"     # 只忽略data下的created_at
  - "request_id"          # 忽略所有request_id字段
  - "trace_id"            # 忽略所有trace_id字段
```

## 5. 实现原理

### 5.1 执行流程

```
1. 执行批量测试
   ├─ 在benchmark开始前，捕获一个样本响应
   ├─ 运行性能测试
   └─ 收集测试结果和样本响应

2. 执行jsondiff对比
   ├─ 遍历所有配置的对比对
   ├─ 查找base和target的样本响应
   ├─ 根据compare_field提取要对比的字段
   ├─ 递归对比JSON结构
   └─ 生成对比结果

3. 生成报告
   ├─ 批量测试结果
   ├─ jsondiff对比结果
   └─ 差异详情
```

### 5.2 对比算法

```go
func CompareJSONResponses(base, target []byte, compareField string, ignoreFields []string) JSONDiffResult {
    // 1. 如果指定了compare_field，提取该字段
    if compareField != "" {
        base = extractField(base, compareField)
        target = extractField(target, compareField)
    }
    
    // 2. 解析JSON
    var baseData, targetData interface{}
    json.Unmarshal(base, &baseData)
    json.Unmarshal(target, &targetData)
    
    // 3. 递归对比
    differences := compareValues("", baseData, targetData, ignoreFields)
    
    // 4. 生成结果
    return JSONDiffResult{
        Passed: len(differences) == 0,
        Differences: differences,
    }
}
```

### 5.3 数据结构

```go
// JSONDiffConfig 定义JSON对比配置
type JSONDiffConfig struct {
    Pairs []JSONDiffPair `yaml:"pairs" json:"pairs"`
}

// JSONDiffPair 表示一对需要对比的测试
type JSONDiffPair struct {
    Name         string   `yaml:"name,omitempty" json:"name,omitempty"`
    Base         string   `yaml:"base" json:"base"`
    Target       string   `yaml:"target" json:"target"`
    CompareField string   `yaml:"compare_field,omitempty" json:"compare_field,omitempty"`
    IgnoreFields []string `yaml:"ignore_fields,omitempty" json:"ignore_fields,omitempty"`
}

// JSONDiffResult 表示对比结果
type JSONDiffResult struct {
    PairName    string
    BaseName    string
    TargetName  string
    Passed      bool
    Message     string
    Differences []FieldDifference
}

// FieldDifference 表示字段差异
type FieldDifference struct {
    Field       string
    BaseValue   string
    TargetValue string
}
```

## 6. 输出报告

### 6.1 文本报告格式

```
=== Batch Test Report ===

Total Tests: 2
Success Rate: 100.00%
Total Time: 10.5s

=== Test Results ===

1. A版本-用户登录
   Duration: 5.2s
   Status: SUCCESS
   Requests: 500
   RPS: 96.15

2. B版本-用户登录
   Duration: 5.3s
   Status: SUCCESS
   Requests: 530
   RPS: 100.00

=== JSON Diff Results ===

Comparison: 登录接口A/B对比
  Base: A版本-用户登录
  Target: B版本-用户登录
  Status: PASSED
  Message: field 'data' is identical in both responses

JSONDiff Summary: 1 passed, 0 failed
```

### 6.2 失败时的差异详情

```
=== JSON Diff Results ===

Comparison: 用户信息对比
  Base: A版本-用户信息
  Target: B版本-用户信息
  Status: FAILED
  Message: found 2 difference(s)
  Differences:
    - Field: data.user.age
      Base:   25
      Target: 26
    - Field: data.user.city
      Base:   Beijing
      Target: Shanghai

JSONDiff Summary: 0 passed, 1 failed
```

## 7. 使用建议

### 7.1 最佳实践

1. **优先使用 compare_field**
   - 使用 `compare_field: "data"` 只对比业务数据
   - 避免对比元数据（code、msg等）

2. **合理设置 ignore_fields**
   - 始终忽略时间戳字段：`timestamp`、`created_at`、`updated_at`
   - 忽略请求追踪字段：`request_id`、`trace_id`
   - 忽略服务器信息：`server_time`、`server_id`

3. **渐进式测试**
   - 先测试一个对比对，验证配置正确
   - 再逐步添加更多对比对

4. **使用 verbose 模式**
   - 运行时添加 `--verbose` 查看详细对比过程
   - 帮助调试配置问题

### 7.2 性能考虑

1. **样本响应捕获**
   - 在benchmark开始前捕获一次样本响应
   - 不影响性能测试的准确性

2. **内存使用**
   - 每个测试只保存一个样本响应
   - 大响应体可能占用较多内存

3. **对比性能**
   - 对比在所有测试完成后执行
   - 不影响并发测试的执行

## 8. 限制和注意事项

### 8.1 当前限制

1. **数组顺序敏感**
   - 数组元素必须按相同顺序排列
   - 未来版本将支持无序数组对比

2. **浮点数精确比较**
   - 浮点数必须完全相等
   - 未来版本将支持容差配置

3. **单样本对比**
   - 只对比一个样本响应
   - 未来版本将支持多样本统计对比

### 8.2 注意事项

1. **确保端点可访问**
   - 样本响应在benchmark前捕获
   - 确保测试端点在执行时可访问

2. **响应格式一致**
   - base和target应返回相同格式的JSON
   - 不同格式会导致对比失败

3. **字段路径正确**
   - `compare_field` 必须在两个响应中都存在
   - 路径错误会导致对比失败

## 9. 示例文件

项目提供了以下示例文件：

- `examples/batch-jsondiff-simple.yaml` - 最简单的使用示例
- `examples/batch-jsondiff.yaml` - 完整功能示例
- `examples/batch-jsondiff-codemsg.yaml` - code-msg-data结构示例
- `examples/JSONDIFF_README.md` - 详细使用文档

## 10. 命令行使用

```bash
# 基本使用
gurl --batch-config batch-jsondiff.yaml

# 查看详细对比过程
gurl --batch-config batch-jsondiff.yaml --verbose

# 生成JSON报告
gurl --batch-config batch-jsondiff.yaml --batch-report json > results.json

# 生成CSV报告
gurl --batch-config batch-jsondiff.yaml --batch-report csv > results.csv
```

## 11. 未来扩展

### 11.1 计划功能

1. **无序数组对比**
   ```yaml
   compare_mode: "array_unordered"
   ```

2. **数值容差**
   ```yaml
   numeric_tolerance: 0.001
   ```

3. **正则表达式匹配**
   ```yaml
   field_patterns:
     - field: "data.order_no"
       pattern: "^ORD\\d{10}$"
   ```

4. **自定义对比函数**
   ```yaml
   custom_comparators:
     - field: "data.amount"
       function: "compare_currency"
   ```

5. **多样本统计对比**
   ```yaml
   sampling:
     enabled: true
     count: 10
     strategy: "random"
   ```

## 12. 相关文档

- [Batch Testing 设计文档](./batch-design.md)
- [Assert 设计文档](./assert-design.md)
- [API 文档](./API.md)

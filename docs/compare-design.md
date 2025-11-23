# Compare 功能设计文档

本文档记录在 `gurl` 中新增一个用于对比多个 HTTP 请求响应结果的 `compare` 功能的设计草案。

目标是：

- 方便用户直接粘贴 Postman 导出的 `curl` 命令，定义多个请求。
- 定义一个或多个 `compare` 场景，对不同请求的响应进行对比。
- 支持按 **gjson 路径（`gjson` 关键字，底层使用 [gjson](https://github.com/tidwall/gjson)）** 粒度进行响应体字段的细粒度比较。
- 借鉴 **hurl** 风格的断言语法，对响应和对比结果进行断言。
- 支持一对一、一对多、多对多等批量对比模式。

---

## 1. 配置文件整体结构概览

配置文件建议使用 YAML（也可扩展为 JSON），典型结构如下：

```yaml
requests:
  - name: cache_on
    curl: |
      curl --location 'https://example.com/api/data?user_id=123' \
        --header 'Accept: application/json' \
        --header 'Cache-Control: max-age=600'

  - name: cache_off
    curl: |
      curl --location 'https://example.com/api/data?user_id=123' \
        --header 'Accept: application/json' \
        --header 'Cache-Control: no-cache'

compare:
  - name: cache_behavior_compare
    base: cache_on
    target: cache_off

    response_compare: |
      status == status
      header "Content-Type" == header "Content-Type"
      header "Date" ignore
      gjson "data.user.id" == gjson "data.user.id"
      gjson "data.items" == gjson "data.items"
      gjson "meta.request_id" exists
```

- `requests`：定义要发送的 HTTP 请求，直接嵌入 `curl` 字符串，支持从 Postman 粘贴。
- `compare`：定义一个或多个对比场景：
  - 指定基准请求 `base` 和目标请求 `target`（或 `targets`）。
  - 通过 `response_compare`（hurl 风格断言）定义更细粒度的对比规则，其中 JSON 字段对比通过 `gjson` 关键字调用 gjson 路径表达式。

---

## 2. requests：直接使用 curl 命令

### 2.1 设计目标

- 用户可以从 Postman 或其他工具直接复制 `curl` 命令到配置文件中。
- 不要求用户手动拆解 method / url / headers / body。
- 内部在执行 compare 时，将 `curl` 解析为内部请求结构。

### 2.2 YAML 示例

```yaml
requests:
  - name: cache_on
    curl: |
      curl --location 'https://example.com/api/data?user_id=123' \
        --header 'Accept: application/json' \
        --header 'Cache-Control: max-age=600' \
        --header 'X-Trace-Id: req-cache-on'

  - name: cache_off
    curl: |
      curl --location 'https://example.com/api/data?user_id=123' \
        --header 'Accept: application/json' \
        --header 'Cache-Control: no-cache' \
        --header 'X-Trace-Id: req-cache-off'
```

### 2.3 内部解析思路（实现建议）

- 提供一个工具函数，例如：
  - `ParseCurlToRequest(curl string) (*Request, error)`
- 解析内容包括：
  - HTTP method
  - URL（含 query）
  - headers
  - body（例如 `--data`, `--data-raw`, `--data-binary` 等）
- 兼容常见 curl 导出风格（Postman / 浏览器 DevTools 等）。

---

## 3. compare：基本一对一对比场景

### 3.1 基本结构

```yaml
compare:
  - name: cache_behavior_compare
    base: cache_on
    target: cache_off

    response_compare: |
      status == status
      header "Content-Type" == header "Content-Type"
      header "ETag" == header "ETag"
      header "Date" ignore
      header "Server" ignore
      gjson "data.user.id" == gjson "data.user.id"
      gjson "data.items" == gjson "data.items"
      gjson "meta.request_id" exists
      gjson "meta.env" exists
```

- `base`：作为基准的请求名，对应 `requests` 中的 `name`。
- `target`：要与基准对比的请求名。
- `response_compare`：使用 hurl 风格的断言语法，作为多行字符串（类似 `asserts: |`），每一行是一个断言表达式，支持对 status、headers、gjson 路径等进行精细对比。

---

## 4. hurl 风格断言语法设计

### 4.1 对比两个响应的断言

断言的主体是一个字符串，例如：

- `"status == status"`
- `"header[Content-Type] == header[Content-Type]"`
- `"header[Date] ignore"`
- `"gjson data.user.id == gjson data.user.id"`
- `"gjson meta.request_id exists"`

语义（约定）：

- `status == status`：`base.status == target.status`。
- `header[Name] == header[Name]`：两个响应中同名 header 值必须相等。
- `header[Name] ignore`：该 header 差异被忽略（常用于 `Date`、`Server` 等）。
- `gjson <path> == gjson <path>`：
  - `gjson.GetBytes(base.body, path) == gjson.GetBytes(target.body, path)`（这里的 `path` 使用 gjson 路径语法，例如 `data.user.id`）。
- `gjson <path> ~= gjson <path> ± <tolerance>`：
  - 取出的值为数值型，比较绝对差值是否不超过 `tolerance`。
- `gjson <path> exists`：
  - 仅验证字段存在（通常用于 `meta.request_id` 等）。

### 4.2 内部结构建议

解析后可映射为内部结构体，例如：

```go
type CompareAssert struct {
    Scope     string  // "status" / "header" / "gjson"
    Key       string  // header name 或 gjson 路径
    Operator  string  // "==", "~=", "exists", "ignore"
    Tolerance float64 // 近似比较的容忍度，可选
}
```

执行时：

- 先根据 `Scope` 决定取什么值：
  - `status`：`int`
  - `header`：`map[string]string`
  - `gjson`：使用 [gjson](https://github.com/tidwall/gjson) 从响应体 JSON 中按路径取值（内部可先解析为 `[]byte`，再用 gjson；路径语法为 gjson path，如 `data.items.0.id`）。
- 再根据 `Operator` 执行相应比较逻辑。

### 4.3 与单请求断言的关系

在 `compare` 模式中，`response_compare` 中的每一条 `assert` 表达式，都是在 **base 与 target 两个响应之间** 做比较，例如：

- `"status == status"` 表示 `base.status == target.status`。
- `"gjson data.user.id == gjson data.user.id"` 表示 `gjson.GetBytes(base.body, "data.user.id") == gjson.GetBytes(target.body, "data.user.id")`。
- `"gjson meta.request_id exists"` 表示 base/target 响应中都要求存在该字段。

在单请求断言模式中，我们采用同样的 hurl 风格断言语法，只是比较对象变成了「单个响应 vs 期望值」，并通过 YAML 中的 `asserts` 多行字符串来承载，例如：

```yaml
version: "1.0"

tests:
  - name: "用户登录API"
    curl: 'curl -X POST https://api.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 10
    duration: "5s"
    threads: 1

    asserts: |
      status == 200
      header "Content-Type" contains "application/json"
      gjson "success" == true
      body contains "login success"
      duration_ms < 500
```

两者共享同一套目标（`status` / `header` / `gjson` / `body` / `duration_ms`）和操作符语义，只是 compare 模式下会在内部同时取 base/target 两侧的值进行对比，而单请求模式下只针对单个响应进行断言。

---

## 6. 批量多个请求的对比模式

### 6.1 一对多：一个 base，对多个 target

适用于一个基准请求与多个变体比较（例如不同缓存策略）：

```yaml
requests:
  - name: base_cache_on
    curl: |
      curl 'https://example.com/api/data?user_id=123' \
        --header 'Cache-Control: max-age=600'

  - name: cache_off
    curl: |
      curl 'https://example.com/api/data?user_id=123' \
        --header 'Cache-Control: no-cache'

  - name: cache_short
    curl: |
      curl 'https://example.com/api/data?user_id=123' \
        --header 'Cache-Control: max-age=60'

  - name: cache_very_long
    curl: |
      curl 'https://example.com/api/data?user_id=123' \
        --header 'Cache-Control: max-age=3600'

compare:
  - name: cache_batch_compare
    base: base_cache_on

    response_compare: |
      status == status
      gjson "data" == gjson "data"
      gjson "meta.request_id" exists
```

执行时：

- 先发一次 `base_cache_on` 请求，得到 `base_response`。
- 对 `targets` 中的每个 request：
  - 发请求得到 `target_response`。
  - 按 `response_compare` 逐条断言，形成每一对的对比结果。

### 6.2 多对多（按索引配对）

适用于「环境 A 和环境 B」各有一组请求，顺序一一对应的场景：

```yaml
request_sets:
  left:
    - name: a_1
      curl: |
        curl 'https://a-env/api/1'
    - name: a_2
      curl: |
        curl 'https://a-env/api/2'

  right:
    - name: b_1
      curl: |
        curl 'https://b-env/api/1'
    - name: b_2
      curl: |
        curl 'https://b-env/api/2'

compare:
  - name: batch_pair_by_index
    mode: pair_by_index
    left_set: left
    right_set: right

    response_compare: |
      status == status
      gjson "data" == gjson "data"
```

执行时：

- `left[0]` vs `right[0]`，`left[1]` vs `right[1]`，依此类推。

### 6.3 按 group / role 分组对比

适用于多条请求按会话或业务分组的场景：

```yaml
requests:
  - name: s1_cache_on
    group: session_1
    role: base
    curl: |
      curl 'https://example.com/api/s1' --header 'Cache-Control: max-age=600'

  - name: s1_cache_off
    group: session_1
    role: variant
    curl: |
      curl 'https://example.com/api/s1' --header 'Cache-Control: no-cache'

  - name: s2_cache_on
    group: session_2
    role: base
    curl: |
      curl 'https://example.com/api/s2' --header 'Cache-Control: max-age=600'

  - name: s2_cache_off
    group: session_2
    role: variant
    curl: |
      curl 'https://example.com/api/s2' --header 'Cache-Control: no-cache'

compare:
  - name: group_by_session
    mode: group_by_field
    group_field: group
    base_role: base
    target_role: variant

    response_compare: |
      status == status
      gjson "data" == gjson "data"
```

执行时：

- 按 `group_field`（如 `group`）分组，在每个组内：
  - 找到 `base_role`（如 `base`）对应请求作为基准。
  - 找到 `target_role`（如 `variant`）请求作为对比对象。

---

## 7. 命令行使用方式

### 7.1 基本用法

```bash
gurl compare -f cache_compare.yaml -n cache_behavior_compare
```

- `-f`：指定配置文件路径（如 `cache_compare.yaml`）。
- `-n`：指定要执行的 compare 场景名称（`compare.name`）。

执行流程示意：

1. 读取配置文件，解析 `requests`、`compare` 等。
2. 根据场景（`base/target` 或批量模式）确定要发送的请求对。
3. 用 `curl` 字符串解析成内部请求结构并发起 HTTP 请求。
4. 对每一对响应按 `compare_fields` 和 `response_compare`/`gjson_compare` 进行比对。
5. 输出每一对的断言结果和汇总。

### 7.2 预期输出示意

```text
Scenario: cache_behavior_compare
Base   : cache_on
Target : cache_off

[OK]  status == status
     base:   200
     target: 200

[OK]  header[Content-Type] == header[Content-Type]
     base:   application/json
     target: application/json

[OK]  gjson data.user.id == gjson data.user.id
     base:   "123"
     target: "123"

[FAIL] gjson data.items == gjson data.items
     base:   [{"id":1,"name":"A"},{"id":2,"name":"B"}]
     target: [{"id":1,"name":"A"},{"id":2,"name":"B"},{"id":3,"name":"C"}]

Summary: 4 passed, 1 failed
Exit code: 1
```

在批量模式下，建议按“每一对 / 每个 target”为一级分组输出，例如：

```text
=== Pair: base_cache_on vs cache_off ===
[OK]  status == status
[OK]  gjson data == gjson data

=== Pair: base_cache_on vs cache_short ===
[FAIL] gjson data == gjson data
  base : ...
  target: ...

Summary:
  total pairs: 3
  passed     : 2
  failed     : 1
```

---

## 8. 后续实现建议

- **解析层**：
  - 使用已有的 YAML/JSON 解析库，将配置映射到 Go struct。
  - 实现 `ParseCurlToRequest`，将 `curl` 字符串解析为内部请求结构（method/url/headers/body）。
- **执行层**：
  - 根据 compare 场景决定请求对的组合方式（一对一、一对多、多对多、分组）。
  - 为每个请求对发送 HTTP 请求，收集响应（status/headers/body）。
- **断言层**：
  - 实现 hurl 风格 `assert` 字符串解析为 `CompareAssert`。
  - 使用 gjson 对 body 进行路径提取，并执行 equals/approx/exists/ignore 等逻辑。
- **输出层**：
  - 支持人类可读的文本输出，显示每条断言的 [OK]/[FAIL] 与差异细节。
  - 可考虑增加 machine-readable 格式（如 JSON），用于 CI/自动化。

---

## 9. response_compare 的断言语法

`response_compare` 中的断言语法与单请求断言对齐，支持以下几种类型（其中 `gjson` 使用 gjson 路径）：

- `status == 200`：状态码检查
- `header "Content-Type" contains "application/json"`：头部检查
  - `gjson "data.id" > 0`：基于 gjson 路径的 JSON 字段检查
- `body contains "success"`：响应体检查
- `duration_ms < 500`：响应时间检查

具体语法见 [`docs/assert-design.md`](./assert-design.md)。

---

## 10. 单请求断言与 compare 模式的统一设计

除了本文件重点介绍的「多请求比较（compare）」模式之外，gurl 还支持对 **单个 HTTP 请求** 做断言检查，两者在语法上保持一致：

- **请求描述**：都使用 curl 命令字符串定义 HTTP 请求。
- **断言语法**：都采用 hurl 风格的断言表达式（`gjson` 关键字底层使用 gjson）：
  - `status == 200`
  - `header "Content-Type" contains "application/json"`
  - `gjson "data.id" > 0`
  - `body contains "success"`
  - `duration_ms < 500`

在 compare 模式中，断言嵌在 `response_compare` 多行字符串中，例如：

```yaml
response_compare: |
  status == status
  gjson "data" == gjson "data"
  gjson "meta.request_id" exists
```

在单请求模式中，断言则写在 YAML 的 `asserts` 字段中，例如：

```yaml
version: "1.0"

tests:
  - name: "用户登录API"
    curl: 'curl -X POST https://api.example.com/login -H "Content-Type: application/json" -d "{\"username\":\"test\",\"password\":\"123456\"}"'
    connections: 10
    duration: "5s"
    threads: 1

    asserts: |
      status == 200
      header "Content-Type" contains "application/json"
      gjson "success" == true
      body contains "login success"
      duration_ms < 500
```

可以看到：

- **compare 模式**：`assert` 表达式用于「base ↔ target」两边的字段对比；
- **单请求模式**：`asserts` 中的每行表达式用于「请求 ↔ 期望值」的检查。

二者共享同一套目标与操作符定义，详细语法见 [`docs/assert-design.md`](./assert-design.md)。

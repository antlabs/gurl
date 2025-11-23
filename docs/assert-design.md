# HTTP 断言设计（gurl）

本设计说明 gurl 在 YAML 配置中如何为 **单个 HTTP 请求** 描述断言。
目标是：**请求部分沿用 curl，断言部分使用 hurl 风格语法，但嵌入在 YAML 中。**

---

## 1. 基本结构

在现有的 batch 配置中，每个测试项维持原有结构：

- 使用 `curl` 字符串描述请求
- 使用 `connections`、`duration`、`threads` 等字段描述压测参数

在此基础上，为每个测试项增加一个可选字段 `asserts`：

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

- **请求部分**：完全兼容现有的 `batch-config.yaml` 风格，仅使用 `curl` 字符串。
- **断言部分**：`asserts` 是一个多行字符串，每一行是一条 **hurl 风格** 的断言语句。
- 行首 `#` 视为注释，空行会被忽略。

---

## 2. 断言语法概览

断言语句的整体风格参考 hurl，简化为以下形式：

```text
<target> <operator> <expected>
```

其中：

- `<target>`：断言目标（状态码 / Header / JSON 字段（gjson 路径）/ Body / 时延等）
- `<operator>`：比较或匹配操作符
- `<expected>`：期望值（数字 / 字符串 / 正则 / 布尔）

### 2.1 支持的 target

gurl 当前计划支持以下断言目标：

- `status`  
  - HTTP 状态码  
  - 示例：`status == 200`

- `header "<Name>"`  
- 指定响应头，对应 `http.Header.Get(Name)` 的字符串值。  
- **字符串比较：** 使用字符串运算符 `==`、`!=`、`contains`、`not_contains`、`starts_with`、`ends_with`，右侧字符串必须用双引号包裹，例如：  
  - `header "Content-Type" == "application/json"`  
  - `header "Content-Type" contains "json"`  
- **存在性检查：** 支持无右值的 `exists` / `not_exists`：  
  - `header "X-Request-Id" exists`  
  - `header "X-Debug-Flag" not_exists`  
- **不支持：** `header` 目前不支持 `matches /.../` 正则运算符（正则仅在 `gjson` 目标中可用）。

- `gjson "<path>"`  
  - 针对 JSON 响应体，使用 [gjson](https://github.com/tidwall/gjson) 路径表达式进行字段提取（例如 `friends.#(last=="Murphy").first`、`data.items.0.name`）  
  - 示例：  
    - `gjson "success" == true`  
    - `gjson "data.id" > 0`  
    - `gjson "items.0.name" matches /^[A-Z].+$/`

- `body`  
  - 原始响应体文本（字符串）  
  - 示例：  
    - `body contains "OK"`  
    - `body not_contains "error"`

- `duration_ms`  
  - 单次请求的耗时（毫秒）  
  - 示例：  
    - `duration_ms < 500`  
    - `duration_ms <= 1000`

> 说明：  
> - `gjson` 使用 [gjson](https://github.com/tidwall/gjson) 解析，仅在响应体为 JSON 时生效。  
> - `duration_ms` 取自实际请求的测量时间。

### 2.2 支持的 operator

根据 target 类型不同，可使用的运算符包括：

- **比较运算**（数值 / 字符串 / 布尔）
  - `==`、`!=`
  - `>`、`>=`、`<`、`<=`（主要用于数值，如状态码、duration、JSON 数值字段）

- **字符串运算**
  - `contains` / `not_contains`
  - `starts_with` / `ends_with`

- **正则匹配（仅 `gjson` 支持）**
-  - `matches`  
-  - 右侧为正则表达式，使用 `/.../` 包裹，仅用于 `gjson` 目标：
-    - `gjson "token" matches /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$/`

- **存在性检查**
  - `exists` / `not_exists`  
  - 常用于 Header / JSON 字段（gjson 路径）：
    - `header "X-Trace-Id" exists`
    - `gjson "data" not_exists`

### 2.3 expected 值的写法

- **数字**：直接写  
  - `status == 200`  
  - `duration_ms < 500`  
  - `gjson "data.id" > 0`

- **字符串**：用双引号包裹  
  - `header "Content-Type" contains "json"`  
  - `body contains "success"`

- **布尔**：`true` / `false`  
  - `gjson "success" == true`

- **正则**：使用 `/.../` 包裹  
  - `gjson "token" matches /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$/`

---

## 3. 示例

### 3.1 简单 GET 接口断言

```yaml
version: "1.0"

tests:
  - name: "HTTP GET Test"
    curl: 'curl https://httpbin.org/get'
    connections: 5
    duration: "5s"
    threads: 1

    asserts: |
      status == 200
      header "Content-Type" contains "application/json"
      gjson "url" contains "httpbin.org/get"
      duration_ms < 1000
```

### 3.2 带 JSON Body 的 POST 接口断言

```yaml
version: "1.0"

tests:
  - name: "创建订单API"
    curl: 'curl -X POST https://api.example.com/orders -H "Content-Type: application/json" -d "{\"product_id\":1,\"quantity\":2}"'
    connections: 10
    duration: "10s"

    asserts: |
      status == 201
      header "Content-Type" contains "application/json"
      gjson "order_id" > 0
      gjson "status" == "created"
      body not_contains "error"
      duration_ms < 800
```

---

## 4. 与 hurl 的关系

- 相同点：
  - 使用类似 hurl 的断言关键字和表达式风格（`status == 200`、`header`、`gjson`、`body` 等）。
  - 每一行即一条断言，支持注释与空行。

- 不同点：
  - **请求部分**：gurl 仍然通过 `curl` 字符串定义请求；  
    而 hurl 在同一个文件中用专门的 HTTP 请求块（`GET /path HTTP/1.1` + 头 + Body）。
  - **嵌入形式**：hurl 直接是 `.hurl` 文件；  
    gurl 将断言作为 YAML 中的一个字段 `asserts`（多行字符串）。

这种设计保证：

- 对已有 `batch-config.yaml` 兼容：只需在单个 test 上添加 `asserts` 字段即可。
- 对 hurl 用户友好：断言语法接近 hurl，无需重新学习一套完全不同 DSL。

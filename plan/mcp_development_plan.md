# gurl MCP 功能开发计划

## 1. 项目概述
为 gurl 添加 MCP (Model Context Protocol) 服务端支持，使其能够作为 MCP 服务端与支持 MCP 的客户端（如 Claude 桌面应用）通信，提供 gurl 的功能作为 AI 工具使用。

## 2. 核心功能列表
- 实现 MCP 服务端，遵循 Model Context Protocol 规范。
- 将 gurl 的现有功能（HTTP 请求、压测、批量测试、模板变量等）封装为 MCP 工具。
- 提供友好的交互方式，允许用户通过 MCP 客户端配置和执行 gurl 命令。

## 3. 技术选型
- **MCP SDK**: 使用 Go 语言实现，采用 [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 库来简化 MCP 服务端的开发。
- **依赖管理**: 将 MCP 相关代码作为新的模块添加到 `internal/mcp`，不干扰现有功能。
- **启动方式**: 添加一个新的子命令 `gurl mcp` 来启动 MCP 服务。

## 4. 开发计划

### 阶段一：基础环境与框架搭建 (已完成)
1.  **MCP 规范研究**:
    *   深入研究 [MCP 规范](https://modelcontextprotocol.io/)，理解其核心概念（如 `initialize`, `tools/list`, `tools/call` 等）。
    *   确定需要实现的核心功能点。
2.  **依赖引入**:
    *   在 `go.mod` 中添加 [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 依赖。
    *   确认 `internal/mcp` 目录已创建用于存放所有 MCP 相关代码。
3.  **基础服务端框架**:
    *   使用 `mark3labs/mcp-go` 库实现一个基础的 MCP 服务端框架。
    *   实现 `initialize` 和 `shutdown` 等基础生命周期方法。
    *   添加 `gurl mcp` 子命令用于启动服务。
    *   实现了第一个工具 `gurl.http_request`，可以执行基本的 HTTP 请求。

### 阶段二：工具注册与发现 (部分完成)
1.  **定义 gurl 工具**:
    *   将 gurl 的核心功能（如执行单次 HTTP 请求、执行压测、执行批量测试）封装为 MCP 工具。
    *   定义每个工具的名称、描述、输入参数（对应 gurl 的命令行参数）。
    *   例如：`gurl.http_request`, `gurl.benchmark`, `gurl.batch_test`。
    *   已完成 `gurl.http_request` 工具的定义和实现。
2.  **实现 `tools/list`**:
    *   实现 `tools/list` 方法，返回所有已定义的 gurl 工具列表。
    *   通过 `mark3labs/mcp-go` 库自动处理。
3.  **实现 `tools/call` (基础)**:
    *   实现 `tools/call` 方法的基本框架，能够接收工具调用请求并根据工具名称分发到对应的处理函数。
    *   已完成 `gurl.http_request` 工具的调用处理。

### 阶段三：功能实现与集成 (进行中)
1.  **实现 `gurl.http_request` 工具**:
    *   创建处理函数，接收 HTTP 请求相关的参数（URL, method, headers, body 等）。
    *   调用现有的 `parser` 和 `nethttp_client` 模块执行请求。
    *   将响应结果（状态码、headers、body）格式化后返回给 MCP 客户端。
    *   **已完成基础实现**。
2.  **实现 `gurl.benchmark` 工具**:
    *   创建处理函数，接收压测相关的参数（URL, connections, duration, threads, curl command 等）。
    *   调用现有的 `benchmark` 模块执行压测。
    *   将压测结果（统计数据）格式化后返回。
    *   **已完成基础实现**。
3.  **实现 `gurl.batch_test` 工具**:
    *   创建处理函数，接收批量测试相关的参数（配置文件路径、并发数等）。
    *   调用现有的 `batch` 模块执行批量测试。
    *   将批量测试报告格式化后返回。
    *   **已完成基础实现**。
4.  **实现 `gurl.template_variables` 工具 (可选)**:
    *   提供查看和管理模板变量的功能。
    *   **待实现**。

### 阶段四：优化与测试 (已完成)
1.  **错误处理**:
    *   完善各工具调用的错误处理逻辑，确保错误信息能清晰地传递给 MCP 客户端。
    *   **已完成基础错误处理**。
2.  **日志记录**:
    *   添加适当的日志记录，方便调试和监控.
    *   **已完成基础日志记录**。
3.  **交互优化**:
    *   优化与 MCP 客户端的交互体验，例如提供更友好的参数提示。
    *   **已完成工具描述和参数提示优化**。
4.  **测试**:
    *   使用支持 MCP 的客户端（如 Claude 桌面应用）进行手动测试。
    *   编写单元测试覆盖核心逻辑。
    *   **已完成基础单元测试**。

## 5. 预期成果
- gurl 可以作为独立的 MCP 服务端运行。
- 用户可以通过支持 MCP 的客户端调用 gurl 的各种功能。
- 提供与命令行类似的能力，但通过图形化界面或 AI 代理进行交互。
- **当前进展**: 已完整实现 MCP 服务端框架，并完成了 `gurl.http_request`、`gurl.benchmark` 和 `gurl.batch_test` 工具的实现，可以执行基本的 HTTP 请求、压测和批量测试。已完成基础的错误处理、日志记录、交互优化和单元测试。

## 6. 后续扩展 (可选)
- 提供更丰富的上下文信息，例如历史请求记录。
- 支持更复杂的交互，例如实时查看压测进度。
- 实现 `gurl.template_variables` 工具，提供查看和管理模板变量的功能。
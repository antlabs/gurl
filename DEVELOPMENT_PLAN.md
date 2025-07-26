# murl 开发计划

## 项目概述
murl 是一个类似 wrk 的高性能 HTTP 压测工具，支持 curl 命令解析和多种压测模式。本文档描述了项目的后续开发计划。

## 开发路线图

### 1. 完善单元测试和集成测试 🧪
**优先级**: 高  
**预计工期**: 2-3天

#### 功能描述
- 完善现有模块的单元测试覆盖率
- 创建集成测试验证端到端功能
- 实现测试用HTTP服务端用于集成测试

#### 技术实现
- 为所有核心模块添加单元测试 (`*_test.go`)
- 在测试中创建轻量级HTTP服务端
- 实现基准测试验证性能
- 添加CI/CD流水线自动化测试

#### 测试覆盖范围
- `internal/config/` - 配置解析测试
- `internal/parser/` - curl命令解析测试
- `internal/benchmark/` - 压测逻辑测试
- `internal/stats/` - 统计功能测试
- 集成测试 - 端到端功能验证

#### 文件结构
```
internal/config/config_test.go
internal/parser/parser_test.go
internal/benchmark/
├── benchmark_test.go
├── nethttp_client_test.go
└── pulse_client_test.go
internal/stats/stats_test.go
tests/
├── integration_test.go  # 集成测试
├── testserver.go       # 测试用HTTP服务端
└── benchmark_test.go   # 基准测试
```

---

### 2. 批量压测配置文件支持 📋
**优先级**: 高  
**预计工期**: 2-3天

#### 功能描述
- 支持从YAML/JSON配置文件导入多个curl命令
- 实现批量压测功能
- 提供详细的批量测试报告

#### 技术实现
- 扩展 `internal/config/` 模块支持配置文件解析
- 创建批量测试配置格式规范
- 实现并发批量测试执行器
- 生成汇总报告

#### 配置文件格式示例
```yaml
# murl-batch.yaml
version: "1.0"
tests:
  - name: "用户登录API"
    curl: 'curl -X POST https://api.example.com/login -d "username=test&password=123"'
    connections: 100
    duration: "30s"
    
  - name: "获取用户信息"
    curl: 'curl -H "Authorization: Bearer token" https://api.example.com/user/profile'
    connections: 50
    duration: "60s"
```

#### 文件结构
```
internal/batch/
├── config.go          # 批量配置解析
├── executor.go        # 批量执行器
└── report.go          # 批量报告生成

examples/
└── batch-config.yaml  # 配置文件示例
```

---

### 3. URL模板变量支持 🔧
**优先级**: 中  
**预计工期**: 2-3天

#### 功能描述
- 支持curl命令中的URL模板变量
- 实现变量替换和随机生成
- 支持多种变量类型 (随机数、UUID、时间戳等)

#### 技术实现
- 扩展 `internal/parser/` 模块支持模板解析
- 实现变量生成器
- 支持自定义变量函数

#### 模板语法示例
```bash
# 支持的模板变量
murl -c 100 -d 30s --parse-curl 'curl https://api.example.com/user/{{.UserID}}/posts/{{.PostID}}'

# 变量定义
--var UserID=random:1-1000
--var PostID=uuid
--var Timestamp=now
```

#### 支持的变量类型
- `{{.random:min-max}}` - 随机数
- `{{.uuid}}` - UUID
- `{{.timestamp}}` - 时间戳
- `{{.sequence:start}}` - 递增序列
- `{{.choice:opt1,opt2,opt3}}` - 随机选择

#### 文件结构
```
internal/template/
├── parser.go          # 模板解析器
├── variables.go       # 变量生成器
└── functions.go       # 内置函数
```

---

### 4. Web UI界面 🎨
**优先级**: 中  
**预计工期**: 5-7天

#### 功能描述
- 创建现代化的Web界面
- 实时显示压测进度和结果
- 支持配置管理和历史记录查看

#### 技术栈
- **后端**: Go + Gin框架 + embed静态资源
- **前端**: 原生HTML/CSS/JavaScript + Chart.js
- **实时通信**: WebSocket
- **样式**: Tailwind CSS (CDN)
- **图表**: Chart.js (CDN)

#### 功能特性
- 📊 实时压测监控面板
- ⚙️ 可视化配置编辑器
- 📈 交互式结果图表
- 📝 测试历史记录
- 🔄 批量测试管理
- 📱 响应式设计

#### 技术实现
- 使用Go 1.16+ embed指令将静态资源打包到二进制文件
- 原生JavaScript实现，无需构建工具和Node.js依赖
- 使用CDN加载外部库（Tailwind CSS、Chart.js）
- WebSocket实现实时数据推送
- 单二进制文件部署，无需额外配置

#### 文件结构
```
web/
├── server.go          # Web服务器主程序
├── handlers/
│   ├── api.go         # API处理器
│   └── websocket.go   # WebSocket处理
├── static/            # 静态资源文件
│   ├── index.html     # 主页面
│   ├── app.js         # 主逻辑
│   ├── style.css      # 样式文件
│   └── components/    # 组件文件
└── embed.go           # Go embed静态资源
```

---

### 5. 接口级别TPS统计 📊
**优先级**: 中  
**预计工期**: 2-3天

#### 功能描述
- 实现按接口URL分组的TPS统计
- 提供详细的性能分析报告
- 支持多维度数据展示

#### 技术实现
- 扩展 `internal/stats/` 模块
- 实现URL分组统计
- 添加百分位延迟统计
- 生成详细的性能报告

#### 统计维度
- 每个接口的TPS (事务/秒)
- 响应时间分布 (P50, P90, P95, P99)
- 错误率统计
- 状态码分布
- 请求/响应字节统计

#### 报告格式示例
```
接口性能统计报告
==================

API: POST /api/login
├── TPS: 1,234.56 req/s
├── 平均响应时间: 45.2ms
├── P95响应时间: 89.1ms
├── 成功率: 99.8%
└── 错误分布: 404(0.1%), 500(0.1%)

API: GET /api/user/profile  
├── TPS: 2,456.78 req/s
├── 平均响应时间: 23.4ms
├── P95响应时间: 56.7ms
├── 成功率: 100%
└── 错误分布: 无
```

#### 文件结构
```
internal/stats/
├── collector.go       # 统计数据收集
├── analyzer.go        # 数据分析器
├── reporter.go        # 报告生成器
└── interface_stats.go # 接口级统计
```

---

## 开发时间线

### 第一阶段 (Week 1-2)
- ✅ 完善单元测试和集成测试
- ✅ 实现批量压测配置文件支持

### 第二阶段 (Week 3-4)  
- ✅ 实现URL模板变量支持
- ✅ 完成接口级别TPS统计

### 第三阶段 (Week 5-6)
- ✅ 开发Web UI界面
- ✅ 集成所有功能模块

### 第四阶段 (Week 7)
- ✅ 完整测试和文档更新
- ✅ 发布新版本

---

## 技术债务和优化

### 性能优化
- 优化内存使用，减少GC压力
- 改进连接池管理
- 优化统计数据收集效率

### 代码质量
- 增加单元测试覆盖率
- 完善错误处理机制
- 改进日志记录

### 文档完善
- 更新README文档
- 添加使用示例
- 创建API文档

---

## 发布计划

### v2.0.0 - 批量测试版本
- 完善测试覆盖率
- 批量压测支持
- URL模板变量

### v2.1.0 - 统计增强版本  
- 接口级TPS统计
- 详细性能报告

### v3.0.0 - Web UI版本
- 完整Web界面
- 实时监控
- 历史记录管理

---

## 贡献指南

欢迎社区贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与开发。

### 开发环境要求
- Go 1.21+
- Git

### 开发流程
1. Fork项目
2. 创建功能分支
3. 提交代码
4. 创建Pull Request

---

*最后更新: 2025-07-26*

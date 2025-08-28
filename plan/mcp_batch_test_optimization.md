# gurl MCP 批量测试优化计划

## 项目概述
优化 gurl 的 MCP 批量测试功能，使其无需配置文件即可直接执行批量测试任务。

## 当前状态
已完成！gurl 的 MCP 批量测试功能现在支持两种模式：
1. 使用配置文件的传统模式
2. 直接通过参数传递测试任务的内联模式

## 优化目标
实现一个可以直接通过 MCP 工具参数传递批量测试任务的方案，无需单独的配置文件。

## 技术实现方案

### 方案一：内联批量测试配置
修改 `gurl.batch_test` 工具，支持直接传递测试任务列表：

```json
{
  "name": "gurl.batch_test",
  "arguments": {
    "tests": [
      {
        "name": "API 1",
        "url": "https://httpbin.org/get",
        "connections": 10,
        "duration": "30s"
      },
      {
        "name": "API 2",
        "curl": "curl -X POST https://httpbin.org/post -d 'data'",
        "connections": 20,
        "duration": "60s"
      }
    ],
    "concurrency": 2
  }
}
```

### 方案二：混合模式
保持对配置文件的支持，同时增加内联配置的支持：

```json
{
  "name": "gurl.batch_test",
  "arguments": {
    "config": "/path/to/batch-config.yaml",
    // 或者
    "tests": [...]
  }
}
```

## 实现细节

### 1. 修改工具定义
更新 `gurl.batch_test` 工具的输入模式，支持直接传递测试任务。

### 2. 修改处理逻辑
在 `handleBatchTest` 函数中，解析直接传递的测试任务而不是读取配置文件。

### 3. 创建测试任务结构
定义测试任务的数据结构，与现有的 `config.BatchTest` 兼容。

### 4. 实现任务执行
复用现有的批量执行器逻辑，但传入直接定义的任务列表。

## 预期收益
- ✅ 简化 MCP 客户端的使用流程
- ✅ 无需管理额外的配置文件
- ✅ 提高工具的易用性和灵活性

## 兼容性考虑
- ✅ 保持与现有 API 的兼容性
- ✅ 支持两种模式：配置文件和内联测试定义

## 优先级
高

## 预计工期
1-2天

## 实际完成时间
2025-08-29
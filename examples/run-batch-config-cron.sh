#!/usr/bin/env bash

set -euo pipefail

# 使用 batch-config 定时运行批量测试示例
# 这里默认使用 examples/asserts-fast.yaml，并每秒运行一次

./build/gurl \
  --batch-config examples/asserts-fast.yaml \
  --batch-concurrency 3 \
  --batch-report text \
  --schedule-cron "*/1 * * * * *"

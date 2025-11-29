#!/usr/bin/env bash

set -euo pipefail

./build/gurl --parse-curl-file examples/multi-curls.txt -c 1 -n 1 --schedule-cron "*/1 * * * * *"

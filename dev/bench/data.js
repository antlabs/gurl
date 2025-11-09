window.BENCHMARK_DATA = {
  "lastUpdate": 1762668502415,
  "repoUrl": "https://github.com/antlabs/gurl",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "guonaihong@qq.com",
            "name": "guonaihong",
            "username": "guonaihong"
          },
          "committer": {
            "email": "guonaihong@qq.com",
            "name": "guonaihong",
            "username": "guonaihong"
          },
          "distinct": true,
          "id": "b596b64cec3970241e3db812eff9d1f541c310a9",
          "message": "fix: resolve request body reuse causing data race in concurrent benchmarks\n\n- Store body content as string and recreate io.Reader for each request to prevent shared state\n- Clone request with fresh body reader before each HTTP call to ensure thread safety\n- Simplify body reader initialization in request builder",
          "timestamp": "2025-11-09T01:08:58+08:00",
          "tree_id": "ed6719350ec9d3abf3a2a907f3a9f3ed3e98842e",
          "url": "https://github.com/antlabs/gurl/commit/b596b64cec3970241e3db812eff9d1f541c310a9"
        },
        "date": 1762668501664,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100079083,
            "unit": "ns/op\t 5177151 B/op\t   63863 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100079083,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5177151,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 63863,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 400.6,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14965449 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 400.6,
            "unit": "ns/op",
            "extra": "14965449 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14965449 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14965449 times\n4 procs"
          }
        ]
      }
    ]
  }
}
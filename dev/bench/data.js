window.BENCHMARK_DATA = {
  "lastUpdate": 1762671246899,
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
      },
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
          "id": "74b8559a46e30a77c529e327843ff7afc60da712",
          "message": "refactor: improve error handling in HTTP request panic recovery\n\n- Changed handleHTTPRequest to use named return values for proper error propagation\n- Updated BuildRequest to use io.Reader interface instead of concrete *strings.Reader type",
          "timestamp": "2025-11-09T14:22:57+08:00",
          "tree_id": "05806cf78b25f326aec769c1b29940412bc851fe",
          "url": "https://github.com/antlabs/gurl/commit/74b8559a46e30a77c529e327843ff7afc60da712"
        },
        "date": 1762669419620,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100091781,
            "unit": "ns/op\t 5219493 B/op\t   64417 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100091781,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5219493,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 64417,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 403.8,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14989899 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 403.8,
            "unit": "ns/op",
            "extra": "14989899 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14989899 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14989899 times\n4 procs"
          }
        ]
      },
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
          "id": "35ea90eaec24b711320dec3f5e8a0edc87449a08",
          "message": "chore: upgrade Go version to 1.25.x in CI workflows",
          "timestamp": "2025-11-09T14:31:20+08:00",
          "tree_id": "088c4749998da5afa3f3d766133cb80056d0d5b5",
          "url": "https://github.com/antlabs/gurl/commit/35ea90eaec24b711320dec3f5e8a0edc87449a08"
        },
        "date": 1762669937358,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100063357,
            "unit": "ns/op\t 5105725 B/op\t   62964 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100063357,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5105725,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62964,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 391.1,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15178192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 391.1,
            "unit": "ns/op",
            "extra": "15178192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15178192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15178192 times\n4 procs"
          }
        ]
      },
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
          "id": "1abea35e188ecaa87ffd4674b3763d613690977e",
          "message": "chore: update Go version to 1.25",
          "timestamp": "2025-11-09T14:42:36+08:00",
          "tree_id": "fa52dc8291838750c1ae446c97afbdbbdce5a011",
          "url": "https://github.com/antlabs/gurl/commit/1abea35e188ecaa87ffd4674b3763d613690977e"
        },
        "date": 1762671246665,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100069053,
            "unit": "ns/op\t 5182832 B/op\t   63897 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100069053,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5182832,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 63897,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 394.5,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15135516 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 394.5,
            "unit": "ns/op",
            "extra": "15135516 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15135516 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15135516 times\n4 procs"
          }
        ]
      }
    ]
  }
}
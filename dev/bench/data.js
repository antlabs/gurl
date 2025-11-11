window.BENCHMARK_DATA = {
  "lastUpdate": 1762875902300,
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
          "id": "4884bd4a7b5dbe76d4bed60771355b57325d0bd6",
          "message": "chore: disable lint job and downgrade Go version requirement\n\n- Commented out golangci-lint job in CI workflow\n- Removed lint dependency from build job\n- Downgraded Go version from 1.25 to 1.24.1",
          "timestamp": "2025-11-09T14:59:55+08:00",
          "tree_id": "ae4cff2cb3697050e0fe32e0c5eb4a635e5c1890",
          "url": "https://github.com/antlabs/gurl/commit/4884bd4a7b5dbe76d4bed60771355b57325d0bd6"
        },
        "date": 1762671630302,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100067283,
            "unit": "ns/op\t 7656909 B/op\t   94490 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100067283,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 7656909,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 94490,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 412.7,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14530971 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 412.7,
            "unit": "ns/op",
            "extra": "14530971 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14530971 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14530971 times\n4 procs"
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
          "id": "a782d830cb02d6585474780fa3566ac67012cd92",
          "message": "fix: configure custom temp directory for Go tests\n\n- Set GOTMPDIR in CI workflows and Makefile to use project-local temp directories\n- Update test files to create temp directories in current working directory instead of /tmp\n- Add temp directory patterns to .gitignore to keep repository clean",
          "timestamp": "2025-11-09T15:12:37+08:00",
          "tree_id": "ad4453ded27d29aa1d6356d8053a48de4ae890e1",
          "url": "https://github.com/antlabs/gurl/commit/a782d830cb02d6585474780fa3566ac67012cd92"
        },
        "date": 1762672391819,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100074104,
            "unit": "ns/op\t 7178920 B/op\t   88574 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100074104,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 7178920,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 88574,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 425.1,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "13691923 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 425.1,
            "unit": "ns/op",
            "extra": "13691923 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "13691923 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "13691923 times\n4 procs"
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
          "id": "cb29d51649b43454c85856de6e5c1e595f029a18",
          "message": "chore: temporarily disable Windows CI builds",
          "timestamp": "2025-11-09T15:15:50+08:00",
          "tree_id": "ed09cd12be984a6656a281c103cb572804661c8b",
          "url": "https://github.com/antlabs/gurl/commit/cb29d51649b43454c85856de6e5c1e595f029a18"
        },
        "date": 1762672586386,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100083753,
            "unit": "ns/op\t 4970049 B/op\t   61291 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100083753,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 4970049,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 61291,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 398.9,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15067323 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 398.9,
            "unit": "ns/op",
            "extra": "15067323 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15067323 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15067323 times\n4 procs"
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
          "id": "c4b07f930deb176d4810d408894fdf5ac048d351",
          "message": "chore: update pulse dependency and fix code formatting\n\n- Updated antlabs/pulse to version e31dbf68d422 for latest improvements\n- Standardized struct field alignment and removed trailing whitespace across codebase\n- Modified pulse implementation to only handle HTTP (not HTTPS) connections",
          "timestamp": "2025-11-09T16:30:02+08:00",
          "tree_id": "982e580b146150ccce080e467360b7faf72cb7e9",
          "url": "https://github.com/antlabs/gurl/commit/c4b07f930deb176d4810d408894fdf5ac048d351"
        },
        "date": 1762677100952,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100080682,
            "unit": "ns/op\t 4937854 B/op\t   60846 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100080682,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 4937854,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 60846,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 399.5,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15175929 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 399.5,
            "unit": "ns/op",
            "extra": "15175929 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15175929 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15175929 times\n4 procs"
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
          "id": "f741522c949243b2279714a3b10516ca97949e08",
          "message": "fix: suppress output in LiveUI mode and extract sampling logic\n\n- Conditionally hide initial benchmark info when LiveUI is enabled to prevent UI interference\n- Extract sampling goroutine into shared StartSampling function to eliminate code duplication\n- Remove debug print statements from pulse_client that would disrupt LiveUI display",
          "timestamp": "2025-11-09T16:47:25+08:00",
          "tree_id": "ac47de7a13f192dee3a6f13a5825dd42169b4763",
          "url": "https://github.com/antlabs/gurl/commit/f741522c949243b2279714a3b10516ca97949e08"
        },
        "date": 1762678098497,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100069862,
            "unit": "ns/op\t 5092284 B/op\t   62782 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100069862,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5092284,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62782,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 393.5,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15389095 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 393.5,
            "unit": "ns/op",
            "extra": "15389095 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15389095 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15389095 times\n4 procs"
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
          "id": "b58f03a2c94226e6e6bfb07d5f49cbe92f4d90c1",
          "message": "refactor: improve UI initialization timing in benchmark\n\n- Move sampling goroutine start before connection establishment to capture metrics earlier\n- Add initial render call to prevent blank screen during UI startup\n- Include debug logging placeholder for troubleshooting metrics updates",
          "timestamp": "2025-11-09T16:56:11+08:00",
          "tree_id": "1587c62a2dd4add42e53e2ef9d185324f9023add",
          "url": "https://github.com/antlabs/gurl/commit/b58f03a2c94226e6e6bfb07d5f49cbe92f4d90c1"
        },
        "date": 1762678599793,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100082744,
            "unit": "ns/op\t 5030609 B/op\t   62043 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100082744,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5030609,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62043,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 400.4,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15006542 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 400.4,
            "unit": "ns/op",
            "extra": "15006542 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15006542 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15006542 times\n4 procs"
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
          "id": "67b52418f62f2476890a3aa1976dfc257469fd36",
          "message": "feat: add RESTful API server for remote benchmark management\n\n- Implemented HTTP API server with endpoints for submitting, monitoring, and retrieving benchmark results\n- Added UI theme customization with auto-detection, dark, and light modes (including Ultimate Light Theme V4.0)\n- Created comprehensive API documentation with examples for all endpoints",
          "timestamp": "2025-11-11T23:44:22+08:00",
          "tree_id": "748da1b6e5d213fc4427a43d39a9f64453918291",
          "url": "https://github.com/antlabs/gurl/commit/67b52418f62f2476890a3aa1976dfc257469fd36"
        },
        "date": 1762875901481,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100131303,
            "unit": "ns/op\t 5034700 B/op\t   62104 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100131303,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5034700,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62104,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 405.2,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15016854 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 405.2,
            "unit": "ns/op",
            "extra": "15016854 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15016854 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15016854 times\n4 procs"
          }
        ]
      }
    ]
  }
}
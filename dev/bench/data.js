window.BENCHMARK_DATA = {
  "lastUpdate": 1764259655095,
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
          "id": "3dc2805a60747417e45a95eed0372cf5f4084983",
          "message": "feat: configure logging to error level only\n\n- Set default log level to error in main application to reduce noise\n- Configure pulse benchmark to only show error logs, preventing INFO logs from interfering with UI display",
          "timestamp": "2025-11-11T23:55:54+08:00",
          "tree_id": "44e518aeea3d802a5e32c02ab816b1502d69be71",
          "url": "https://github.com/antlabs/gurl/commit/3dc2805a60747417e45a95eed0372cf5f4084983"
        },
        "date": 1762876642297,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100091817,
            "unit": "ns/op\t 5222180 B/op\t   64452 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100091817,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5222180,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 64452,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 397.5,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15195699 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 397.5,
            "unit": "ns/op",
            "extra": "15195699 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15195699 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15195699 times\n4 procs"
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
          "id": "ea74ac82a0594141349ffc2db7cc8228e71a6682",
          "message": "fix: use actual runtime duration and improve chart number formatting\n\n- Changed results.Duration to use actual elapsed time instead of configured duration to support early interruption\n- Added compact number formatting (K/M suffixes) for request chart to prevent display overflow\n- Adjusted bar chart spacing and added dynamic title to indicate format when using large numbers",
          "timestamp": "2025-11-12T00:13:20+08:00",
          "tree_id": "29de0454114efd715965dad2b7f467ab1f47d47a",
          "url": "https://github.com/antlabs/gurl/commit/ea74ac82a0594141349ffc2db7cc8228e71a6682"
        },
        "date": 1762877633565,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100079848,
            "unit": "ns/op\t 4939548 B/op\t   60930 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100079848,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 4939548,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 60930,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 409.6,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "13996408 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 409.6,
            "unit": "ns/op",
            "extra": "13996408 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "13996408 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "13996408 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8b0def42a7c4480db6e429d7f72685975ec2de11",
          "message": "Merge pull request #10 from antlabs/dependabot/github_actions/actions/checkout-5\n\nchore(deps): bump actions/checkout from 4 to 5",
          "timestamp": "2025-11-17T23:35:08+08:00",
          "tree_id": "db084de82299fc165f693c9a11fedffd37c21f81",
          "url": "https://github.com/antlabs/gurl/commit/8b0def42a7c4480db6e429d7f72685975ec2de11"
        },
        "date": 1763393754901,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100087127,
            "unit": "ns/op\t 5097979 B/op\t   62849 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100087127,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5097979,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62849,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 391.5,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15404581 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 391.5,
            "unit": "ns/op",
            "extra": "15404581 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15404581 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15404581 times\n4 procs"
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
          "id": "cca6c2bd1b7be3847fb618e3f7b773bee4c7f992",
          "message": "feat: add compare mode for response comparison between base and target endpoints\n\n- Added --compare-config/-f and --compare-name/-n flags to enable compare mode execution\n- Implemented runCompare function to load configuration, execute scenarios, and display results\n- Added result grouping by PairLabel to support batch comparison modes (one_to_many, pair_by_index, group_by_field)\n- Display comparison details including pass/fail status, base/target values, and failure reasons\n- Compare mode takes",
          "timestamp": "2025-11-23T15:25:38+08:00",
          "tree_id": "91235ce534d5312d987c37d230532a4a39b68fd3",
          "url": "https://github.com/antlabs/gurl/commit/cca6c2bd1b7be3847fb618e3f7b773bee4c7f992"
        },
        "date": 1763883643610,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100079815,
            "unit": "ns/op\t 5072578 B/op\t   62558 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100079815,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5072578,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 62558,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 408.4,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14979574 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 408.4,
            "unit": "ns/op",
            "extra": "14979574 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14979574 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14979574 times\n4 procs"
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
          "id": "aff7589f2a1a116ebe0b5f5861ac8d0261515abf",
          "message": "feat: support multiple HTTP methods per path in mock server\n\n- Changed mock server port from 8080 to 9191 in example configuration\n- Refactored route registration to group routes by path and dispatch by method within handler\n- Fixed conflict when registering multiple methods for the same path pattern\n- Moved method matching logic inside handler to avoid duplicate pattern registration errors",
          "timestamp": "2025-11-23T22:29:55+08:00",
          "tree_id": "71f617a4e8adb1f54bb33845af22051ecad95c49",
          "url": "https://github.com/antlabs/gurl/commit/aff7589f2a1a116ebe0b5f5861ac8d0261515abf"
        },
        "date": 1763908229579,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100084224,
            "unit": "ns/op\t 4977303 B/op\t   61361 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100084224,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 4977303,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 61361,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 413.9,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14701317 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 413.9,
            "unit": "ns/op",
            "extra": "14701317 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14701317 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14701317 times\n4 procs"
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
          "id": "100aae52c1490b4e970972ab624c976f8c3bda2d",
          "message": "feat: improve assertion error handling and add quoted string literal support\n\n- Added unquoteMaybe helper to support quoted string literals in header and body assertions\n- Fixed success rate calculation to account for per-request assertion failures in addition to top-level errors\n- Updated batch reporter to display stats (requests/RPS/latency) for failed tests with assertion errors\n- Added comments clarifying that tests are considered successful only when both top-level error and stats errors are absent",
          "timestamp": "2025-11-23T23:25:50+08:00",
          "tree_id": "a5299cc4da45f60412aacd628b82f8ec0273ed6c",
          "url": "https://github.com/antlabs/gurl/commit/100aae52c1490b4e970972ab624c976f8c3bda2d"
        },
        "date": 1763911582694,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100081160,
            "unit": "ns/op\t 5126351 B/op\t   63264 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100081160,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5126351,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 63264,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 399.9,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14873554 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 399.9,
            "unit": "ns/op",
            "extra": "14873554 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14873554 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14873554 times\n4 procs"
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
          "id": "b8d086a517e3ed0c02d2467884fcd8df46168e6b",
          "message": "docs: add header assertion documentation and implement assertion support in pulse_client\n\n- Documented header assertion syntax with string comparison operators (==, !=, contains, starts_with, ends_with)\n- Clarified that header assertions use http.Header.Get() string values and require quoted strings\n- Added existence checks (exists/not_exists) for header assertions\n- Noted that regex matches operator is only supported for gjson targets, not headers\n- Implemented assertion evaluation in pulse_client to",
          "timestamp": "2025-11-23T23:53:44+08:00",
          "tree_id": "bb830caededaf28387b9b1bace8511e5c13be5ba",
          "url": "https://github.com/antlabs/gurl/commit/b8d086a517e3ed0c02d2467884fcd8df46168e6b"
        },
        "date": 1763914614357,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100092003,
            "unit": "ns/op\t 4986079 B/op\t   61486 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100092003,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 4986079,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 61486,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 409.4,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "12963630 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 409.4,
            "unit": "ns/op",
            "extra": "12963630 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "12963630 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "12963630 times\n4 procs"
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
          "id": "1b3bf3fe6222b47524a2a5af9a91a5eba9f1a5e9",
          "message": "feat: add total request limit support with -n/--requests flag\n\n- Added --requests/-n flag to limit total number of requests (0=unlimited, duration-limited)\n- Implemented atomic CAS-based request counting in nethttp_client and pulse_client to enforce request limit\n- Modified worker goroutines to cancel context when request limit is reached\n- Split request counting logic: CAS increment when Requests>0, simple increment when Requests=0\n- Improved error reporting to group identical errors and show count",
          "timestamp": "2025-11-25T00:18:10+08:00",
          "tree_id": "ff38271cdebb551bdbeb4ccdeda940f024c3d8aa",
          "url": "https://github.com/antlabs/gurl/commit/1b3bf3fe6222b47524a2a5af9a91a5eba9f1a5e9"
        },
        "date": 1764001131113,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100093181,
            "unit": "ns/op\t 5168361 B/op\t   63751 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100093181,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5168361,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 63751,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 411,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14704760 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 411,
            "unit": "ns/op",
            "extra": "14704760 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14704760 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14704760 times\n4 procs"
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
          "id": "53731b9f34246be1a301293d5479c2efb1eff7d4",
          "message": "feat: add multi-request support to pulse client and improve write traffic tracking\n\n- Added PulseBenchmarkMulti struct to support multiple requests in pulse client\n- Implemented NewPulseBenchmarkWithMultipleRequests to create pulse benchmark with request pool\n- Modified HTTPClientHandler to support both single and multi-request modes via requestPool field\n- Added write traffic tracking with AddWriteBytes calls after each request write operation\n- Enhanced output to display write traffic statistics (",
          "timestamp": "2025-11-26T00:30:23+08:00",
          "tree_id": "e80dcd0f58e54e7cd01defeacca664c5ce9b9292",
          "url": "https://github.com/antlabs/gurl/commit/53731b9f34246be1a301293d5479c2efb1eff7d4"
        },
        "date": 1764088326257,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100068863,
            "unit": "ns/op\t 7957087 B/op\t   98226 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100068863,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 7957087,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 98226,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 419.8,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14308395 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 419.8,
            "unit": "ns/op",
            "extra": "14308395 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14308395 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14308395 times\n4 procs"
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
          "id": "c629e99897d430933065853a601399377183907e",
          "message": "docs: remove outdated development plan document\n\n- Deleted DEVELOPMENT_PLAN.md containing obsolete roadmap and feature plans\n- Removed documentation for completed features (batch testing, URL templates, interface-level TPS stats, Web UI)\n- Cleaned up outdated technical debt, timeline, and release planning sections",
          "timestamp": "2025-11-26T00:45:40+08:00",
          "tree_id": "fdf7e683b8a55d6ddf8840fd16491617479815ca",
          "url": "https://github.com/antlabs/gurl/commit/c629e99897d430933065853a601399377183907e"
        },
        "date": 1764089180759,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100060114,
            "unit": "ns/op\t 5172271 B/op\t   63795 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100060114,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5172271,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 63795,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 405.3,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "15032593 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 405.3,
            "unit": "ns/op",
            "extra": "15032593 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "15032593 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "15032593 times\n4 procs"
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
          "id": "1c7d0955b6a6804672274f4ab6b4efdea1f669d2",
          "message": "refactor: rename TotalBytes to ReadBytes for clarity in traffic statistics\n\n- Renamed TotalBytes field to ReadBytes in EndpointStats struct to distinguish from WriteBytes\n- Updated totalBytes to totalReadBytes in Results struct for consistency\n- Changed output labels from \"Data:\" to \"Read:\" in endpoint statistics display\n- Added latency percentiles output (p50, p75, p90, p95, p99) in PrintResults\n- Modified API handler and all related functions to use ReadBytes instead of TotalBytes",
          "timestamp": "2025-11-27T00:16:21+08:00",
          "tree_id": "b1dfe9c992f3465756f433663a2d3dbff54f860e",
          "url": "https://github.com/antlabs/gurl/commit/1c7d0955b6a6804672274f4ab6b4efdea1f669d2"
        },
        "date": 1764256311121,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100117052,
            "unit": "ns/op\t 7361822 B/op\t   90820 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100117052,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 7361822,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 90820,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 422,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14219050 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 422,
            "unit": "ns/op",
            "extra": "14219050 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14219050 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14219050 times\n4 procs"
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
          "id": "5cd9fb67194cef8d0b852824d0bdc53b9b0c03ed",
          "message": "feat: add scheduled execution support with --schedule-cron flag for all run modes\n\n- Added --schedule-cron flag to accept cron expressions for scheduled benchmark runs\n- Implemented runBenchmarkWithCron, runBatchTestWithCron, and runCompareWithCron functions\n- Added signal handling (SIGINT/SIGTERM) to gracefully stop scheduled runs\n- Integrated scheduler.ParseDailyCron to parse cron expressions and calculate next run times\n- Modified main function to check for ScheduleCron flag and route to appropriate schedule",
          "timestamp": "2025-11-28T00:06:57+08:00",
          "tree_id": "732c6a892eb70915ca6ec95df993a476e98bf07c",
          "url": "https://github.com/antlabs/gurl/commit/5cd9fb67194cef8d0b852824d0bdc53b9b0c03ed"
        },
        "date": 1764259654245,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkNetHTTPClient",
            "value": 100073645,
            "unit": "ns/op\t 5240975 B/op\t   64676 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - ns/op",
            "value": 100073645,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - B/op",
            "value": 5240975,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkNetHTTPClient - allocs/op",
            "value": 64676,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing",
            "value": 401.6,
            "unit": "ns/op\t     880 B/op\t       5 allocs/op",
            "extra": "14924128 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - ns/op",
            "value": 401.6,
            "unit": "ns/op",
            "extra": "14924128 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - B/op",
            "value": 880,
            "unit": "B/op",
            "extra": "14924128 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPRequestParsing - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14924128 times\n4 procs"
          }
        ]
      }
    ]
  }
}
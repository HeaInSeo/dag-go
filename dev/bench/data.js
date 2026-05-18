window.BENCHMARK_DATA = {
  "lastUpdate": 1779106715921,
  "repoUrl": "https://github.com/HeaInSeo/dag-go",
  "entries": {
    "dag-go benchmarks": [
      {
        "commit": {
          "author": {
            "name": "HeaInSeo",
            "username": "icgseoy",
            "email": "seoyhaein@gmail.com"
          },
          "committer": {
            "name": "HeaInSeo",
            "username": "icgseoy",
            "email": "seoyhaein@gmail.com"
          },
          "id": "252a9d9d532e9d50ee14bdf67db84fe88ffcc698",
          "message": "chore(ci): Node.js 24 opt-in (enforcement 2026-06-02)\n\nFORCE_JAVASCRIPT_ACTIONS_TO_NODE24=true 를 4개 workflow에 추가.\nactions/checkout@v4, setup-go@v5, upload/download-artifact@v4 모두 적용됨.\n\nCo-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>",
          "timestamp": "2026-05-18T12:14:40Z",
          "url": "https://github.com/HeaInSeo/dag-go/commit/252a9d9d532e9d50ee14bdf67db84fe88ffcc698"
        },
        "date": 1779106714993,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 3625,
            "unit": "ns/op\t    4624 B/op\t      57 allocs/op",
            "extra": "1206022 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 3625,
            "unit": "ns/op",
            "extra": "1206022 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4624,
            "unit": "B/op",
            "extra": "1206022 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 57,
            "unit": "allocs/op",
            "extra": "1206022 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 128803,
            "unit": "ns/op\t  135864 B/op\t    1812 allocs/op",
            "extra": "27982 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 128803,
            "unit": "ns/op",
            "extra": "27982 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 135864,
            "unit": "B/op",
            "extra": "27982 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "27982 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 5417261,
            "unit": "ns/op\t 3927851 B/op\t   53078 allocs/op",
            "extra": "676 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 5417261,
            "unit": "ns/op",
            "extra": "676 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3927851,
            "unit": "B/op",
            "extra": "676 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 53078,
            "unit": "allocs/op",
            "extra": "676 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 1317,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2911456 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 1317,
            "unit": "ns/op",
            "extra": "2911456 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2911456 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2911456 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 35841,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100740 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 35841,
            "unit": "ns/op",
            "extra": "100740 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100740 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100740 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 876798,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4018 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 876798,
            "unit": "ns/op",
            "extra": "4018 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4018 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4018 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 5499,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "638094 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 5499,
            "unit": "ns/op",
            "extra": "638094 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "638094 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "638094 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 19024,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "190082 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 19024,
            "unit": "ns/op",
            "extra": "190082 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "190082 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "190082 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 25156,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "146767 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 25156,
            "unit": "ns/op",
            "extra": "146767 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "146767 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "146767 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 13052,
            "unit": "ns/op\t    4355 B/op\t      83 allocs/op",
            "extra": "268804 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 13052,
            "unit": "ns/op",
            "extra": "268804 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 4355,
            "unit": "B/op",
            "extra": "268804 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 83,
            "unit": "allocs/op",
            "extra": "268804 times\n4 procs"
          }
        ]
      }
    ]
  }
}
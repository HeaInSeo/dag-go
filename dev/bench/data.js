window.BENCHMARK_DATA = {
  "lastUpdate": 1785134777900,
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
      },
      {
        "commit": {
          "author": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "committer": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "distinct": true,
          "id": "0d484be2f26f376d46d30b5d9bffe3fe8b8bc6b6",
          "message": "perf(bench): 베이스라인 재설정 및 bench.yml push-to-main 자동 트리거\n\n- bench_compare.sh 베이스라인: 2회 실행 평균값으로 교체\n  - PreFlight: 27685→18330 ns/op (ErrDependencyBlocked 센티넬로 -41% 실질 개선)\n  - PreFlight allocs: 89→43 (fmt.Errorf 제거)\n  - ToMermaid/CopyDag/DetectCycle: 현재 환경 측정값 반영\n- 기본 임계치: 10%→15% (공유 서버 스케줄러 노이즈 허용 마진)\n- NOISY_BENCHMARKS에 CopyDag_Small, ToMermaid_Small 추가 (절대값 <10μs)\n- bench.yml: workflow_dispatch 전용→push[main]+workflow_dispatch 자동 트리거\n  fail-on-alert: false 유지 (노이즈 오탐이 CI 차단 방지)\n\nCo-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>",
          "timestamp": "2026-05-25T16:42:09+09:00",
          "tree_id": "e0eaec714e15e182f58eae01713a21d398a0f1b0",
          "url": "https://github.com/HeaInSeo/dag-go/commit/0d484be2f26f376d46d30b5d9bffe3fe8b8bc6b6"
        },
        "date": 1779694999023,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 3397,
            "unit": "ns/op\t    4872 B/op\t      62 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 3397,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4872,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 62,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 112066,
            "unit": "ns/op\t  133264 B/op\t    1774 allocs/op",
            "extra": "30030 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 112066,
            "unit": "ns/op",
            "extra": "30030 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 133264,
            "unit": "B/op",
            "extra": "30030 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1774,
            "unit": "allocs/op",
            "extra": "30030 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 5439844,
            "unit": "ns/op\t 3907362 B/op\t   52772 allocs/op",
            "extra": "672 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 5439844,
            "unit": "ns/op",
            "extra": "672 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3907362,
            "unit": "B/op",
            "extra": "672 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 52772,
            "unit": "allocs/op",
            "extra": "672 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 1101,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3356235 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 1101,
            "unit": "ns/op",
            "extra": "3356235 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3356235 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3356235 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 30566,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "114686 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 30566,
            "unit": "ns/op",
            "extra": "114686 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "114686 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "114686 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 737512,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4860 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 737512,
            "unit": "ns/op",
            "extra": "4860 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4860 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4860 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 5534,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "628825 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 5534,
            "unit": "ns/op",
            "extra": "628825 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "628825 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "628825 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 17410,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "207348 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 17410,
            "unit": "ns/op",
            "extra": "207348 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "207348 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "207348 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 22614,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "154425 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 22614,
            "unit": "ns/op",
            "extra": "154425 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "154425 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "154425 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 9487,
            "unit": "ns/op\t    1873 B/op\t      43 allocs/op",
            "extra": "376701 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 9487,
            "unit": "ns/op",
            "extra": "376701 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 1873,
            "unit": "B/op",
            "extra": "376701 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 43,
            "unit": "allocs/op",
            "extra": "376701 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "committer": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "distinct": true,
          "id": "fc76263a0e1a5c82bbb32f00fedf4c69030c7fe7",
          "message": "docs: PROGRESS_LOG Stage 16 기록 및 아키텍처 노트 수정 (v1.2.1 릴리스 마감)\n\n- Stage 16 완료 항목 추가: 버그 6개, CI 자동화, 커버리지 93.0%\n- 아키텍처 노트 오류 수정: preFlight goroutine bounding\n  eg.SetLimit(10)+TryGo → eg.Go unbounded (Stage 13 수정 반영)\n- ErrDependencyBlocked, DependencySkipped 아키텍처 노트 추가\n- Known Follow-ups 섹션 추가: bTimeout, Node.js 20 deprecation\n- 헤더 업데이트: v1.2.1 / 2026-05-25 / 93.0%\n\nCo-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>",
          "timestamp": "2026-05-25T16:56:09+09:00",
          "tree_id": "08a72f40667aad6bd745554c7e54731590facb3d",
          "url": "https://github.com/HeaInSeo/dag-go/commit/fc76263a0e1a5c82bbb32f00fedf4c69030c7fe7"
        },
        "date": 1779695856893,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 3209,
            "unit": "ns/op\t    4272 B/op\t      53 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 3209,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4272,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 53,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 123196,
            "unit": "ns/op\t  134488 B/op\t    1791 allocs/op",
            "extra": "28876 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 123196,
            "unit": "ns/op",
            "extra": "28876 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 134488,
            "unit": "B/op",
            "extra": "28876 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1791,
            "unit": "allocs/op",
            "extra": "28876 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 5110456,
            "unit": "ns/op\t 3928714 B/op\t   53116 allocs/op",
            "extra": "711 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 5110456,
            "unit": "ns/op",
            "extra": "711 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3928714,
            "unit": "B/op",
            "extra": "711 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 53116,
            "unit": "allocs/op",
            "extra": "711 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 1262,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2924610 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 1262,
            "unit": "ns/op",
            "extra": "2924610 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2924610 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2924610 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 34004,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "104065 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 34004,
            "unit": "ns/op",
            "extra": "104065 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "104065 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "104065 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 863250,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 863250,
            "unit": "ns/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 5410,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "648537 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 5410,
            "unit": "ns/op",
            "extra": "648537 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "648537 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "648537 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 18040,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "199395 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 18040,
            "unit": "ns/op",
            "extra": "199395 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "199395 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "199395 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 23614,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "152949 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 23614,
            "unit": "ns/op",
            "extra": "152949 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "152949 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "152949 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 8217,
            "unit": "ns/op\t    1873 B/op\t      43 allocs/op",
            "extra": "432680 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 8217,
            "unit": "ns/op",
            "extra": "432680 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 1873,
            "unit": "B/op",
            "extra": "432680 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 43,
            "unit": "allocs/op",
            "extra": "432680 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "committer": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "distinct": true,
          "id": "a20d4d64627af39ca94ad85a80f52e5c51930343",
          "message": "Add CodeQL advanced setup",
          "timestamp": "2026-07-16T19:09:57+09:00",
          "tree_id": "79ca102b157dd7f4e44eae4d8b5f5c206b25550f",
          "url": "https://github.com/HeaInSeo/dag-go/commit/a20d4d64627af39ca94ad85a80f52e5c51930343"
        },
        "date": 1784196662233,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 2196,
            "unit": "ns/op\t    4632 B/op\t      58 allocs/op",
            "extra": "1714924 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 2196,
            "unit": "ns/op",
            "extra": "1714924 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4632,
            "unit": "B/op",
            "extra": "1714924 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 58,
            "unit": "allocs/op",
            "extra": "1714924 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 77012,
            "unit": "ns/op\t  133904 B/op\t    1783 allocs/op",
            "extra": "47580 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 77012,
            "unit": "ns/op",
            "extra": "47580 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 133904,
            "unit": "B/op",
            "extra": "47580 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1783,
            "unit": "allocs/op",
            "extra": "47580 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 3760596,
            "unit": "ns/op\t 3917858 B/op\t   52941 allocs/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 3760596,
            "unit": "ns/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3917858,
            "unit": "B/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 52941,
            "unit": "allocs/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 843,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4397874 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 843,
            "unit": "ns/op",
            "extra": "4397874 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4397874 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4397874 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 20948,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "161551 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 20948,
            "unit": "ns/op",
            "extra": "161551 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "161551 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "161551 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 452758,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8397 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 452758,
            "unit": "ns/op",
            "extra": "8397 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8397 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8397 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 3725,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "910712 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 3725,
            "unit": "ns/op",
            "extra": "910712 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "910712 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "910712 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 11642,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "304603 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 11642,
            "unit": "ns/op",
            "extra": "304603 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "304603 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "304603 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 15304,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "234517 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 15304,
            "unit": "ns/op",
            "extra": "234517 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "234517 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "234517 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 6234,
            "unit": "ns/op\t    1873 B/op\t      43 allocs/op",
            "extra": "563793 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 6234,
            "unit": "ns/op",
            "extra": "563793 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 1873,
            "unit": "B/op",
            "extra": "563793 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 43,
            "unit": "allocs/op",
            "extra": "563793 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "126678443+HeaInSeo@users.noreply.github.com",
            "name": "HeaInSeo",
            "username": "HeaInSeo"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0ffc259c8e986cbf2d2f247f12cabcf93e051942",
          "message": "fix(security): wire lint-security/vuln hard gates into CI (#2)\n\nFull CI audit: dag-go had zero automatic security scanning (same gap\nas the rest of this sweep) — gosec/govulncheck only in\nsecurity-observe.yml (workflow_dispatch only). Added\nlint-security-check/vuln-check (hard-failing) and wired both into\ntest.yml.\n\ngovulncheck: 0 findings even on go.mod's declared go1.25.5 — this\nlibrary doesn't touch the stdlib packages (crypto/tls, net/url, etc.)\ncarrying reachable CVEs elsewhere in this sweep, so no go.mod bump\nneeded here.\n\ngosec found 4 findings, all false positives, all fixed:\n- 3x G104 \"errors unhandled\" in dag_test.go on Close() calls that\n  already had //nolint:errcheck (a different linter than gosec's\n  G104) — extended to //nolint:errcheck,gosec across all 64 such call\n  sites in the file, not just the ones flagged in this run (same\n  pattern repeats throughout the file; fixing only the reported subset\n  would just surface the rest on the next lint run).\n- 1x G404 \"weak random number generator\" in node_bench_test.go — a\n  math/rand jitter for a benchmark's simulated async delay, not\n  security-sensitive.\n\nVerified: go build/vet/test, make lint (0 issues), lint-security-check\n(0 issues), vuln-check (0 findings) all clean.\n\nCo-authored-by: HeaInSeo <seoyhaein@gmail.com>",
          "timestamp": "2026-07-27T15:32:25+09:00",
          "tree_id": "bac1e3ad425620f0e93a1974345b2d8ff9c41f77",
          "url": "https://github.com/HeaInSeo/dag-go/commit/0ffc259c8e986cbf2d2f247f12cabcf93e051942"
        },
        "date": 1785134016404,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 3215,
            "unit": "ns/op\t    4264 B/op\t      50 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 3215,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4264,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 124225,
            "unit": "ns/op\t  133936 B/op\t    1781 allocs/op",
            "extra": "29420 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 124225,
            "unit": "ns/op",
            "extra": "29420 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 133936,
            "unit": "B/op",
            "extra": "29420 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1781,
            "unit": "allocs/op",
            "extra": "29420 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 5267051,
            "unit": "ns/op\t 3925154 B/op\t   53070 allocs/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 5267051,
            "unit": "ns/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3925154,
            "unit": "B/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 53070,
            "unit": "allocs/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 1272,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2562478 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 1272,
            "unit": "ns/op",
            "extra": "2562478 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2562478 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2562478 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 34073,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "105206 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 34073,
            "unit": "ns/op",
            "extra": "105206 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "105206 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "105206 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 861247,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4173 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 861247,
            "unit": "ns/op",
            "extra": "4173 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4173 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4173 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 5521,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "634628 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 5521,
            "unit": "ns/op",
            "extra": "634628 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "634628 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "634628 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 18364,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "191824 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 18364,
            "unit": "ns/op",
            "extra": "191824 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "191824 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "191824 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 23899,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "150657 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 23899,
            "unit": "ns/op",
            "extra": "150657 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "150657 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "150657 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 8637,
            "unit": "ns/op\t    1873 B/op\t      43 allocs/op",
            "extra": "402796 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 8637,
            "unit": "ns/op",
            "extra": "402796 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 1873,
            "unit": "B/op",
            "extra": "402796 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 43,
            "unit": "allocs/op",
            "extra": "402796 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "committer": {
            "email": "seoyhaein@gmail.com",
            "name": "HeaInSeo",
            "username": "icgseoy"
          },
          "distinct": true,
          "id": "79e36562a30d695c9f7d1bd2b8fc4d731f5f5490",
          "message": "fix(ci): rename lint job to disambiguate from test.yml's identical 'Run on Ubuntu' check name\n\nBoth lint.yml and test.yml used the job display name 'Run on Ubuntu',\nproducing two check-runs with the identical name on every PR. A\nbranch ruleset's required status check matches by name only, so\nrequiring 'Run on Ubuntu' would be satisfiable by either job\nsucceeding — not a real guarantee that both lint and tests passed.",
          "timestamp": "2026-07-27T15:45:08+09:00",
          "tree_id": "440c7caa87829f795af059249a68f8d76d56047c",
          "url": "https://github.com/HeaInSeo/dag-go/commit/79e36562a30d695c9f7d1bd2b8fc4d731f5f5490"
        },
        "date": 1785134777601,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCopyDag_Small",
            "value": 3224,
            "unit": "ns/op\t    4336 B/op\t      52 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - ns/op",
            "value": 3224,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - B/op",
            "value": 4336,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Small - allocs/op",
            "value": 52,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium",
            "value": 129080,
            "unit": "ns/op\t  140024 B/op\t    1858 allocs/op",
            "extra": "28257 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - ns/op",
            "value": 129080,
            "unit": "ns/op",
            "extra": "28257 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - B/op",
            "value": 140024,
            "unit": "B/op",
            "extra": "28257 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Medium - allocs/op",
            "value": 1858,
            "unit": "allocs/op",
            "extra": "28257 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large",
            "value": 5278231,
            "unit": "ns/op\t 3912506 B/op\t   52868 allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - ns/op",
            "value": 5278231,
            "unit": "ns/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - B/op",
            "value": 3912506,
            "unit": "B/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkCopyDag_Large - allocs/op",
            "value": 52868,
            "unit": "allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small",
            "value": 1365,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2782957 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - ns/op",
            "value": 1365,
            "unit": "ns/op",
            "extra": "2782957 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2782957 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Small - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2782957 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium",
            "value": 34195,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100614 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - ns/op",
            "value": 34195,
            "unit": "ns/op",
            "extra": "100614 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100614 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Medium - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100614 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large",
            "value": 861200,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4178 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - ns/op",
            "value": 861200,
            "unit": "ns/op",
            "extra": "4178 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4178 times\n4 procs"
          },
          {
            "name": "BenchmarkDetectCycle_Large - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4178 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small",
            "value": 5470,
            "unit": "ns/op\t    2528 B/op\t      51 allocs/op",
            "extra": "642318 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - ns/op",
            "value": 5470,
            "unit": "ns/op",
            "extra": "642318 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "642318 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Small - allocs/op",
            "value": 51,
            "unit": "allocs/op",
            "extra": "642318 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium",
            "value": 18122,
            "unit": "ns/op\t   12178 B/op\t     172 allocs/op",
            "extra": "199066 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - ns/op",
            "value": 18122,
            "unit": "ns/op",
            "extra": "199066 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - B/op",
            "value": 12178,
            "unit": "B/op",
            "extra": "199066 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Medium - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "199066 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large",
            "value": 23484,
            "unit": "ns/op\t   15587 B/op\t     215 allocs/op",
            "extra": "151929 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - ns/op",
            "value": 23484,
            "unit": "ns/op",
            "extra": "151929 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - B/op",
            "value": 15587,
            "unit": "B/op",
            "extra": "151929 times\n4 procs"
          },
          {
            "name": "BenchmarkToMermaid_Large - allocs/op",
            "value": 215,
            "unit": "allocs/op",
            "extra": "151929 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight",
            "value": 8751,
            "unit": "ns/op\t    1873 B/op\t      43 allocs/op",
            "extra": "403195 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - ns/op",
            "value": 8751,
            "unit": "ns/op",
            "extra": "403195 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - B/op",
            "value": 1873,
            "unit": "B/op",
            "extra": "403195 times\n4 procs"
          },
          {
            "name": "BenchmarkPreFlight - allocs/op",
            "value": 43,
            "unit": "allocs/op",
            "extra": "403195 times\n4 procs"
          }
        ]
      }
    ]
  }
}
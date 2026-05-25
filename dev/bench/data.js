window.BENCHMARK_DATA = {
  "lastUpdate": 1779695857625,
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
      }
    ]
  }
}
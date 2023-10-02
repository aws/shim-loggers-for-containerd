window.BENCHMARK_DATA = {
  "lastUpdate": 1696269673844,
  "repoUrl": "https://github.com/aws/shim-loggers-for-containerd",
  "entries": {
    "Benchmark for awslogs": [
      {
        "commit": {
          "author": {
            "email": "ningziwe@amazon.com",
            "name": "Ziwen Ning",
            "username": "ningziwen"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d03a1489d3c2d44645260123e40e106ea89a4977",
          "message": "ci: add basic benchmarking (#82)\n\n*Description of changes:*\r\n\r\nadd benchmarking for time and memory of sending a 1MB log\r\n\r\nResult page:\r\nhttps://aws.github.io/shim-loggers-for-containerd/dev/bench/\r\n\r\nUX improvements will come later in the gh-pages branch.\r\n\r\n\r\nBy submitting this pull request, I confirm that you can use, modify,\r\ncopy, and redistribute this contribution, under the terms of your\r\nchoice.\r\n\r\nSigned-off-by: Ziwen Ning <ningziwe@amazon.com>",
          "timestamp": "2023-09-29T15:15:33-07:00",
          "tree_id": "b2ed61664e0b68613dc0ae063baafe02520c112b",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/d03a1489d3c2d44645260123e40e106ea89a4977"
        },
        "date": 1696025833179,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFluentd - ns/op",
            "value": 9327039149,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkFluentd - B/op",
            "value": 40993112,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkFluentd - allocs/op",
            "value": 73526,
            "unit": "allocs/op",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ningziwe@amazon.com",
            "name": "Ziwen Ning",
            "username": "ningziwen"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0fe081af7810db59dd851915ee82e9604e725b34",
          "message": "ci: add lint for e2e and unit (#84)\n\n*Issue #, if available:*\r\n\r\n*Description of changes:*\r\n\r\n\r\nBy submitting this pull request, I confirm that you can use, modify,\r\ncopy, and redistribute this contribution, under the terms of your\r\nchoice.\r\n\r\nSigned-off-by: Ziwen Ning <ningziwe@amazon.com>",
          "timestamp": "2023-10-02T10:59:37-07:00",
          "tree_id": "0ea6ef376743fc77631ee2a4d138b5ef57a2d3f2",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/0fe081af7810db59dd851915ee82e9604e725b34"
        },
        "date": 1696269673317,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFluentd - ns/op",
            "value": 8880350431,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkFluentd - B/op",
            "value": 43102688,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkFluentd - allocs/op",
            "value": 73491,
            "unit": "allocs/op",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}
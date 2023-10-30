window.BENCHMARK_DATA = {
  "lastUpdate": 1698635157791,
  "repoUrl": "https://github.com/aws/shim-loggers-for-containerd",
  "entries": {
    "Benchmark for splunk": [
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
        "date": 1696025992736,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 9691451447,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42110712,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 73478,
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
        "date": 1696269784173,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8058176551,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42075216,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 73364,
            "unit": "allocs/op",
            "extra": "1 times\n2 procs"
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
          "id": "d2f8bb4879d80c164283cf75f7233d60934e264e",
          "message": "chore(deps): bump golang.org/x/net from 0.10.0 to 0.17.0 (#85)\n\nBumps [golang.org/x/net](https://github.com/golang/net) from 0.10.0 to\r\n0.17.0.\r\n<details>\r\n<summary>Commits</summary>\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/b225e7ca6dde1ef5a5ae5ce922861bda011cfabd\"><code>b225e7c</code></a>\r\nhttp2: limit maximum handler goroutines to MaxConcurrentStreams</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/88194ad8ab44a02ea952c169883c3f57db6cf9f4\"><code>88194ad</code></a>\r\ngo.mod: update golang.org/x dependencies</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/2b60a61f1e4cf3a5ecded0bd7e77ea168289e6de\"><code>2b60a61</code></a>\r\nquic: fix several bugs in flow control accounting</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/73d82efb96cacc0c378bc150b56675fc191894b9\"><code>73d82ef</code></a>\r\nquic: handle DATA_BLOCKED frames</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/5d5a036a503f8accd748f7453c0162115187be13\"><code>5d5a036</code></a>\r\nquic: handle streams moving from the data queue to the meta queue</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/350aad2603e57013fafb1a9e2089a382fe67dc80\"><code>350aad2</code></a>\r\nquic: correctly extend peer's flow control window after MAX_DATA</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/21814e71db756f39b69fb1a3e06350fa555a79b1\"><code>21814e7</code></a>\r\nquic: validate connection id transport parameters</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/a600b3518eed7a9a4e24380b4b249cb986d9b64d\"><code>a600b35</code></a>\r\nquic: avoid redundant MAX_DATA updates</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/ea633599b58dc6a50d33c7f5438edfaa8bc313df\"><code>ea63359</code></a>\r\nhttp2: check stream body is present on read timeout</li>\r\n<li><a\r\nhref=\"https://github.com/golang/net/commit/ddd8598e5694aa5e966e44573a53e895f6fa5eb2\"><code>ddd8598</code></a>\r\nquic: version negotiation</li>\r\n<li>Additional commits viewable in <a\r\nhref=\"https://github.com/golang/net/compare/v0.10.0...v0.17.0\">compare\r\nview</a></li>\r\n</ul>\r\n</details>\r\n<br />\r\n\r\n\r\n[![Dependabot compatibility\r\nscore](https://dependabot-badges.githubapp.com/badges/compatibility_score?dependency-name=golang.org/x/net&package-manager=go_modules&previous-version=0.10.0&new-version=0.17.0)](https://docs.github.com/en/github/managing-security-vulnerabilities/about-dependabot-security-updates#about-compatibility-scores)\r\n\r\nDependabot will resolve any conflicts with this PR as long as you don't\r\nalter it yourself. You can also trigger a rebase manually by commenting\r\n`@dependabot rebase`.\r\n\r\n[//]: # (dependabot-automerge-start)\r\n[//]: # (dependabot-automerge-end)\r\n\r\n---\r\n\r\n<details>\r\n<summary>Dependabot commands and options</summary>\r\n<br />\r\n\r\nYou can trigger Dependabot actions by commenting on this PR:\r\n- `@dependabot rebase` will rebase this PR\r\n- `@dependabot recreate` will recreate this PR, overwriting any edits\r\nthat have been made to it\r\n- `@dependabot merge` will merge this PR after your CI passes on it\r\n- `@dependabot squash and merge` will squash and merge this PR after\r\nyour CI passes on it\r\n- `@dependabot cancel merge` will cancel a previously requested merge\r\nand block automerging\r\n- `@dependabot reopen` will reopen this PR if it is closed\r\n- `@dependabot close` will close this PR and stop Dependabot recreating\r\nit. You can achieve the same result by closing it manually\r\n- `@dependabot show <dependency name> ignore conditions` will show all\r\nof the ignore conditions of the specified dependency\r\n- `@dependabot ignore this major version` will close this PR and stop\r\nDependabot creating any more for this major version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this minor version` will close this PR and stop\r\nDependabot creating any more for this minor version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this dependency` will close this PR and stop\r\nDependabot creating any more for this dependency (unless you reopen the\r\nPR or upgrade to it yourself)\r\nYou can disable automated security fix PRs for this repo from the\r\n[Security Alerts\r\npage](https://github.com/aws/shim-loggers-for-containerd/network/alerts).\r\n\r\n</details>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2023-10-29T20:02:20-07:00",
          "tree_id": "e965326ced4128febe38348f47763884ddf45baf",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/d2f8bb4879d80c164283cf75f7233d60934e264e"
        },
        "date": 1698635156986,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8571449964,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42105112,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 73533,
            "unit": "allocs/op",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}
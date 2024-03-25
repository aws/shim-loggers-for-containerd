window.BENCHMARK_DATA = {
  "lastUpdate": 1711381598413,
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
          "id": "53b6d6b72a4e19bc1dc236c8a753b9c68fdec750",
          "message": "chore(deps): bump google.golang.org/grpc from 1.53.0 to 1.56.3 (#86)\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from\r\n1.53.0 to 1.56.3.\r\n<details>\r\n<summary>Release notes</summary>\r\n<p><em>Sourced from <a\r\nhref=\"https://github.com/grpc/grpc-go/releases\">google.golang.org/grpc's\r\nreleases</a>.</em></p>\r\n<blockquote>\r\n<h2>Release 1.56.3</h2>\r\n<h1>Security</h1>\r\n<ul>\r\n<li>\r\n<p>server: prohibit more than MaxConcurrentStreams handlers from running\r\nat once (CVE-2023-44487)</p>\r\n<p>In addition to this change, applications should ensure they do not\r\nleave running tasks behind related to the RPC before returning from\r\nmethod handlers, or should enforce appropriate limits on any such\r\nwork.</p>\r\n</li>\r\n</ul>\r\n<h2>Release 1.56.2</h2>\r\n<ul>\r\n<li>status: To fix a panic, <code>status.FromError</code> now returns an\r\nerror with <code>codes.Unknown</code> when the error implements the\r\n<code>GRPCStatus()</code> method, and calling <code>GRPCStatus()</code>\r\nreturns <code>nil</code>. (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6374\">#6374</a>)</li>\r\n</ul>\r\n<h2>Release 1.56.1</h2>\r\n<ul>\r\n<li>client: handle empty address lists correctly in\r\naddrConn.updateAddrs</li>\r\n</ul>\r\n<h2>Release 1.56.0</h2>\r\n<h1>New Features</h1>\r\n<ul>\r\n<li>client: support channel idleness using <code>WithIdleTimeout</code>\r\ndial option (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6263\">#6263</a>)\r\n<ul>\r\n<li>This feature is currently disabled by default, but will be enabled\r\nwith a 30 minute default in the future.</li>\r\n</ul>\r\n</li>\r\n<li>client: when using pickfirst, keep channel state in\r\nTRANSIENT_FAILURE until it becomes READY (<a\r\nhref=\"https://github.com/grpc/proposal/blob/master/A62-pick-first.md\">gRFC\r\nA62</a>) (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6306\">#6306</a>)</li>\r\n<li>xds: Add support for Custom LB Policies (<a\r\nhref=\"https://github.com/grpc/proposal/blob/master/A52-xds-custom-lb-policies.md\">gRFC\r\nA52</a>) (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6224\">#6224</a>)</li>\r\n<li>xds: support pick_first Custom LB policy (<a\r\nhref=\"https://github.com/grpc/proposal/blob/master/A62-pick-first.md\">gRFC\r\nA62</a>) (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6314\">#6314</a>)\r\n(<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6317\">#6317</a>)</li>\r\n<li>client: add support for pickfirst address shuffling (<a\r\nhref=\"https://github.com/grpc/proposal/blob/master/A62-pick-first.md\">gRFC\r\nA62</a>) (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6311\">#6311</a>)</li>\r\n<li>xds: Add support for String Matcher Header Matcher in RDS (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6313\">#6313</a>)</li>\r\n<li>xds/outlierdetection: Add Channelz Logger to Outlier Detection LB\r\n(<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6145\">#6145</a>)\r\n<ul>\r\n<li>Special Thanks: <a\r\nhref=\"https://github.com/s-matyukevich\"><code>@​s-matyukevich</code></a></li>\r\n</ul>\r\n</li>\r\n<li>xds: enable RLS in xDS by default (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6343\">#6343</a>)</li>\r\n<li>orca: add support for application_utilization field and missing\r\nrange checks on several metrics setters</li>\r\n<li>balancer/weightedroundrobin: add new LB policy for balancing between\r\nbackends based on their load reports (<a\r\nhref=\"https://github.com/grpc/proposal/blob/master/A58-client-side-weighted-round-robin-lb-policy.md\">gRFC\r\nA58</a>) (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6241\">#6241</a>)</li>\r\n<li>authz: add conversion of json to RBAC Audit Logging config (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6192\">#6192</a>)</li>\r\n<li>authz: add support for stdout logger (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6230\">#6230</a>\r\nand <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6298\">#6298</a>)</li>\r\n<li>authz: support customizable audit functionality for authorization\r\npolicy (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6192\">#6192</a> <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6230\">#6230</a> <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6298\">#6298</a> <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6158\">#6158</a> <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6304\">#6304</a>\r\nand <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6225\">#6225</a>)</li>\r\n</ul>\r\n<h1>Bug Fixes</h1>\r\n<ul>\r\n<li>orca: fix a race at startup of out-of-band metric subscriptions that\r\nwould cause the report interval to request 0 (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6245\">#6245</a>)</li>\r\n<li>xds/xdsresource: Fix Outlier Detection Config Handling and correctly\r\nset xDS Defaults (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6361\">#6361</a>)</li>\r\n<li>xds/outlierdetection: Fix Outlier Detection Config Handling by\r\nsetting defaults in ParseConfig() (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6361\">#6361</a>)</li>\r\n</ul>\r\n<h1>API Changes</h1>\r\n<ul>\r\n<li>orca: allow a ServerMetricsProvider to be passed to the ORCA service\r\nand ServerOption (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6223\">#6223</a>)</li>\r\n</ul>\r\n<h2>Release 1.55.1</h2>\r\n<ul>\r\n<li>status: To fix a panic, <code>status.FromError</code> now returns an\r\nerror with <code>codes.Unknown</code> when the error implements the\r\n<code>GRPCStatus()</code> method, and calling <code>GRPCStatus()</code>\r\nreturns <code>nil</code>. (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6374\">#6374</a>)</li>\r\n</ul>\r\n<h2>Release 1.55.0</h2>\r\n<h1>Behavior Changes</h1>\r\n<ul>\r\n<li>xds: enable federation support by default (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6151\">#6151</a>)</li>\r\n<li>status: <code>status.Code</code> and <code>status.FromError</code>\r\nhandle wrapped errors (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6031\">#6031</a>\r\nand <a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6150\">#6150</a>)</li>\r\n</ul>\r\n<!-- raw HTML omitted -->\r\n</blockquote>\r\n<p>... (truncated)</p>\r\n</details>\r\n<details>\r\n<summary>Commits</summary>\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/1055b481ed2204a29d233286b9b50c42b63f8825\"><code>1055b48</code></a>\r\nUpdate version.go to 1.56.3 (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6713\">#6713</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/5efd7bd73e11fea58d1c7f1c110902e78a286299\"><code>5efd7bd</code></a>\r\nserver: prohibit more than MaxConcurrentStreams handlers from running at\r\nonce...</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/bd1f038e7234580c2694e433bec5cd97e7b7f662\"><code>bd1f038</code></a>\r\nUpgrade version.go to 1.56.3-dev (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6434\">#6434</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/faab8736bf73291f92b867d5dae31c927d53d508\"><code>faab873</code></a>\r\nUpdate version.go to v1.56.2 (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6432\">#6432</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/6b0b291d79831b1c8caafceec268b82c92253f96\"><code>6b0b291</code></a>\r\nstatus: fix panic when servers return a wrapped error with status OK (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6374\">#6374</a>)\r\n...</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/ed56401aa514462d5371713b8ec5c889da33953c\"><code>ed56401</code></a>\r\n[PSM interop] Don't fail target if sub-target already failed (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6390\">#6390</a>)\r\n(<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6405\">#6405</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/cd6a794f0bdcf9a216e8f4d3c5717faf96d9fd78\"><code>cd6a794</code></a>\r\nUpdate version.go to v1.56.2-dev (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6387\">#6387</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/5b67e5ea449ef0686a0c0b6de48cd4cb63e3db2a\"><code>5b67e5e</code></a>\r\nUpdate version.go to v1.56.1 (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6386\">#6386</a>)</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/d0f5150384a87f9fcac488a9c18727a55b7354c1\"><code>d0f5150</code></a>\r\nclient: handle empty address lists correctly in addrConn.updateAddrs (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6354\">#6354</a>)\r\n...</li>\r\n<li><a\r\nhref=\"https://github.com/grpc/grpc-go/commit/997c1ea101cc5d496d2b148388f1df49632a9171\"><code>997c1ea</code></a>\r\nChange version to 1.56.1-dev (<a\r\nhref=\"https://redirect.github.com/grpc/grpc-go/issues/6345\">#6345</a>)</li>\r\n<li>Additional commits viewable in <a\r\nhref=\"https://github.com/grpc/grpc-go/compare/v1.53.0...v1.56.3\">compare\r\nview</a></li>\r\n</ul>\r\n</details>\r\n<br />\r\n\r\n\r\n[![Dependabot compatibility\r\nscore](https://dependabot-badges.githubapp.com/badges/compatibility_score?dependency-name=google.golang.org/grpc&package-manager=go_modules&previous-version=1.53.0&new-version=1.56.3)](https://docs.github.com/en/github/managing-security-vulnerabilities/about-dependabot-security-updates#about-compatibility-scores)\r\n\r\nDependabot will resolve any conflicts with this PR as long as you don't\r\nalter it yourself. You can also trigger a rebase manually by commenting\r\n`@dependabot rebase`.\r\n\r\n[//]: # (dependabot-automerge-start)\r\n[//]: # (dependabot-automerge-end)\r\n\r\n---\r\n\r\n<details>\r\n<summary>Dependabot commands and options</summary>\r\n<br />\r\n\r\nYou can trigger Dependabot actions by commenting on this PR:\r\n- `@dependabot rebase` will rebase this PR\r\n- `@dependabot recreate` will recreate this PR, overwriting any edits\r\nthat have been made to it\r\n- `@dependabot merge` will merge this PR after your CI passes on it\r\n- `@dependabot squash and merge` will squash and merge this PR after\r\nyour CI passes on it\r\n- `@dependabot cancel merge` will cancel a previously requested merge\r\nand block automerging\r\n- `@dependabot reopen` will reopen this PR if it is closed\r\n- `@dependabot close` will close this PR and stop Dependabot recreating\r\nit. You can achieve the same result by closing it manually\r\n- `@dependabot show <dependency name> ignore conditions` will show all\r\nof the ignore conditions of the specified dependency\r\n- `@dependabot ignore this major version` will close this PR and stop\r\nDependabot creating any more for this major version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this minor version` will close this PR and stop\r\nDependabot creating any more for this minor version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this dependency` will close this PR and stop\r\nDependabot creating any more for this dependency (unless you reopen the\r\nPR or upgrade to it yourself)\r\nYou can disable automated security fix PRs for this repo from the\r\n[Security Alerts\r\npage](https://github.com/aws/shim-loggers-for-containerd/network/alerts).\r\n\r\n</details>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2023-10-29T20:02:59-07:00",
          "tree_id": "06ebe0247e0691a199794baa4585130aa0abe2d9",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/53b6d6b72a4e19bc1dc236c8a753b9c68fdec750"
        },
        "date": 1698635183857,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7734681361,
            "unit": "ns/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 41800280,
            "unit": "B/op",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 73471,
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
          "id": "37a56b1cc2e3fd50fc8a3aa2424b48c3abadd7ff",
          "message": "chore(deps): bump google.golang.org/protobuf from 1.30.0 to 1.33.0 (#89)\n\nBumps google.golang.org/protobuf from 1.30.0 to 1.33.0.\r\n\r\n\r\n[![Dependabot compatibility\r\nscore](https://dependabot-badges.githubapp.com/badges/compatibility_score?dependency-name=google.golang.org/protobuf&package-manager=go_modules&previous-version=1.30.0&new-version=1.33.0)](https://docs.github.com/en/github/managing-security-vulnerabilities/about-dependabot-security-updates#about-compatibility-scores)\r\n\r\nDependabot will resolve any conflicts with this PR as long as you don't\r\nalter it yourself. You can also trigger a rebase manually by commenting\r\n`@dependabot rebase`.\r\n\r\n[//]: # (dependabot-automerge-start)\r\n[//]: # (dependabot-automerge-end)\r\n\r\n---\r\n\r\n<details>\r\n<summary>Dependabot commands and options</summary>\r\n<br />\r\n\r\nYou can trigger Dependabot actions by commenting on this PR:\r\n- `@dependabot rebase` will rebase this PR\r\n- `@dependabot recreate` will recreate this PR, overwriting any edits\r\nthat have been made to it\r\n- `@dependabot merge` will merge this PR after your CI passes on it\r\n- `@dependabot squash and merge` will squash and merge this PR after\r\nyour CI passes on it\r\n- `@dependabot cancel merge` will cancel a previously requested merge\r\nand block automerging\r\n- `@dependabot reopen` will reopen this PR if it is closed\r\n- `@dependabot close` will close this PR and stop Dependabot recreating\r\nit. You can achieve the same result by closing it manually\r\n- `@dependabot show <dependency name> ignore conditions` will show all\r\nof the ignore conditions of the specified dependency\r\n- `@dependabot ignore this major version` will close this PR and stop\r\nDependabot creating any more for this major version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this minor version` will close this PR and stop\r\nDependabot creating any more for this minor version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this dependency` will close this PR and stop\r\nDependabot creating any more for this dependency (unless you reopen the\r\nPR or upgrade to it yourself)\r\nYou can disable automated security fix PRs for this repo from the\r\n[Security Alerts\r\npage](https://github.com/aws/shim-loggers-for-containerd/network/alerts).\r\n\r\n</details>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-03-21T08:17:17-07:00",
          "tree_id": "5906b797643d30ba566c944f6750ae8565b8ada5",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/37a56b1cc2e3fd50fc8a3aa2424b48c3abadd7ff"
        },
        "date": 1711034430407,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8654735708,
            "unit": "ns/op\t40577072 B/op\t   66255 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8654735708,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 40577072,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66255,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
          "id": "2de4d68856005881922b0c450bc067c87d2552bd",
          "message": "chore(deps): bump github.com/opencontainers/runc from 1.1.5 to 1.1.12 (#88)\n\nBumps\r\n[github.com/opencontainers/runc](https://github.com/opencontainers/runc)\r\nfrom 1.1.5 to 1.1.12.\r\n<details>\r\n<summary>Release notes</summary>\r\n<p><em>Sourced from <a\r\nhref=\"https://github.com/opencontainers/runc/releases\">github.com/opencontainers/runc's\r\nreleases</a>.</em></p>\r\n<blockquote>\r\n<h2>runc 1.1.12 -- &quot;Now you're thinking with Portals™!&quot;</h2>\r\n<p>This is the twelfth patch release in the 1.1.z release branch of\r\nrunc.\r\nIt fixes a high-severity container breakout vulnerability involving\r\nleaked file descriptors, and users are strongly encouraged to update as\r\nsoon as possible.</p>\r\n<ul>\r\n<li>\r\n<p>Fix <a\r\nhref=\"https://github.com/opencontainers/runc/security/advisories/GHSA-xr7r-f8xq-vfvv\">CVE-2024-21626</a>,\r\na container breakout attack that took advantage of\r\na file descriptor that was leaked internally within runc (but never\r\nleaked to the container process).</p>\r\n<p>In addition to fixing the leak, several strict hardening measures\r\nwere\r\nadded to ensure that future internal leaks could not be used to break\r\nout in this manner again.</p>\r\n<p>Based on our research, while no other container runtime had a similar\r\nleak, none had any of the hardening steps we've introduced (and some\r\nruntimes would not check for any file descriptors that a calling\r\nprocess may have leaked to them, allowing for container breakouts due\r\nto basic user error).</p>\r\n</li>\r\n</ul>\r\n<h3>Static Linking Notices</h3>\r\n<p>The <code>runc</code> binary distributed with this release are\r\n<em>statically linked</em> with\r\nthe following <a\r\nhref=\"https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html\">GNU\r\nLGPL-2.1</a> licensed libraries, with <code>runc</code> acting\r\nas a &quot;work that uses the Library&quot;:</p>\r\n<ul>\r\n<li><a href=\"https://github.com/seccomp/libseccomp\">libseccomp</a></li>\r\n</ul>\r\n<p>The versions of these libraries were not modified from their upstream\r\nversions,\r\nbut in order to comply with the LGPL-2.1 (§6(a)), we have attached the\r\ncomplete source code for those libraries which (when combined with the\r\nattached\r\nrunc source code) may be used to exercise your rights under the\r\nLGPL-2.1.</p>\r\n<p>However we strongly suggest that you make use of your distribution's\r\npackages\r\nor download them from the authoritative upstream sources, especially\r\nsince\r\nthese libraries are related to the security of your containers.</p>\r\n<!-- raw HTML omitted -->\r\n<p>Thanks to all of the contributors who made this release possible:</p>\r\n<ul>\r\n<li>Aleksa Sarai <a\r\nhref=\"mailto:cyphar@cyphar.com\">cyphar@cyphar.com</a></li>\r\n<li>hang.jiang <a\r\nhref=\"mailto:hang.jiang@daocloud.io\">hang.jiang@daocloud.io</a></li>\r\n<li>lfbzhm <a\r\nhref=\"mailto:lifubang@acmcoder.com\">lifubang@acmcoder.com</a></li>\r\n</ul>\r\n<p>Signed-off-by: Aleksa Sarai <a\r\nhref=\"mailto:cyphar@cyphar.com\">cyphar@cyphar.com</a></p>\r\n<!-- raw HTML omitted -->\r\n</blockquote>\r\n<p>... (truncated)</p>\r\n</details>\r\n<details>\r\n<summary>Changelog</summary>\r\n<p><em>Sourced from <a\r\nhref=\"https://github.com/opencontainers/runc/blob/v1.1.12/CHANGELOG.md\">github.com/opencontainers/runc's\r\nchangelog</a>.</em></p>\r\n<blockquote>\r\n<h2>[1.1.12] - 2024-01-31</h2>\r\n<blockquote>\r\n<p>Now you're thinking with Portals™!</p>\r\n</blockquote>\r\n<h3>Security</h3>\r\n<ul>\r\n<li>Fix <a\r\nhref=\"https://github.com/opencontainers/runc/security/advisories/GHSA-xr7r-f8xq-vfvv\">CVE-2024-21626</a>,\r\na container breakout attack that took\r\nadvantage of a file descriptor that was leaked internally within runc\r\n(but\r\nnever leaked to the container process). In addition to fixing the leak,\r\nseveral strict hardening measures were added to ensure that future\r\ninternal\r\nleaks could not be used to break out in this manner again. Based on our\r\nresearch, while no other container runtime had a similar leak, none had\r\nany\r\nof the hardening steps we've introduced (and some runtimes would not\r\ncheck\r\nfor any file descriptors that a calling process may have leaked to them,\r\nallowing for container breakouts due to basic user error).</li>\r\n</ul>\r\n<h2>[1.1.11] - 2024-01-01</h2>\r\n<blockquote>\r\n<p>Happy New Year!</p>\r\n</blockquote>\r\n<h3>Fixed</h3>\r\n<ul>\r\n<li>Fix several issues with userns path handling. (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4122\">#4122</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4124\">#4124</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4134\">#4134</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4144\">#4144</a>)</li>\r\n</ul>\r\n<h3>Changed</h3>\r\n<ul>\r\n<li>Support memory.peak and memory.swap.peak in cgroups v2.\r\nAdd <code>swapOnlyUsage</code> in <code>MemoryStats</code>. This field\r\nreports swap-only usage.\r\nFor cgroupv1, <code>Usage</code> and <code>Failcnt</code> are set by\r\nsubtracting memory usage\r\nfrom memory+swap usage. For cgroupv2, <code>Usage</code>,\r\n<code>Limit</code>, and <code>MaxUsage</code>\r\nare set. (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4000\">#4000</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4010\">#4010</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4131\">#4131</a>)</li>\r\n<li>build(deps): bump github.com/cyphar/filepath-securejoin. (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4140\">#4140</a>)</li>\r\n</ul>\r\n<h2>[1.1.10] - 2023-10-31</h2>\r\n<blockquote>\r\n<p>Śruba, przykręcona we śnie, nie zmieni sytuacji, jaka panuje na\r\njawie.</p>\r\n</blockquote>\r\n<h3>Added</h3>\r\n<ul>\r\n<li>Support for <code>hugetlb.&lt;pagesize&gt;.rsvd</code> limiting and\r\naccounting. Fixes the\r\nissue of postres failing when hugepage limits are set. (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/3859\">#3859</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4077\">#4077</a>)</li>\r\n</ul>\r\n<h3>Fixed</h3>\r\n<ul>\r\n<li>Fixed permissions of a newly created directories to not depend on\r\nthe value\r\nof umask in tmpcopyup feature implementation. (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/3991\">#3991</a>,\r\n<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4060\">#4060</a>)</li>\r\n<li>libcontainer: cgroup v1 GetStats now ignores missing\r\n<code>kmem.limit_in_bytes</code>\r\n(fixes the compatibility with Linux kernel 6.1+). (<a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4028\">#4028</a>)</li>\r\n</ul>\r\n<!-- raw HTML omitted -->\r\n</blockquote>\r\n<p>... (truncated)</p>\r\n</details>\r\n<details>\r\n<summary>Commits</summary>\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/51d5e94601ceffbbd85688df1c928ecccbfa4685\"><code>51d5e94</code></a>\r\nVERSION: release 1.1.12</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/2a4ed3e75b9e80d93d1836a9c4c1ebfa2b78870e\"><code>2a4ed3e</code></a>\r\nmerge 1.1-ghsa-xr7r-f8xq-vfvv into release-1.1</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/e9665f4d606b64bf9c4652ab2510da368bfbd951\"><code>e9665f4</code></a>\r\ninit: don't special-case logrus fds</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/683ad2ff3b01fb142ece7a8b3829de17150cf688\"><code>683ad2f</code></a>\r\nlibcontainer: mark all non-stdio fds O_CLOEXEC before spawning init</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/b6633f48a8c970433737b9be5bfe4f25d58a5aa7\"><code>b6633f4</code></a>\r\ncgroup: plug leaks of /sys/fs/cgroup handle</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/284ba3057e428f8d6c7afcc3b0ac752e525957df\"><code>284ba30</code></a>\r\ninit: close internal fds before execve</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/fbe3eed1e568a376f371d2ced1b4ac16b7d7adde\"><code>fbe3eed</code></a>\r\nsetns init: do explicit lookup of execve argument early</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/0994249a5ec4e363bfcf9af58a87a722e9a3a31b\"><code>0994249</code></a>\r\ninit: verify after chdir that cwd is inside the container</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/506552a88bd3455e80a9b3829568e94ec0160309\"><code>506552a</code></a>\r\nFix File to Close</li>\r\n<li><a\r\nhref=\"https://github.com/opencontainers/runc/commit/099ff69336840fecf3fc0ab13aab4c3aded640c3\"><code>099ff69</code></a>\r\nmerge <a\r\nhref=\"https://redirect.github.com/opencontainers/runc/issues/4177\">#4177</a>\r\ninto opencontainers/runc:release-1.1</li>\r\n<li>Additional commits viewable in <a\r\nhref=\"https://github.com/opencontainers/runc/compare/v1.1.5...v1.1.12\">compare\r\nview</a></li>\r\n</ul>\r\n</details>\r\n<br />\r\n\r\n\r\n[![Dependabot compatibility\r\nscore](https://dependabot-badges.githubapp.com/badges/compatibility_score?dependency-name=github.com/opencontainers/runc&package-manager=go_modules&previous-version=1.1.5&new-version=1.1.12)](https://docs.github.com/en/github/managing-security-vulnerabilities/about-dependabot-security-updates#about-compatibility-scores)\r\n\r\nYou can trigger a rebase of this PR by commenting `@dependabot rebase`.\r\n\r\n[//]: # (dependabot-automerge-start)\r\n[//]: # (dependabot-automerge-end)\r\n\r\n---\r\n\r\n<details>\r\n<summary>Dependabot commands and options</summary>\r\n<br />\r\n\r\nYou can trigger Dependabot actions by commenting on this PR:\r\n- `@dependabot rebase` will rebase this PR\r\n- `@dependabot recreate` will recreate this PR, overwriting any edits\r\nthat have been made to it\r\n- `@dependabot merge` will merge this PR after your CI passes on it\r\n- `@dependabot squash and merge` will squash and merge this PR after\r\nyour CI passes on it\r\n- `@dependabot cancel merge` will cancel a previously requested merge\r\nand block automerging\r\n- `@dependabot reopen` will reopen this PR if it is closed\r\n- `@dependabot close` will close this PR and stop Dependabot recreating\r\nit. You can achieve the same result by closing it manually\r\n- `@dependabot show <dependency name> ignore conditions` will show all\r\nof the ignore conditions of the specified dependency\r\n- `@dependabot ignore this major version` will close this PR and stop\r\nDependabot creating any more for this major version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this minor version` will close this PR and stop\r\nDependabot creating any more for this minor version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this dependency` will close this PR and stop\r\nDependabot creating any more for this dependency (unless you reopen the\r\nPR or upgrade to it yourself)\r\nYou can disable automated security fix PRs for this repo from the\r\n[Security Alerts\r\npage](https://github.com/aws/shim-loggers-for-containerd/network/alerts).\r\n\r\n</details>\r\n\r\n> **Note**\r\n> Automatic rebases have been disabled on this pull request as it has\r\nbeen open for over 30 days.\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-03-21T08:24:42-07:00",
          "tree_id": "07495f9e35ce152564782912196d7381f22c6472",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/2de4d68856005881922b0c450bc067c87d2552bd"
        },
        "date": 1711034867472,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7914640257,
            "unit": "ns/op\t42057152 B/op\t   65609 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7914640257,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42057152,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65609,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
          "id": "d4df5dd67316689b02ca909baa796a5ba7dac30d",
          "message": "chore(deps): bump github.com/containerd/containerd from 1.6.18 to 1.6.26 (#87)\n\nBumps\r\n[github.com/containerd/containerd](https://github.com/containerd/containerd)\r\nfrom 1.6.18 to 1.6.26.\r\n<details>\r\n<summary>Release notes</summary>\r\n<p><em>Sourced from <a\r\nhref=\"https://github.com/containerd/containerd/releases\">github.com/containerd/containerd's\r\nreleases</a>.</em></p>\r\n<blockquote>\r\n<h2>containerd 1.6.26</h2>\r\n<p>Welcome to the v1.6.26 release of containerd!</p>\r\n<p>The twenty-sixth patch release for containerd 1.6 contains various\r\nfixes and updates.</p>\r\n<h3>Notable Updates</h3>\r\n<ul>\r\n<li><strong>Fix windows default path overwrite issue</strong> (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9441\">#9441</a>)</li>\r\n<li><strong>Update push to inherit distribution sources from\r\nparent</strong> (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9453\">#9453</a>)</li>\r\n<li><strong>Mask <code>/sys/devices/virtual/powercap</code> path in\r\nruntime spec and deny in default apparmor profile</strong> (<a\r\nhref=\"https://github.com/containerd/containerd/security/advisories/GHSA-7ww5-4wqc-m92c\">GHSA-7ww5-4wqc-m92c</a>)</li>\r\n</ul>\r\n<h3>Deprecation Warnings</h3>\r\n<ul>\r\n<li><strong>Emit deprecation warning for AUFS snapshotter usage</strong>\r\n(<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9448\">#9448</a>)</li>\r\n<li><strong>Emit deprecation warning for v1 runtime usage</strong> (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9468\">#9468</a>)</li>\r\n<li><strong>Emit deprecation warning for CRI v1alpha1 usage</strong> (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9468\">#9468</a>)</li>\r\n</ul>\r\n<p>See the changelog for complete list of changes</p>\r\n<p>Please try out the release binaries and report any issues at\r\n<a\r\nhref=\"https://github.com/containerd/containerd/issues\">https://github.com/containerd/containerd/issues</a>.</p>\r\n<h3>Contributors</h3>\r\n<ul>\r\n<li>Samuel Karp</li>\r\n<li>Derek McGowan</li>\r\n<li>Kohei Tokunaga</li>\r\n<li>Phil Estes</li>\r\n<li>Bjorn Neergaard</li>\r\n<li>Sebastiaan van Stijn</li>\r\n<li>Brian Goff</li>\r\n<li>Charity Kathure</li>\r\n<li>Kazuyoshi Kato</li>\r\n<li>Milas Bowman</li>\r\n<li>Wei Fu</li>\r\n<li>ruiwen-zhao</li>\r\n</ul>\r\n<h3>Changes</h3>\r\n<!-- raw HTML omitted -->\r\n<ul>\r\n<li>[release/1.6] Prepare release notes for v1.6.26 (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9490\">#9490</a>)\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/ac5c5d3e03ab3c5b8103a1c0bd9931389f7a8fcf\"><code>ac5c5d3e0</code></a>\r\nPrepare release notes for v1.6.26</li>\r\n</ul>\r\n</li>\r\n<li>Github Security Advisory <a\r\nhref=\"https://github.com/containerd/containerd/security/advisories/GHSA-7ww5-4wqc-m92c\">GHSA-7ww5-4wqc-m92c</a>\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/02f07fe1994a3ddda3626c1ede2e32bc82b8e426\"><code>02f07fe19</code></a>\r\ncontrib/apparmor: deny /sys/devices/virtual/powercap</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/c94577e78d2924ddeb90d1601e31b50ee3acac48\"><code>c94577e78</code></a>\r\noci/spec: deny /sys/devices/virtual/powercap</li>\r\n</ul>\r\n</li>\r\n<li>[release/1.6] update to go1.20.12, test go1.21.5 (<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9472\">#9472</a>)\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/7cbdfc92ef38f789f1a2773fa6fac405d361a6cc\"><code>7cbdfc92e</code></a>\r\nupdate to go1.20.12, test go1.21.5</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/024b1cce6b27f10e00bb9bde33a5fe9563545f8d\"><code>024b1cce6</code></a>\r\nupdate to go1.20.11, test go1.21.4</li>\r\n</ul>\r\n</li>\r\n<li>[release/1.6] Add cri-api v1alpha2 usage warning to all api calls\r\n(<a\r\nhref=\"https://redirect.github.com/containerd/containerd/pull/9484\">#9484</a>)</li>\r\n</ul>\r\n<!-- raw HTML omitted -->\r\n</blockquote>\r\n<p>... (truncated)</p>\r\n</details>\r\n<details>\r\n<summary>Commits</summary>\r\n<ul>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/3dd1e886e55dd695541fdcd67420c2888645a495\"><code>3dd1e88</code></a>\r\nMerge pull request <a\r\nhref=\"https://redirect.github.com/containerd/containerd/issues/9490\">#9490</a>\r\nfrom dmcgowan/prepare-1.6.26</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/746b910f05855c8bfdb4415a1c0f958b234910e5\"><code>746b910</code></a>\r\nMerge pull request from GHSA-7ww5-4wqc-m92c</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/ac5c5d3e03ab3c5b8103a1c0bd9931389f7a8fcf\"><code>ac5c5d3</code></a>\r\nPrepare release notes for v1.6.26</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/e7ca005043f6974536c3f8e0da42f93b5bdc2879\"><code>e7ca005</code></a>\r\nMerge pull request <a\r\nhref=\"https://redirect.github.com/containerd/containerd/issues/9472\">#9472</a>\r\nfrom thaJeztah/1.6_update_golang_1.20.12</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/7cbdfc92ef38f789f1a2773fa6fac405d361a6cc\"><code>7cbdfc9</code></a>\r\nupdate to go1.20.12, test go1.21.5</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/024b1cce6b27f10e00bb9bde33a5fe9563545f8d\"><code>024b1cc</code></a>\r\nupdate to go1.20.11, test go1.21.4</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/2e404598e7da93f4ad8b13bb6119441a5e3c83b0\"><code>2e40459</code></a>\r\nMerge pull request <a\r\nhref=\"https://redirect.github.com/containerd/containerd/issues/9484\">#9484</a>\r\nfrom ruiwen-zhao/cri-api-warning-1.6</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/64e56bfde95828660971673d20952f275cc2c0ba\"><code>64e56bf</code></a>\r\nAdd cri-api v1alpha2 usage warning to all api calls</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/c566b7d46668de23d2881eb86ce1b76ca23c8a75\"><code>c566b7d</code></a>\r\nMerge pull request <a\r\nhref=\"https://redirect.github.com/containerd/containerd/issues/9468\">#9468</a>\r\nfrom samuelkarp/deprecation-warning-runtime-1.6</li>\r\n<li><a\r\nhref=\"https://github.com/containerd/containerd/commit/efefd3bf334b5df0e97bff0be61ba906a9b3b528\"><code>efefd3b</code></a>\r\ntasks: emit warning for runc v1 runtime</li>\r\n<li>Additional commits viewable in <a\r\nhref=\"https://github.com/containerd/containerd/compare/v1.6.18...v1.6.26\">compare\r\nview</a></li>\r\n</ul>\r\n</details>\r\n<br />\r\n\r\n\r\n[![Dependabot compatibility\r\nscore](https://dependabot-badges.githubapp.com/badges/compatibility_score?dependency-name=github.com/containerd/containerd&package-manager=go_modules&previous-version=1.6.18&new-version=1.6.26)](https://docs.github.com/en/github/managing-security-vulnerabilities/about-dependabot-security-updates#about-compatibility-scores)\r\n\r\nYou can trigger a rebase of this PR by commenting `@dependabot rebase`.\r\n\r\n[//]: # (dependabot-automerge-start)\r\n[//]: # (dependabot-automerge-end)\r\n\r\n---\r\n\r\n<details>\r\n<summary>Dependabot commands and options</summary>\r\n<br />\r\n\r\nYou can trigger Dependabot actions by commenting on this PR:\r\n- `@dependabot rebase` will rebase this PR\r\n- `@dependabot recreate` will recreate this PR, overwriting any edits\r\nthat have been made to it\r\n- `@dependabot merge` will merge this PR after your CI passes on it\r\n- `@dependabot squash and merge` will squash and merge this PR after\r\nyour CI passes on it\r\n- `@dependabot cancel merge` will cancel a previously requested merge\r\nand block automerging\r\n- `@dependabot reopen` will reopen this PR if it is closed\r\n- `@dependabot close` will close this PR and stop Dependabot recreating\r\nit. You can achieve the same result by closing it manually\r\n- `@dependabot show <dependency name> ignore conditions` will show all\r\nof the ignore conditions of the specified dependency\r\n- `@dependabot ignore this major version` will close this PR and stop\r\nDependabot creating any more for this major version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this minor version` will close this PR and stop\r\nDependabot creating any more for this minor version (unless you reopen\r\nthe PR or upgrade to it yourself)\r\n- `@dependabot ignore this dependency` will close this PR and stop\r\nDependabot creating any more for this dependency (unless you reopen the\r\nPR or upgrade to it yourself)\r\nYou can disable automated security fix PRs for this repo from the\r\n[Security Alerts\r\npage](https://github.com/aws/shim-loggers-for-containerd/network/alerts).\r\n\r\n</details>\r\n\r\n> **Note**\r\n> Automatic rebases have been disabled on this pull request as it has\r\nbeen open for over 30 days.\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-03-21T08:35:38-07:00",
          "tree_id": "6fbea479e21f5df4ff2b8f4ea39099c741b37d7a",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/d4df5dd67316689b02ca909baa796a5ba7dac30d"
        },
        "date": 1711035517496,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7399885744,
            "unit": "ns/op\t41915688 B/op\t   66010 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7399885744,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 41915688,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66010,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b50076ae1eacfd4d57f7111cb180460cddd010cc",
          "message": "chore: add dependabot config (#67)\n\n*Issue #, if available:*\r\nN/A\r\n\r\n*Description of changes:*\r\nThis PR adds explicit configuration for dependabot for Go dependency and\r\nGitHub Actions dependency updates\r\n\r\nBy submitting this pull request, I confirm that you can use, modify,\r\ncopy, and redistribute this contribution, under the terms of your\r\nchoice.\r\n\r\nSigned-off-by: Austin Vazquez <macedonv@amazon.com>",
          "timestamp": "2024-03-22T22:30:12-07:00",
          "tree_id": "17b1dbaaef5575fbcb0d4de860e44ba7e547cc74",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/b50076ae1eacfd4d57f7111cb180460cddd010cc"
        },
        "date": 1711171986193,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7839188673,
            "unit": "ns/op\t42806664 B/op\t   66044 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7839188673,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42806664,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66044,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "617b48d51a60ef781fcbb3448dd263841a4779d1",
          "message": "chore(deps): bump actions/checkout from 3 to 4\n\nBumps [actions/checkout](https://github.com/actions/checkout) from 3 to 4.\n- [Release notes](https://github.com/actions/checkout/releases)\n- [Changelog](https://github.com/actions/checkout/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/actions/checkout/compare/v3...v4)\n\n---\nupdated-dependencies:\n- dependency-name: actions/checkout\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:38:40-07:00",
          "tree_id": "778fb30d0760fb0d371260672b951bf4764666b6",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/617b48d51a60ef781fcbb3448dd263841a4779d1"
        },
        "date": 1711172494423,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7847067793,
            "unit": "ns/op\t42816584 B/op\t   66101 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7847067793,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42816584,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66101,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "b2117287609243c3254d72c7ae0711f546a22e59",
          "message": "chore(deps): bump golangci/golangci-lint-action from 3 to 4\n\nBumps [golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action) from 3 to 4.\n- [Release notes](https://github.com/golangci/golangci-lint-action/releases)\n- [Commits](https://github.com/golangci/golangci-lint-action/compare/v3...v4)\n\n---\nupdated-dependencies:\n- dependency-name: golangci/golangci-lint-action\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:40:20-07:00",
          "tree_id": "57136a13e94f53e0cca574b28a4601e91875affc",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/b2117287609243c3254d72c7ae0711f546a22e59"
        },
        "date": 1711172591322,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8455131754,
            "unit": "ns/op\t40734656 B/op\t   65810 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8455131754,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 40734656,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65810,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "5bc6a14b00c95d5b8032e9d54ddc1cdc6df36018",
          "message": "chore(deps): bump github.com/google/uuid from 1.3.0 to 1.6.0\n\nBumps [github.com/google/uuid](https://github.com/google/uuid) from 1.3.0 to 1.6.0.\n- [Release notes](https://github.com/google/uuid/releases)\n- [Changelog](https://github.com/google/uuid/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/google/uuid/compare/v1.3.0...v1.6.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/google/uuid\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:41:10-07:00",
          "tree_id": "3128b1ba1638da2439cb1d2c728488b19338a5b4",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/5bc6a14b00c95d5b8032e9d54ddc1cdc6df36018"
        },
        "date": 1711172642508,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8349658417,
            "unit": "ns/op\t40519728 B/op\t   66009 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8349658417,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 40519728,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66009,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "b5d843d6d275e77c25e7ffcd34e03df89ee6e067",
          "message": "chore(deps): bump github.com/aws/smithy-go from 1.14.2 to 1.20.1\n\nBumps [github.com/aws/smithy-go](https://github.com/aws/smithy-go) from 1.14.2 to 1.20.1.\n- [Release notes](https://github.com/aws/smithy-go/releases)\n- [Changelog](https://github.com/aws/smithy-go/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/aws/smithy-go/compare/v1.14.2...v1.20.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/aws/smithy-go\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:42:15-07:00",
          "tree_id": "fc4a5bf3ad11b37c44cb1d7b8fab673742f80135",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/b5d843d6d275e77c25e7ffcd34e03df89ee6e067"
        },
        "date": 1711172706912,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8187971074,
            "unit": "ns/op\t41759120 B/op\t   65813 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8187971074,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 41759120,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65813,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "f1187b86483699457255ecafe8d651079d1752df",
          "message": "chore(deps): bump actions/setup-go from 4 to 5\n\nBumps [actions/setup-go](https://github.com/actions/setup-go) from 4 to 5.\n- [Release notes](https://github.com/actions/setup-go/releases)\n- [Commits](https://github.com/actions/setup-go/compare/v4...v5)\n\n---\nupdated-dependencies:\n- dependency-name: actions/setup-go\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:43:25-07:00",
          "tree_id": "d6861be9181ee5d5f0cdf1058632619031d20a82",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/f1187b86483699457255ecafe8d651079d1752df"
        },
        "date": 1711172781219,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7925396914,
            "unit": "ns/op\t42861656 B/op\t   66215 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7925396914,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42861656,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66215,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "8d6879748f8c6dd3c0baccdd200afaa815c6283e",
          "message": "chore(deps): bump github.com/onsi/gomega from 1.27.8 to 1.32.0\n\nBumps [github.com/onsi/gomega](https://github.com/onsi/gomega) from 1.27.8 to 1.32.0.\n- [Release notes](https://github.com/onsi/gomega/releases)\n- [Changelog](https://github.com/onsi/gomega/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/onsi/gomega/compare/v1.27.8...v1.32.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/onsi/gomega\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:47:11-07:00",
          "tree_id": "73e861aa5067bbb8e0abf162a901d0dc38b35e2f",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/8d6879748f8c6dd3c0baccdd200afaa815c6283e"
        },
        "date": 1711173004434,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8203830019,
            "unit": "ns/op\t40522368 B/op\t   65831 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8203830019,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 40522368,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65831,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "373cb3f7e5d293239f45317b783a30c64444c9e3",
          "message": "chore(deps): bump github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs\n\nBumps [github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs](https://github.com/aws/aws-sdk-go-v2) from 1.23.5 to 1.35.0.\n- [Release notes](https://github.com/aws/aws-sdk-go-v2/releases)\n- [Changelog](https://github.com/aws/aws-sdk-go-v2/blob/service/s3/v1.35.0/CHANGELOG.md)\n- [Commits](https://github.com/aws/aws-sdk-go-v2/compare/v1.23.5...service/s3/v1.35.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:48:09-07:00",
          "tree_id": "e4732d594279db86d3de54c19c8a2ffdfb1b9857",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/373cb3f7e5d293239f45317b783a30c64444c9e3"
        },
        "date": 1711173064056,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7872682902,
            "unit": "ns/op\t42868192 B/op\t   66103 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7872682902,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42868192,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66103,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "dea99cfffc7145ea49bab351d3f1e37654bd2a4d",
          "message": "chore(deps): bump github.com/aws/aws-sdk-go-v2/config\n\nBumps [github.com/aws/aws-sdk-go-v2/config](https://github.com/aws/aws-sdk-go-v2) from 1.18.39 to 1.27.9.\n- [Release notes](https://github.com/aws/aws-sdk-go-v2/releases)\n- [Commits](https://github.com/aws/aws-sdk-go-v2/compare/config/v1.18.39...config/v1.27.9)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/aws/aws-sdk-go-v2/config\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-22T22:53:28-07:00",
          "tree_id": "78f2539cb94fa772ff35104173dbcd7de8518a1e",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/dea99cfffc7145ea49bab351d3f1e37654bd2a4d"
        },
        "date": 1711173386476,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 8311114855,
            "unit": "ns/op\t38863176 B/op\t   66959 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 8311114855,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 38863176,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66959,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "d831847f9140b8fdf3b6517ce9f4493c3cedb5bc",
          "message": "chore(deps): bump github.com/stretchr/testify from 1.8.4 to 1.9.0\n\nBumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.8.4 to 1.9.0.\n- [Release notes](https://github.com/stretchr/testify/releases)\n- [Commits](https://github.com/stretchr/testify/compare/v1.8.4...v1.9.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/stretchr/testify\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-25T08:32:47-07:00",
          "tree_id": "92d7d31cbdb0ce5f635e190ee74887126a6eb029",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/d831847f9140b8fdf3b6517ce9f4493c3cedb5bc"
        },
        "date": 1711380937666,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7449077144,
            "unit": "ns/op\t42825608 B/op\t   66074 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7449077144,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 42825608,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 66074,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "41ac43af025ff6922e1f3d3ea22b018cba52cb9c",
          "message": "chore(deps): bump github.com/docker/go-units from 0.4.0 to 0.5.0\n\nBumps [github.com/docker/go-units](https://github.com/docker/go-units) from 0.4.0 to 0.5.0.\n- [Commits](https://github.com/docker/go-units/compare/v0.4.0...v0.5.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/docker/go-units\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-25T08:36:51-07:00",
          "tree_id": "231add9f27fd85c86a8585f8795f0e0eaeb2960b",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/41ac43af025ff6922e1f3d3ea22b018cba52cb9c"
        },
        "date": 1711381195540,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7925851151,
            "unit": "ns/op\t38970368 B/op\t   65932 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7925851151,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 38970368,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65932,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "33da11705e72e30ecd7f47cf601cd1f65864e2ae",
          "message": "chore(deps): bump golang.org/x/sync from 0.5.0 to 0.6.0\n\nBumps [golang.org/x/sync](https://github.com/golang/sync) from 0.5.0 to 0.6.0.\n- [Commits](https://github.com/golang/sync/compare/v0.5.0...v0.6.0)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/sync\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-25T08:38:20-07:00",
          "tree_id": "acb27b3f9365ddf4199d78c56e8850f7d97a02ff",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/33da11705e72e30ecd7f47cf601cd1f65864e2ae"
        },
        "date": 1711381282057,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7759401056,
            "unit": "ns/op\t39583384 B/op\t   65800 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7759401056,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 39583384,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65800,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
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
            "email": "55906459+austinvazquez@users.noreply.github.com",
            "name": "Austin Vazquez",
            "username": "austinvazquez"
          },
          "distinct": true,
          "id": "f594421443abecb3f1ea0431a6eb6f41a47fc89f",
          "message": "chore(deps): bump github.com/spf13/viper from 1.4.0 to 1.18.2\n\nBumps [github.com/spf13/viper](https://github.com/spf13/viper) from 1.4.0 to 1.18.2.\n- [Release notes](https://github.com/spf13/viper/releases)\n- [Commits](https://github.com/spf13/viper/compare/v1.4.0...v1.18.2)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/spf13/viper\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-03-25T08:43:43-07:00",
          "tree_id": "12428fc1d4b04cbeff5d7a5cd1fd73524cac9f46",
          "url": "https://github.com/aws/shim-loggers-for-containerd/commit/f594421443abecb3f1ea0431a6eb6f41a47fc89f"
        },
        "date": 1711381597828,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSplunk",
            "value": 7862179148,
            "unit": "ns/op\t39864296 B/op\t   65975 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - ns/op",
            "value": 7862179148,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - B/op",
            "value": 39864296,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkSplunk - allocs/op",
            "value": 65975,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
          }
        ]
      }
    ]
  }
}
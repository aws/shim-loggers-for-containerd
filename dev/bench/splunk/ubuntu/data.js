window.BENCHMARK_DATA = {
  "lastUpdate": 1698635184378,
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
      }
    ]
  }
}
# Amazon ECS shim loggers for containerd
Amazon ECS shim loggers for containerd is a collection of [containerd](https://github.com/containerd/containerd) compatible logger
implementations that send container logs to various destinations. The following destinations are currently supported:
* [Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html)
* [Splunk](https://www.splunk.com/en_us/central-log-management.html)
* [Fluentd](http://www.fluentd.org/)

## Build
Make sure you have [golang](https://golang.org) installed. Then simply run `make build` to build the respective binaries. You might need to execute `make get-deps` to install some of the dependencies.

## Usage
Containerd supports shim plugins that redirect container output to a custom binary on Linux using STDIO URIs with
[runc v2 runtime](https://github.com/containerd/containerd/tree/release/1.3/runtime/v2). These loggers can be used
either programmatically or with the [ctr](https://github.com/projectatomic/containerd/blob/master/docs/cli.md) tool.

### When using the `NewTask` API
When using the [`NewTask`](https://github.com/containerd/containerd/blob/release/1.3/container.go#L208) API
to start a container, simply provide the path to the built binary file `shim-loggers-for-containerd` and required
arguments. Note it's a good practice to clean up container resources with
[`Delete`](https://github.com/containerd/containerd/blob/release/1.3/task.go#L287) API call after container exited
as the container IO pipes are not closed if the shim process is still running.

Example:
```
NewTask(context, cio.BinaryIO("/usr/bin/shim-loggers-for-containerd", args))
```

### When using the `ctr` tool
When using [ctr](https://github.com/projectatomic/containerd/blob/master/docs/cli.md) tool to run
a container, provide the URI path to the binary file `shim-loggers-for-containerd` and required arguments as part of
the path.

Example:
```
ctr run \
    --runtime io.containerd.runc.v2 \
    --log-uri "binary:///usr/bin/shim-loggers-for-containerd?--log-driver=awslogs&--arg1=value1&-args2=value2" \
    docker.io/library/redis:alpine \
    redis
```

## Arguments

### Common arguments
The following list of arguments apply to all of the shim logger binaries in this repo:

|Name|Required|Description|
|-|-|-|
| log-driver | Yes | The name of the shim logger. Can be any of `awslogs`, `splunk` or `fluentd`. |
| container-id | Yes | The container id |
| container-name | Yes | The name of the container |
| mode | No | Either `blocking` or `non-blocking`. In the `non-blocking` mode, log events are buffered and the application continues to execute even if these logs can't be drained or sent to the destination. Logs could also be lost when the buffer is full. |
| max-buffer-size | No | Only supported in `non-blocking` mode. Set to `1m` (1MiB) by default. Example values: `200`, `4k`, `1m` etc. |
| uid | No | Set a custom uid for the shim logger process. `0` is not supported. |
| gid | No | Set a custom gid for the shim logger process. `0` is not supported. |
| cleanup-time | No | Set a custom time for the shim logger process clean up itself. Set to `5s` (5 seconds) by default. Note the maximum supported value is 12 seconds, since containerd shim sets shim logger cleanup timeout value as 12 seconds. See [reference](https://github.com/containerd/containerd/commit/0dc7c8595627e38ca2b83d17a062b51f384c2025). |
| container-image-id | No | The container image id. This is part of the docker config variables that can be logged by splunk log driver. |
| container-image-name | No | The container image name. This is part of the docker config variables that can be logged by splunk log driver. |
| container-env | No | The container environment variables map in json format. This is part of the docker config variables that can be logged by splunk log driver. |
| container-labels | No | The container labels map in json format. This is part of the docker config variables that can be logged by splunk log driver. |

### Windows specific arguments
The following list of arguments apply to Windows shim logger binaries in this repo:

|Name|Required|Description|
|-|-|-|
| log-file-dir | No | Only supported in Windows. Will be the path where shim logger's log files are written. By default it is `\ProgramData\Amazon\ECS\log\shim-logger`
| proxy-variable | No | Only supported in Windows. The proxy variable will set the `HTTP_PROXY` and `HTTPS_PROXY` environment variables.

### Additional log driver options
#### Amazon CloudWatch Logs
The following additional arguments are supported for the `awslogs` shim logger binary, which can be used to send container logs to [Amazon CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html).

|Name|Required|Description|
|-|-|-|
| awslogs-group | Yes | The [log group](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CloudWatchLogsConcepts.html) in which the log stream for the container will be created.|
| awslogs-stream | Yes | The [log stream name](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CloudWatchLogsConcepts.html) to stream container logs to. |
| awslogs-region | Yes | The region name in which the log group and log stream needs to be created in.|
| awslogs-credentials-endpoint | Yes | The endpoint from which credentials are retrieved from to connect to Amazon CloudWatch Logs.|
| awslogs-create-group | No | Set to `false` by default. If the provided log group name does not exist and this value is set to `false`, the binary will directly exit with an error|
| awslogs-create-stream | No | Set to `true` by default. The log stream will always be created unless this value specified to `false` explicitly.|
| awslogs-multiline-pattern | No | Matches the behavior of the [`awslogs` Docker log driver](https://docs.docker.com/config/containers/logging/awslogs/#amazon-cloudwatch-logs-options#awslogs-multiline-pattern).|
| awslogs-datetime-format | No | Matches the behavior of the [`awslogs` Docker log driver](https://docs.docker.com/config/containers/logging/awslogs/#amazon-cloudwatch-logs-options#awslogs-datetime-format)|

#### Splunk
The following additional arguments are supported for the `splunk` shim logger binary, which can be used to send container logs to [splunk](https://www.splunk.com/en_us/central-log-management.html).
You can find a description of what these parameters are used for [here](https://docs.docker.com/config/containers/logging/splunk/).

|Name|Required|
|-|-|
| splunk-token | Yes |
| splunk-url | Yes |
| splunk-source | No |
| splunk-sourcetype | No |
| splunk-index | No |
| splunk-capath | No |
| splunk-caname | No |
| splunk-insecureskipverify | No |
| splunk-format | No |
| splunk-verify-connection | No |
| splunk-gzip | No |
| splunk-gzip-level | No |
| splunk-tag | No |
| labels | No |
| env | No |
| env-regex | No |

#### Fluentd
The following additional arguments are supported for the `fluentd` shim logger binary, which can be used to send container logs to  [Fluentd](https://www.fluentd.org). Note that all of these are optional arguments.
* fluentd-address: The address of the Fluentd server to connect to. By default, the `localhost:24224` address is used.
* fluentd-async-connect: Specifies if the logger connects to Fluentd in background. Defaults to `false`.
* fluentd-sub-second-precision: Generates logs in nanoseconds. Defaults to `true`. Note that this is in contrast to the default behaviour of fluentd log driver where it defaults to `false`. 
* fluentd-tag: Specifies the tag used for log messages. Defaults to the first 12 characters of container ID.

## License

This project is licensed under the Apache-2.0 License.

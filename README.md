# ShimLoggersForContainerd
This repository is a collection of [containerd](https://github.com/containerd/containerd) compatible logger 
implementations that send container logs to different destinations.

For more information about log drivers, see [Docker logging drivers configuration](https://docs.docker.com/config/containers/logging/configure/).

## Build
Make sure you have [golang](https://golang.org) installed. Then simply run `make build` to build the binary.

## Usage
Containerd supports shim plugins that redirect container output to a custom binary on Linux using STDIO URIs with 
[runc v2 runtime](https://github.com/containerd/containerd/tree/release/1.3/runtime/v2). These loggers can be used 
either programmatically or with the [ctr](https://github.com/projectatomic/containerd/blob/master/docs/cli.md) tool.

* When using containerd [`NewTask`](https://github.com/containerd/containerd/blob/release/1.3/container.go#L208) API 
to start a container, simply provide the path to the built binary file `shim-loggers-for-containerd` and required 
arguments. Note it's a good practice to clean up container resources with 
[`Delete`](https://github.com/containerd/containerd/blob/release/1.3/task.go#L287) API call after container exited 
as the container IO pipes are not closed if the shim process is still running. 
    * Example: 
        `NewTask(context, cio.BinaryIO("/usr/bin/shim-loggers-for-containerd", args))`
* When using [ctr](https://github.com/projectatomic/containerd/blob/master/docs/cli.md) tool to run 
a container, provide the URI path to the binary file `shim-loggers-for-containerd` and required arguments as part of 
the path.
    * Example: 
        ```
        ctr run \ 
            --runtime io.containerd.runc.v2 \ 
            --log-uri "binary:///usr/bin/shim-loggers-for-containerd?--log-driver=awslogs&--arg1=value1&-args2=value2" \
            docker.io/library/redis:alpine \
            redis
        ```

## Arguments
* Required args:
    * log-driver
    * container-id
    * container-name
    * All other required arguments for chosen log driver
* Optional args:
    * mode
    * max-buffer-size
    * uid: set customized uid. Value of zero is not supported.
    * gid: set customized gid. Value of zero is not supported.
    * All other optional arguments for chosen log driver

## Supported log driver options
* `awslogs`: send container logs to [aws cloudwatch logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html). 
You can find more details [here](https://docs.docker.com/config/containers/logging/awslogs/).
    * Required arguments:
        * awslogs-group
        * awslogs-region
        * awslogs-stream
        * awslogs-credentials-endpoint
    * Optional arguments:
        * awslogs-create-group: default to be `false`. If the provided log group name does not exist and this value 
        is set to `false`, the binary will directly exit with errors.
        * awslogs-multiline-pattern: no default value
        * awslogs-datetime-format: no default value
        
* `fluentd`: send container logs to [Fluentd](https://www.fluentd.org).
You can find more details [here](https://docs.docker.com/config/containers/logging/fluentd/).
    * Required arguments: No required arguments
    * Optional arguments:
        * fluentd-address: default to connect to port `24224`
        * fluentd-async-connect: if connect fluentd in background. Default to be false.
        * fluentd-tag: tagging log message. Default to be first 12 characters of container ID.

* `splunk`: send container log to [splunk](www.splunk.com).
You can find more details [here](https://docs.docker.com/config/containers/logging/splunk/).
    * Required arguments:
        * splunk-token
        * splunk-url
    * Optional arguments:
        * splunk-source: Event source.
        * splunk-sourcetype: Event source type.
        * splunk-index: Event index.
        * splunk-capath: Path to root certificate.
        * splunk-caname: Name to use for validating server certificate; by default the hostname of the `splunk-url` is used.
        * splunk-insecureskipverify: Ignore server certificate validation.
        * splunk-format: Message format. Can be `inline`, `json` or `raw`. Defaults to `inline`.
        * splunk-verify-connection: Verify on start, that docker can connect to Splunk server. Defaults to true.
        * splunk-gzip: Enable/disable gzip compression to send events to Splunk Enterprise or Splunk Cloud instance. Defaults to false.
        * splunk-gzip-level: Set compression level for gzip. Valid values are `-1` (default), `0` (no compression), `1` (best speed) ... `9` (best compression). Defaults to DefaultCompression.
        * splunk-tag: Specify tag for message, which interpret some markup. Default value is `{{.ID}}` (12 characters of the container ID). 
        * labels: Comma-separated list of keys of labels, which should be included in message, if these labels are specified for container.
        * env: Comma-separated list of keys of environment variables, which should be included in message, if these variables are specified for container.
        * env-regex: Similar to and compatible with `env`. A regular expression to match logging-related environment variables. Used for advanced log tag options.

## Supported values for mode
* `blocking`: default mode
* `non-blocking`: saving containerd logs to an intermediate buffer consumed by log driver, which unblocks container 
performance if log driver has trouble sending logs to destination. Note in this mode, there may exist chance of losing 
container logs when buffer is full. More info can be found 
[here](https://docs.docker.com/config/containers/logging/configure/#configure-the-delivery-mode-of-log-messages-from-container-to-log-driver).

## Supported values for max-buffer-size
This value is only supported when `non-blocking` mode is enabled. Please provide it in a human readable format.
* Example: `200`, `4k`, `1m`
* Default to be 1 megabytes or `1m`.

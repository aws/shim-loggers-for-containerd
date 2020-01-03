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
arguments.
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

## Supported values for mode
* `blocking`: default mode

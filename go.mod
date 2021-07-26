module github.com/aws/shim-loggers-for-containerd

require (
	github.com/aws/aws-sdk-go v1.26.8 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/containerd/containerd v1.5.4
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/docker/docker v0.7.3-0.20190918143018-ad1b781e44fa
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0
	github.com/fluent/fluent-logger-golang v1.4.0 // indirect
	github.com/golang/mock v1.4.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1
	github.com/tinylib/msgp v1.1.1 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	gotest.tools v2.2.0+incompatible
)

// Workaround for optionally calling create log stream API
replace github.com/docker/docker v0.7.3-0.20190918143018-ad1b781e44fa => github.com/xia-wu/moby v17.12.0-ce-rc1.0.20210308205136-2cd0a2a46d81+incompatible

go 1.13

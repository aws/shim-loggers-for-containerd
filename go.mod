module github.com/aws/shim-loggers-for-containerd

require (
	github.com/aws/aws-sdk-go v1.26.8 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/containerd/containerd v1.5.13
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/docker/docker v20.10.13+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0
	github.com/fluent/fluent-logger-golang v1.9.0 // indirect
	github.com/golang/mock v1.4.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	github.com/tinylib/msgp v1.1.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gotest.tools v2.2.0+incompatible
)

go 1.13

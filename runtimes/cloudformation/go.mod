module github.com/falcosecurity/kilt/runtimes/cloudformation

go 1.18

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.34.27
	github.com/falcosecurity/kilt/pkg v0.0.0-20201012153322-cfbae90c1fbc
	github.com/google/go-containerregistry v0.4.0
	github.com/rs/zerolog v1.19.0
	github.com/stretchr/testify v1.8.1
	github.com/yudai/gojsondiff v1.0.0
)

require (
	github.com/containerd/stargz-snapshotter/estargz v0.0.0-20201223015020-a9a0c2d64694 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v0.0.0-20191017083524-a8ff7f821017 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7 // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/go-akka/configuration v0.0.0-20200606091224-a002c0330665 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/falcosecurity/kilt/pkg => ./../../pkg

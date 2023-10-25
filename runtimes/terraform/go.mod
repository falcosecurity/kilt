module github.com/falcosecurity/kilt/runtimes/terraform

go 1.18

require (
	github.com/falcosecurity/kilt/runtimes/cloudformation v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-framework v0.1.0
	github.com/hashicorp/terraform-plugin-go v0.3.1
)

require (
	github.com/Jeffail/gabs/v2 v2.7.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.14.3 // indirect
	github.com/docker/cli v24.0.2+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/falcosecurity/kilt/pkg v0.0.0-20230606123839-2e4c434d5d88 // indirect
	github.com/go-akka/configuration v0.0.0-20200606091224-a002c0330665 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-containerregistry v0.15.2 // indirect
	github.com/hashicorp/go-hclog v0.0.0-20180709165350-ff2cf002a8dd // indirect
	github.com/hashicorp/go-plugin v1.3.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/rs/zerolog v1.29.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/vbatts/tar-split v0.11.3 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.56.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/falcosecurity/kilt/pkg => ./../../pkg

replace github.com/falcosecurity/kilt/runtimes/cloudformation => ./../../runtimes/cloudformation

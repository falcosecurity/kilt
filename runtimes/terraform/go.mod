module github.com/falcosecurity/kilt/runtimes/terraform

go 1.18

require (
	github.com/falcosecurity/kilt/runtimes/cloudformation v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-framework v0.1.0
	github.com/hashicorp/terraform-plugin-go v0.3.1
)

require (
	github.com/Jeffail/gabs/v2 v2.6.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.0.0-20201223015020-a9a0c2d64694 // indirect
	github.com/docker/cli v0.0.0-20191017083524-a8ff7f821017 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7 // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/falcosecurity/kilt/pkg v0.0.0-20201012153322-cfbae90c1fbc // indirect
	github.com/go-akka/configuration v0.0.0-20200606091224-a002c0330665 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-containerregistry v0.4.0 // indirect
	github.com/hashicorp/go-hclog v0.0.0-20180709165350-ff2cf002a8dd // indirect
	github.com/hashicorp/go-plugin v1.3.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.19.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200527145253-8367513e4ece // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
)

replace github.com/falcosecurity/kilt/pkg => ./../../pkg

replace github.com/falcosecurity/kilt/runtimes/cloudformation => ./../../runtimes/cloudformation

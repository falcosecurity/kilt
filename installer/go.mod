module github.com/falcosecurity/kilt/installer

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.0
	github.com/aws/aws-sdk-go-v2/config v0.3.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v0.1.3
	github.com/aws/aws-sdk-go-v2/service/cloudformation v0.30.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sts v0.30.0
	github.com/markbates/pkger v0.17.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

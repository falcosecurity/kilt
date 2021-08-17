module github.com/falcosecurity/kilt/runtimes/terraform

go 1.16

require (
	github.com/falcosecurity/kilt/runtimes/cloudformation v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-framework v0.1.0
	github.com/hashicorp/terraform-plugin-go v0.3.1
)

replace github.com/falcosecurity/kilt/pkg => ./../../pkg

replace github.com/falcosecurity/kilt/runtimes/cloudformation => ./../../runtimes/cloudformation

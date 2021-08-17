package tf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	ImageAuthSecret types.String `tfsdk:"image_auth_secret"`
}

func (p *provider) GetSchema(ctx context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"image_auth_secret": {
				Type:     types.StringType,
				Required: false,
			},
		},
	}, nil
}

func (p *provider) Configure(ctx context.Context, request tfsdk.ConfigureProviderRequest, response *tfsdk.ConfigureProviderResponse) {
	err := request.Config.Get(ctx, p)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading configuration",
			Detail:   "An unexpected error was encountered while reading configuration: " + err.Error(),
		})
		return
	}
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{
		"kilt_fargate_container_definitions": dataSourceFargateContainerDefinitions{provider: p},
	}, nil
}

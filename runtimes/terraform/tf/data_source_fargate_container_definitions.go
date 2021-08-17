package tf

import (
	"context"
	"encoding/json"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type dataSourceFargateContainerDefinitions struct {
	provider *provider
}

type DataSourceConfig struct {
	ContainerDefinitions       types.String      `tfsdk:"container_definitions"`
	KiltDefinition             types.String      `tfsdk:"kilt_definition"`
	RecipeConfig               map[string]string `tfsdk:"recipe_config"`
	OutputContainerDefinitions types.String      `tfsdk:"output_container_definitions"`
}

func (d dataSourceFargateContainerDefinitions) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest, response *tfsdk.ReadDataSourceResponse) {
	var config DataSourceConfig
	err := request.Config.Get(ctx, &config)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading configuration",
			Detail:   "An unexpected error was encountered while reading configuration: " + err.Error(),
		})
		return
	}

	jsonConf, err := json.Marshal(&config.RecipeConfig)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to serialize recipe config",
			Detail:   "Failed to serialize recipe config: " + err.Error(),
		})
		return
	}

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               config.KiltDefinition.Value,
		ImageAuthSecret:    d.provider.ImageAuthSecret.Value,
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConf),
	}

	config.OutputContainerDefinitions, err = patchFargateTaskDefinition(ctx, &config.ContainerDefinitions, kiltConfig)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error applying kilt patch",
			Detail:   "Error applying kilt patch: " + err.Error(),
		})
		return
	}

	err = response.State.Set(ctx, &config)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting response",
			Detail:   "An unexpected error was encountered while setting the data source response: " + err.Error(),
		})
		return
	}
}

func (d dataSourceFargateContainerDefinitions) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"container_definitions": {
				Type:        types.StringType,
				Description: "the input Fargate container definitions to process with kilt",
				Required:    true,
			},
			"kilt_definition": {
				Type:        types.StringType,
				Description: "the kilt definition to apply to the input",
				Required:    true,
			},
			"recipe_config": {
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Description: "kilt recipe configuration",
				Required:    false,
			},
			"output_container_definitions": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (d dataSourceFargateContainerDefinitions) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, []*tfprotov6.Diagnostic) {
	return dataSourceFargateContainerDefinitions{provider: p.(*provider)}, nil
}

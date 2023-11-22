package tf

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type cfnProperties struct {
	RequiresCompatibilities []string                 `json:"RequiresCompatibilities"`
	ContainerDefinitions    []map[string]interface{} `json:"ContainerDefinitions"`
}

type cfnResource struct {
	ResourceType string        `json:"Type"`
	Properties   cfnProperties `json:"Properties"`
}

type cfnStack struct {
	Resources map[string]cfnResource `json:"Resources"`
}

func patchFargateTaskDefinition(ctx context.Context, containerDefinitions *types.String, kiltConfig *cfnpatcher.Configuration) (patched types.String, err error) {
	var cdefs []map[string]interface{}
	err = json.Unmarshal([]byte(containerDefinitions.Value), &cdefs)
	if err != nil {
		return types.String{Unknown: true}, err
	}

	stack := cfnStack{
		Resources: map[string]cfnResource{
			"kilt": {
				ResourceType: "AWS::ECS::TaskDefinition",
				Properties: cfnProperties{
					RequiresCompatibilities: []string{"FARGATE"},
					ContainerDefinitions:    cdefs,
				},
			},
		},
	}

	patchedStack, err := json.Marshal(stack)
	if err != nil {
		return types.String{Unknown: true}, err
	}

	defer func() {
		if r := recover(); r != nil {
			patched = types.String{Unknown: true}
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

	templateParameters := make([]byte, 0)
	patchedBytes, err := cfnpatcher.Patch(ctx, kiltConfig, patchedStack, templateParameters)
	if err != nil {
		return types.String{Unknown: true}, err
	}

	err = json.Unmarshal(patchedBytes, &stack)
	if err != nil {
		return types.String{Unknown: true}, err
	}

	patchedBytes, err = json.Marshal(stack.Resources["kilt"].Properties.ContainerDefinitions)
	if err != nil {
		return types.String{Unknown: true}, err
	}

	return types.String{Value: string(patchedBytes)}, nil
}

package hocon

import (
	"encoding/json"
	"fmt"
	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/go-akka/configuration"
)

var defaults = `
build {
	entry_point: ${original.entry_point}
	command: ${original.command}
	image: ${original.image}
	environment_variables: ${original.environment_variables}

	mount: []
}
`

type KiltHocon struct {
	definition string
	config 	   string
}

type HoconProvided struct {
	Image string `json:"image"`
}

func NewKiltHocon(definition string) *KiltHocon {
	return NewKiltHoconWithConfig(definition, "{}")
}

func NewKiltHoconWithConfig(definition string, recipeConfig string) *KiltHocon {
	h := new(KiltHocon)
	h.definition = definition
	h.config = recipeConfig
	return h
}

func (k *KiltHocon) prepareFullStringConfig(info *kilt.TargetInfo) (*configuration.Config, error) {
	rawVars, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("could not serialize info: %w", err)
	}

	configString := "original:" + string(rawVars) + "\n" +
		"config:" + k.config + "\n" +
		defaults + k.definition

	return configuration.ParseString(configString), nil
}

func (k *KiltHocon) Build(info *kilt.TargetInfo) (*kilt.Build, error) {
	config, err := k.prepareFullStringConfig(info)
	if err != nil {
		return nil, fmt.Errorf("could not assemble full config: %w", err)
	}

	return extractBuild(config)
}

func (k *KiltHocon) Runtime(info *kilt.TargetInfo) (*kilt.Runtime, error) {
	config, err := k.prepareFullStringConfig(info)
	if err != nil {
		return nil, fmt.Errorf("could not assemble full config: %w", err)
	}
	if ! config.HasPath("runtime"){
		return nil, fmt.Errorf("definition does not have a runtime section")
	}
	return extractRuntime(config)
}


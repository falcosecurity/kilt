package cfnpatcher

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/google/go-containerregistry/pkg/crane"
)

type PartialImageConfig struct {
	Entrypoint []string
	Command []string
}

func GetConfigFromRepository(image string) (*PartialImageConfig,error) {
	ic := new(PartialImageConfig)

	res, err := crane.Config(image)
	if err != nil {
		return nil, fmt.Errorf("could not get defaults about image %s: %w", image, err)
	}
	cont, err := gabs.ParseJSON(res)
	if err != nil {
		return nil, fmt.Errorf("could not parse response from registry for image %s: %w", image, err)
	}

	if cont.Exists("config", "Entrypoint") {
		for _, v := range cont.S("config", "Entrypoint").Children() {
			if a, ok := v.Data().(string); ok {
				ic.Entrypoint = append(ic.Entrypoint, a)
			}
		}
	}

	if cont.Exists("config", "Cmd") {
		for _, v := range cont.S("config", "Cmd").Children() {
			if a, ok := v.Data().(string); ok {
				ic.Command = append(ic.Command, a)
			}
		}
	}

	return ic, nil
}
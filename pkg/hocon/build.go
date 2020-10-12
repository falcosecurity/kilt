package hocon

import (
	"fmt"
	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/go-akka/configuration"
)

func extractBuild(config *configuration.Config) (*kilt.Build, error) {
	b := new(kilt.Build)

	b.Image = config.GetString("build.image")
	b.EntryPoint = config.GetStringList("build.entry_point")
	b.Command = config.GetStringList("build.command")

	b.EnvironmentVariables = extractToStringMap(config, "build.environment_variables")

	if config.IsArray("build.mount") {
		mounts := config.GetValue("build.mount").GetArray()

		for k, m := range mounts {
			if m.IsObject() {
				mount := m.GetObject()

				resource := kilt.BuildResource{
					Name:       mount.GetKey("name").GetString(),
					Image:      mount.GetKey("image").GetString(),
					Volumes:    mount.GetKey("volumes").GetStringList(),
					EntryPoint: mount.GetKey("entry_point").GetStringList(),
				}

				if resource.Image == "" || len(resource.Volumes) == 0 || len(resource.EntryPoint) == 0 {
					return nil, fmt.Errorf("error at build.mount.%d: image, volumes and entry_point are all required ", k)
				}

				b.Resources = append(b.Resources, resource)
			}
		}
	}

	return b, nil
}

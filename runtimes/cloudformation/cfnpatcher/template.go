package cfnpatcher

import (
	"context"
	"os"

	"github.com/Jeffail/gabs/v2"
	"github.com/rs/zerolog/log"

	"github.com/falcosecurity/kilt/pkg/kilt"
)

type TemplateInfo struct {
	TargetInfo           *kilt.TargetInfo
	// Containers are not null when template values are complex
	Name                 *gabs.Container
	Image                *gabs.Container
	EntryPoint           []*gabs.Container
	Command              []*gabs.Container
	EnvironmentVariables map[string]*gabs.Container
}

func GetValueFromTemplate(what *gabs.Container) (string, *gabs.Container) {
	var result string
	var fallback *gabs.Container

	switch v := what.Data().(type) {
	case string:
		result = v
		fallback = nil
	default:
		fallback = what
		result = "placeholder: " + what.String()
	}
	return result, fallback
}

func extractContainerInfo(ctx context.Context, group *gabs.Container, groupName string, container *gabs.Container, configuration *Configuration) *TemplateInfo {
	cfnInfo := new(TemplateInfo)
	info := new(kilt.TargetInfo)
	cfnInfo.TargetInfo = info
	l := log.Ctx(ctx)

	info.ContainerName, cfnInfo.Image = GetValueFromTemplate(container.S("Name"))
	info.ContainerGroupName = groupName
	info.EnvironmentVariables = make(map[string]string)
	cfnInfo.EnvironmentVariables = make(map[string]*gabs.Container)
	info.Metadata = make(map[string]string)

	if container.Exists("Image") {
		info.Image, cfnInfo.Image = GetValueFromTemplate(container.S("Image"))

		os.Setenv("HOME", "/tmp") // crane requires $HOME variable
		repoInfo, err := GetConfigFromRepository(info.Image)
		if err != nil {
			l.Warn().Str("image", info.Image).Err(err).Msg("could not retrieve metadata from repository")
		} else {
			if configuration.UseRepositoryHints {
				l.Info().Str("image", info.Image).Msgf("extracted info from remote repository: %+v", repoInfo)
				if repoInfo.Entrypoint != nil {
					info.EntryPoint = repoInfo.Entrypoint
					cfnInfo.EntryPoint = make([]*gabs.Container, len(info.EntryPoint))
				}
				if repoInfo.Command != nil {
					info.Command = repoInfo.Command
					cfnInfo.Command = make([]*gabs.Container, len(info.Command))
				}
			}
		}
	}

	if container.Exists("EntryPoint") {
		info.EntryPoint = make([]string, 0)
		cfnInfo.EntryPoint = make([]*gabs.Container,0)
		for _, arg := range container.S("EntryPoint").Children() {
			passthrough, templateVal := GetValueFromTemplate(arg)
			cfnInfo.EntryPoint = append(cfnInfo.EntryPoint, templateVal)
			info.EntryPoint = append(info.EntryPoint, passthrough)
		}
	} else {
		l.Warn().Str("image", info.Image).Msg("no EntryPoint was specified")
	}

	if container.Exists("Command") {
		info.Command = make([]string, 0)
		cfnInfo.Command = make([]*gabs.Container,0)
		for _, arg := range container.S("Command").Children() {
			passthrough, templateVal := GetValueFromTemplate(arg)
			cfnInfo.Command = append(cfnInfo.Command, templateVal)
			info.Command = append(info.Command, passthrough)
		}
	} else {
		l.Warn().Str("image", info.Image).Msg("no Command was specified")
	}

	if container.Exists("Environment") {
		for _, env := range container.S("Environment").Children() {
			k, ok := env.S("Name").Data().(string)
			if ! ok {
				l.Fatal().Str("Fragment", env.S("Name").String()).Str("TaskDefinition", groupName).Msg("Environment has an unsupported value type. Expected string")
			}
			passthrough, templateVal := GetValueFromTemplate(env.S("Value"))

			cfnInfo.EnvironmentVariables[k] = templateVal
			info.EnvironmentVariables[k] = passthrough
		}
	}

	if group.Exists("Properties", "Tags") {
		for _, tag := range group.S("Properties", "Tags").Children() {
			if tag.Exists("Key") && tag.Exists("Value") {
				k, ok := tag.S("Key").Data().(string)
				if !ok {
					l.Fatal().Str("Fragment", tag.String()).Str("TaskDefinition", groupName).Msg("Tags has an unsupported key type")
				}

				passthrough, _ := GetValueFromTemplate(tag.S("Value"))
				info.Metadata[k] = passthrough
			}
		}
	}

	// TODO(admiral0): metadata tags

	return cfnInfo
}

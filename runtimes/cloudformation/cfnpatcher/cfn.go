package cfnpatcher

import (
	"context"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Kilt               string
	ImageAuthSecret    string
	OptIn              bool
	RecipeConfig       string
	UseRepositoryHints bool
	LogGroup           string
	ParameterizeEnvars bool
}

type InstrumentationHints struct {
	IgnoreContainersNamed  []string
	IncludeContainersNamed []string
	HasGlobalInclude       bool
}

const KiltIgnoreTag = "kilt-ignore"
const KiltIncludeTag = "kilt-include"
const KiltIgnoreContainersTag = "kilt-ignore-containers"
const KiltIncludeContainersTag = "kilt-include-containers"

var OptTagKeys = []string {KiltIgnoreTag, KiltIncludeTag, KiltIgnoreContainersTag, KiltIncludeContainersTag}

func isOptTagKey(key string) bool {
	for _, v := range OptTagKeys {
		if key == v {
			return true
		}
	}
	return false
}

func isIgnored(tags map[string]string, isOptIn bool) bool {
	_, included := tags[KiltIncludeTag]
	_, ignored := tags[KiltIgnoreTag]
	_, hasNamedContainerIncluded := tags[KiltIncludeContainersTag]

	return !((isOptIn && (included || hasNamedContainerIncluded)) || (!isOptIn && !ignored))
}

func extractContainersFromTag(tags map[string]string, tag string) []string {
	containers := make([]string, 0)
	containerList, hasIgnores := tags[tag]
	if hasIgnores {
		containers = strings.Split(containerList, ":")
	}
	return containers
}

func extractHintsFromTags(tags map[string]string) *InstrumentationHints {
	_, included := tags[KiltIncludeTag]
	return &InstrumentationHints{
		IgnoreContainersNamed:  extractContainersFromTag(tags, KiltIgnoreContainersTag),
		IncludeContainersNamed: extractContainersFromTag(tags, KiltIncludeContainersTag),
		HasGlobalInclude:       included,
	}
}

func Patch(ctx context.Context, configuration *Configuration, fragment, templateParameters []byte) ([]byte, error) {
	l := log.Ctx(ctx)
	template, err := gabs.ParseJSON(fragment)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse input fragment")
		return nil, err
	}

	if configuration.ParameterizeEnvars {
		l.Info().Msg("parameterizing recipe envars")
		applyParametersPatch(ctx, template, configuration)
	}

	var parameters *gabs.Container
	if len(templateParameters) > 0 {
		parameters, err = gabs.ParseJSON(templateParameters)
		if err != nil {
			l.Error().Err(err).Msg("failed to parse input templateParameters")
			return nil, err
		}
	}

	for name, resource := range template.S("Resources").ChildrenMap() {
		if matchFargate(resource) {

			optTags := getOptTags(resource)

			if isIgnored(optTags, configuration.OptIn) {
				l.Info().Str("resource", name).Msg("ignored resource due to tag")
				continue
			}
			l.Info().Str("resource", name).Msg("patching task definition")
			if err != nil {
				l.Error().Err(err).Str("resource", name).Msg("could not generate kilt instructions")
				continue
			}
			hints := extractHintsFromTags(optTags)
			_, err = applyTaskDefinitionPatch(ctx, name, resource, parameters, configuration, hints)
			if err != nil {
				l.Error().Err(err).Str("resource", name).Msgf("could not patch resource")
			}
		}
	}

	return template.Bytes(), nil
}

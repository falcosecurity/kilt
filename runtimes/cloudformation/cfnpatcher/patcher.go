package cfnpatcher

import (
	"context"
	"fmt"

	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/falcosecurity/kilt/pkg/kiltapi"

	"github.com/Jeffail/gabs/v2"
	"github.com/rs/zerolog/log"
)

func containerInConfig(name string, listOfNames []string) bool {
	for _, n := range listOfNames {
		if n == name {
			return true
		}
	}
	return false
}

func shouldSkip(info *kilt.TargetInfo, configuration *Configuration, hints *InstrumentationHints) bool {
	isForceIncluded := containerInConfig(info.ContainerName, hints.IncludeContainersNamed)
	isExcluded := containerInConfig(info.ContainerName, hints.IgnoreContainersNamed)

	return (configuration.OptIn && !isForceIncluded && !hints.HasGlobalInclude) || (!configuration.OptIn && isExcluded)
}

func applyTaskDefinitionPatch(ctx context.Context, name string, resource *gabs.Container, configuration *Configuration, hints *InstrumentationHints) (*gabs.Container, error) {
	l := log.Ctx(ctx)

	successes := 0
	containers := make(map[string]kilt.BuildResource)
	k := kiltapi.NewKiltFromHoconWithConfig(configuration.Kilt, configuration.RecipeConfig)
	if resource.Exists("Properties", "ContainerDefinitions") {
		for _, container := range resource.S("Properties", "ContainerDefinitions").Children() {
			info := extractContainerInfo(resource, name, container)
			l.Info().Msgf("extracted info for container: %v", info)
			if shouldSkip(info, configuration, hints) {
				l.Info().Msgf("skipping container due to hints in tags")
				continue
			}
			patch, err := k.Build(info)
			if err != nil {
				return nil, fmt.Errorf("could not construct kilt patch: %w", err)
			}
			l.Info().Msgf("created patch for container: %v", patch)
			err = applyContainerDefinitionPatch(l.WithContext(ctx), container, patch)
			if err != nil {
				l.Warn().Str("resource", name).Err(err).Msg("skipped patching container in task definition")
			} else {
				successes += 1
			}

			for _, appendResource := range patch.Resources {
				containers[appendResource.Name] = appendResource
			}
		}
		err := appendContainers(resource, containers, configuration.ImageAuthSecret)
		if err != nil {
			return nil, fmt.Errorf("could not append container: %w", err)
		}
	}
	if successes == 0 {
		return resource, fmt.Errorf("could not patch a single container in the task")
	}
	return resource, nil
}

func applyContainerDefinitionPatch(ctx context.Context, container *gabs.Container, patch *kilt.Build) error {
	l := log.Ctx(ctx)

	_, err := container.Set(patch.EntryPoint, "EntryPoint")
	if err != nil {
		return fmt.Errorf("could not set EntryPoint: %w", err)
	}
	_, err = container.Set(patch.Command, "Command")
	if err != nil {
		return fmt.Errorf("could not set Command: %w", err)
	}
	_, err = container.Set(patch.Image, "Image")
	if err != nil {
		return fmt.Errorf("could not set Command: %w", err)
	}

	if !container.Exists("VolumesFrom") {
		_, err = container.Set([]interface{}{}, "VolumesFrom")
		if err != nil {
			return fmt.Errorf("could not set VolumesFrom: %w", err)
		}
	}

	for _, newContainer := range patch.Resources {
		// Skip containers with no volumes - just injecting sidecars
		if len(newContainer.Volumes) == 0 {
			l.Info().Msgf("Skipping injection of %s because it has no volumes specified", newContainer.Name)
			continue
		}
		addVolume := map[string]interface{}{
			"ReadOnly":        true,
			"SourceContainer": newContainer.Name,
		}

		_, err = container.Set(addVolume, "VolumesFrom", "-")
		if err != nil {
			return fmt.Errorf("could not add VolumesFrom directive: %w", err)
		}
	}

	if !(container.Exists("Environment") || len(patch.EnvironmentVariables) == 0) {
		_, err = container.Set([]interface{}{}, "Environment")

		if err != nil {
			return fmt.Errorf("could not add environment variable container: %w", err)
		}
	}
	for k, v := range patch.EnvironmentVariables {
		keyValue := make(map[string]string)
		keyValue["Name"] = k
		keyValue["Value"] = v
		_, err = container.Set(keyValue, "Environment", "-")

		if err != nil {
			return fmt.Errorf("could not add environment variable %v: %w", keyValue, err)
		}

	}

	// We need to add SYS_PTRACE capability to the container
	if !container.Exists("LinuxParameters") {
		emptyMap := make(map[string]interface{})
		_, err = container.Set(emptyMap, "LinuxParameters")
		if err != nil {
			return fmt.Errorf("could not add LinuxParameters: %w", err)
		}
	}

	if !container.Exists("LinuxParameters", "Capabilities") {
		emptyMap := make(map[string]interface{})
		_, err = container.Set(emptyMap, "LinuxParameters", "Capabilities")
		if err != nil {
			return fmt.Errorf("could not add LinuxParameters.Capabilities: %w", err)
		}
	}

	// fargate only supports SYS_PTRACE
	_, err = container.Set([]string{"SYS_PTRACE"}, "LinuxParameters", "Capabilities", "Add")
	if err != nil {
		return fmt.Errorf("could not add LinuxParamaters.Capabilities.Add: %w", err)
	}

	return nil
}

func appendContainers(resource *gabs.Container, containers map[string]kilt.BuildResource, imageAuth string) error {
	for _, inject := range containers {
		appended := map[string]interface{}{
			"Name":       inject.Name,
			"Image":      inject.Image,
			"EntryPoint": inject.EntryPoint,
		}
		if len(imageAuth) > 0 {
			appended["RepositoryCredentials"] = map[string]interface{}{
				"CredentialsParameter": imageAuth,
			}
		}
		_, err := resource.Set(appended, "Properties", "ContainerDefinitions", "-")
		if err != nil {
			return fmt.Errorf("could not inject %s: %w", inject.Name, err)
		}
	}
	return nil
}

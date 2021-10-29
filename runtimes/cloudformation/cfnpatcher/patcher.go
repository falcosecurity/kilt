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
			info := extractContainerInfo(ctx, resource, name, container, configuration)
			l.Info().Msgf("extracted info for container: %+v %+v", info.TargetInfo, info)
			if shouldSkip(info.TargetInfo, configuration, hints) {
				l.Info().Msgf("skipping container due to hints in tags")
				continue
			}
			patch, err := k.Build(info.TargetInfo)
			if err != nil {
				return nil, fmt.Errorf("could not construct kilt patch: %w", err)
			}
			l.Info().Msgf("created patch for container: %v", patch)
			err = applyContainerDefinitionPatch(l.WithContext(ctx), container, patch, info)
			if err != nil {
				l.Warn().Str("resource", name).Err(err).Msg("skipped patching container in task definition")
			} else {
				successes += 1
			}

			for _, appendResource := range patch.Resources {
				containers[appendResource.Name] = appendResource
			}
		}
		err := appendContainers(resource, containers, configuration.ImageAuthSecret, configuration.LogGroup, name)
		if err != nil {
			return nil, fmt.Errorf("could not append container: %w", err)
		}
	}
	if successes == 0 {
		return resource, fmt.Errorf("could not patch a single container in the task")
	}
	return resource, nil
}

func postPatchReplace(patched []string, original []string, parallel []*gabs.Container) []*gabs.Container {
	value := make([]*gabs.Container, 0)
	originalLen := len(original)
	for i, j := 0, 0; i < len(patched); i++ {
		toAssign := gabs.New()
		toAssign, _ = toAssign.Set(patched[i])
		if j < originalLen && patched[i] == original[j] {
			if parallel[j] != nil {
				toAssign = parallel[j]
			}
			j++
		}
		value = append(value, toAssign)
	}
	return value
}

func postPatchSelect(patched string, previous string, original *gabs.Container)  interface{} {
	if patched == previous && original != nil {
		return original
	}
	return patched
}

func applyContainerDefinitionPatch(ctx context.Context, container *gabs.Container, patch *kilt.Build, cfnInfo *TemplateInfo) error {
	l := log.Ctx(ctx)

	finalEntryPoint := postPatchReplace(patch.EntryPoint, cfnInfo.TargetInfo.EntryPoint, cfnInfo.EntryPoint)
	finalCommand := postPatchReplace(patch.Command, cfnInfo.TargetInfo.Command, cfnInfo.Command)

	_, err := container.Set(finalEntryPoint, "EntryPoint")
	if err != nil {
		return fmt.Errorf("could not set EntryPoint: %w", err)
	}
	_, err = container.Set(finalCommand, "Command")
	if err != nil {
		return fmt.Errorf("could not set Command: %w", err)
	}


	_, err = container.Set(postPatchSelect(patch.Image, cfnInfo.TargetInfo.Image, cfnInfo.Image), "Image")
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

	if len(patch.EnvironmentVariables) > 0 {
		_, err = container.Set([]interface{}{}, "Environment")

		if err != nil {
		return fmt.Errorf("could not add environment variable container: %w", err)
	}
	}

	for k, v := range patch.EnvironmentVariables {
		keyValue := make(map[string]interface{})
		keyValue["Name"] = k
		keyValue["Value"] = v
		if v == cfnInfo.TargetInfo.EnvironmentVariables[k] && cfnInfo.EnvironmentVariables[k] != nil {
			keyValue["Value"] = cfnInfo.EnvironmentVariables[k]
		}

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

func appendContainers(resource *gabs.Container, containers map[string]kilt.BuildResource, imageAuth string, logGroup string, name string) error {
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
		if len(logGroup) > 0 {
			appended["LogConfiguration"] = prepareLogConfiguration(name, logGroup)
		}
		_, err := resource.Set(appended, "Properties", "ContainerDefinitions", "-")
		if err != nil {
			return fmt.Errorf("could not inject %s: %w", inject.Name, err)
		}
	}
	return nil
}

func prepareLogConfiguration(taskName string, logGroup string) map[string]interface{} {
	// assuming that all given log configurations are for the awslogs driver
	config := map[string]interface{}{
		"LogDriver": "awslogs",
		"Options": map[string]interface{}{
			"awslogs-region": map[string]interface{}{
				"Ref": "AWS::Region",
			},
			"awslogs-group":         logGroup,
			"awslogs-stream-prefix": taskName,
		},
	}

	return config
}

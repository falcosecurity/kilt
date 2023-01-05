package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/config"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MacroInput struct {
	Region      string          `json:"region"`
	AccountID   string          `json:"accountId"`
	RequestID   string          `json:"requestId"`
	TransformID string          `json:"transformId"`
	Fragment    json.RawMessage `json:"fragment"`
}

type MacroOutput struct {
	RequestID string          `json:"requestId"`
	Status    string          `json:"status"`
	Fragment  json.RawMessage `json:"fragment"`
}

func HandleRequest(configuration *cfnpatcher.Configuration, ctx context.Context, event MacroInput) (MacroOutput, error) {
	l := log.With().
		Str("region", event.Region).
		Str("account", event.AccountID).
		Str("requestId", event.RequestID).
		Str("transformId", event.TransformID).
		Logger()
	loggerCtx := l.WithContext(ctx)
	result, err := cfnpatcher.Patch(loggerCtx, configuration, event.Fragment)
	if err != nil {
		return MacroOutput{event.RequestID, "failure", result}, err
	}
	log.Info().Str("template", string(result)).Msg("processing complete")
	return MacroOutput{event.RequestID, "success", result}, nil
}

func PatchLocalFile(configuration *cfnpatcher.Configuration, ctx context.Context, inputFile string) ([]byte, error) {
	l := log.With().
		Str("region", "local").
		Logger()
	loggerCtx := l.WithContext(ctx)

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		l.Error().Err(err).Msgf("cannot read file %s", inputFile)
		return nil, err
	}

	result, err := cfnpatcher.Patch(loggerCtx, configuration, inputData)
	if err != nil {
		l.Error().Err(err).Msg("failed to patch local file")
		return nil, err
	}

	log.Info().Str("template", string(result)).Msg("processing complete")
	return result, nil
}

func GetConfig() *cfnpatcher.Configuration {
	definition := os.Getenv("KILT_DEFINITION")
	definitionType := os.Getenv("KILT_DEFINITION_TYPE")
	optIn := os.Getenv("KILT_OPT_IN")
	imageAuth := os.Getenv("KILT_IMAGE_AUTH_SECRET")
	recipeConfig := os.Getenv("KILT_RECIPE_CONFIG")
	disableRepoHints := os.Getenv("KILT_DISABLE_REPO_HINTS")
	logGroup := os.Getenv("KILT_LOG_GROUP")
	var fullDefinition string
	switch definitionType {
	case config.S3:
		fullDefinition = config.FromS3(definition, false)
	case config.S3Gz:
		fullDefinition = config.FromS3(definition, true)
	case config.Http:
		fullDefinition = config.FromWeb(definition)
	case config.Base64:
		fullDefinition = config.FromBase64(definition, false)
	case config.Base64Gz:
		fullDefinition = config.FromBase64(definition, true)
	default:
		panic("unrecognized definition type - " + definitionType)
	}

	configuration := &cfnpatcher.Configuration{
		Kilt:               fullDefinition,
		ImageAuthSecret:    imageAuth,
		OptIn:              optIn != "",
		RecipeConfig:       recipeConfig,
		UseRepositoryHints: disableRepoHints == "",
		LogGroup:           logGroup,
	}

	return configuration
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configuration := GetConfig()
	switch os.Getenv("KILT_MODE") {
	case "local":
		result, err := PatchLocalFile(configuration, context.Background(), os.Getenv("KILT_SRC_TEMPLATE"))
		if err != nil {
			panic("cannot patch local file " + os.Getenv("KILT_SRC_TEMPLATE"))
		}

		err = os.WriteFile(os.Getenv("KILT_OUT_TEMPLATE"), result, 0644)
		if err != nil {
			panic("cannot write dst file " + os.Getenv("KILT_OUT_TEMPLATE"))
		}

	default:
		lambda.Start(
			func(ctx context.Context, event MacroInput) (MacroOutput, error) {
				return HandleRequest(configuration, ctx, event)
			})
	}
}

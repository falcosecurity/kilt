package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"

	cfnmacro "github.com/falcosecurity/kilt/installer/runtimes/cfn-macro"
	"github.com/falcosecurity/kilt/installer/util"
)

func initializeAwsBucket(cfg aws.Config) (string, error) {
	awsAccount, err := util.GetAwsAccountName(cfg, nil)
	if err != nil {
		return "", cli.Exit("could not identify aws account ID: %s"+err.Error(), 1)

	}
	awsRegion, err := util.GetRegion(cfg)
	if err != nil {
		return "", cli.Exit("could not autodetect region. please use -r parameter\nerror:"+err.Error(), 1)
	}
	bucket, err := util.GetBucketName(awsAccount, awsRegion)
	if err != nil {
		return "", cli.Exit("could not compute bucket name: "+err.Error(), 1)
	}
	err = util.EnsureBucketExists(bucket, s3.NewFromConfig(cfg))
	if err != nil {
		return "", cli.Exit("could not ensure existence of kilt s3 bucket: " + err.Error(), 1)
	}
	return bucket, nil
}

func registerCfnMacro() []*cli.Command {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("could not load AWS config:\n")
		panic(err)
	}

	return []*cli.Command{
		{
			Name:     "cfn-macro",
			Usage:    "Uses cloud formation macros to alter ECS Task Definitions",
			Category: "runtimes",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "Lists installed CFN macros",
					Action: func(c *cli.Context) error {
						cfg.Region = c.String("region")
						bucket, err := initializeAwsBucket(cfg)
						if err != nil {
							return err
						}
						cfnMacro := cfnmacro.NewCfnMacroRuntime(cfg, bucket)
						err = cfnMacro.List()
						if err != nil {
							return cli.Exit("could not list macros: "+err.Error(), 2)
						}
						return nil
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "region",
							Aliases:     []string{"r"},
							EnvVars:     []string{"AWS_REGION"},
							DefaultText: cfg.Region,
							Value:       cfg.Region,
						},
					},
				},
				{
					Name:  "delete",
					Usage: "Deletes CFN macros",
					Action: func(c *cli.Context) error {
						cfg.Region = c.String("region")
						bucket, err := initializeAwsBucket(cfg)
						if err != nil {
							return err
						}
						cfnMacro := cfnmacro.NewCfnMacroRuntime(cfg, bucket)
						err = cfnMacro.Delete(c.Args().Slice())
						if err != nil {
							return cli.Exit(err.Error(), 2)
						}
						return nil
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "region",
							Aliases:     []string{"r"},
							EnvVars:     []string{"AWS_REGION"},
							DefaultText: cfg.Region,
							Value:       cfg.Region,
						},
					},
				},
				{
					Name:        "install",
					Usage:       "Installs a new kilt CFN macro",
					Description: "Creates a CFN macro named MACRO_NAME that applies KILT_DEFINITION to CFN templates annotated with Transform: MACRO_NAME",
					ArgsUsage:   "KILT_DEFINITION MACRO_NAME",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "opt-in",
							Aliases: []string{"o"},
							Usage:   "Use opt-in logic instead of default opt-out",
						},
						&cli.StringFlag{
							Name:        "region",
							Aliases:     []string{"r"},
							Usage:       "Specify region for installation",
							EnvVars:     []string{"AWS_REGION"},
							DefaultText: cfg.Region,
						},
						&cli.PathFlag{
							Name:    "lambda-zip",
							Aliases: []string{"s"},
							Usage:   "[Advanced] Specify an external lambda payload to use",
						},
						&cli.StringFlag{
							Name:    "lambda-file-name",
							Aliases: []string{"n"},
							Usage:   "[Advanced] Save lambda under a different name",
							Value:   "kilt.zip",
						},
						&cli.PathFlag{
							Name:  "override-cfn-template",
							Usage: "[Advanced] Override the template used to deploy the macro",
						},
						&cli.StringFlag{
							Name:  "recipe-config",
							Usage: "[Advanced] Extra variables passed to the kilt definition (in json format)",
							Value: "{}",
						},
					},
					Action: func(c *cli.Context) error {
						cfg.Region = c.String("region")
						if c.Args().Len() != 2 {
							fmt.Printf("you need to specify 2 arguments to install\n")
							return cli.ShowAppHelp(c)
						}

						bucket, err := initializeAwsBucket(cfg)
						if err != nil {
							return cli.Exit("could not get aws bucket "+err.Error(), 1)
						}
						cfnMacro := cfnmacro.NewCfnMacroRuntime(cfg, bucket)

						var lambdaZip io.ReadCloser
						if c.IsSet("lambda-zip") {
							lambdaZip, err = os.Open(c.Path("lambda-zip"))
						} else {
							lambdaZip, err = pkger.Open("/artifacts/kilt-cfn-macro.zip")
						}
						if err != nil {
							return cli.Exit(fmt.Sprintf("could not open lambda zip: %s", err.Error()), 1)
						}
						defer lambdaZip.Close()

						var cfnTemplate io.ReadCloser
						if c.IsSet("override-cfn-template") {
							cfnTemplate, err = os.Open(c.Path("override-cfn-template"))
						} else {
							cfnTemplate, err = pkger.Open("/cmd/kilt-installer/kilt.yaml")
						}
						if err != nil {
							return cli.Exit(fmt.Sprintf("could not open cfn template: %s", err.Error()), 1)
						}
						defer cfnTemplate.Close()

						kiltDefinition, err := os.Open(c.Args().Get(0))
						if err != nil {
							return cli.Exit("could not open kilt definition: "+err.Error(), 1)
						}
						defer kiltDefinition.Close()

						throwaway := make(map[string]interface{})
						err = json.Unmarshal([]byte(c.String("recipe-config")), &throwaway)
						if err != nil {
							return cli.Exit("invalid recipe config specified: "+err.Error(), 1)
						}

						params := &cfnmacro.InstallationParameters{
							KiltDefinition:     kiltDefinition,
							MacroName:          c.Args().Get(1),
							OptIn:              c.Bool("opt-in"),
							LambdaZip:          lambdaZip,
							ZipDestinationName: c.String("lambda-file-name"),
							CfnTemplate:        cfnTemplate,
							ModelBuilder: func(input cfnmacro.TemplateDefaultModel) interface{} {
								return input
							},
							RecipeConfig: c.String("recipe-config"),
						}

						err = cfnMacro.InstallMacro(params)
						if err != nil {
							return cli.Exit("could not install macro: "+err.Error(), 2)
						}
						return nil
					},
				},
			},
		},
	}
}

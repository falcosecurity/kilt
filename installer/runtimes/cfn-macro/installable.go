package cfn_macro

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/falcosecurity/kilt/installer"
	"github.com/urfave/cli/v2"
)

const (
	macroPrefix = "KiltMacro"
	s3MacroPrefix = "cfn-macro-"
	macroS3Suffix = ".kilt.cfg"

)

type CfnMacroInstaller struct {}


type TemplateDefaultModel struct {
	BucketName string
	MacroName string
	MacroFileName string
	OptIn bool
	KiltZipPath string
	KiltKmsSecret string
}
type TemplateModelBuilderFunc func(input TemplateDefaultModel) interface{}

type Config struct {
	TemplateFile string
	TemplateModelBuilder TemplateModelBuilderFunc
}

func (c *CfnMacroInstaller) GetCommands(hooks *installer.Hooks) []*cli.Command {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		fmt.Printf("could not load AWS config:\n")
		panic(err)
	}

	return []*cli.Command{
		{
			Name:                   "cfn-macro",
			Usage:            		"Uses cloud formation macros to alter ECS Task Definitions",
			Category:               "runtimes",
			Subcommands:            []*cli.Command{
				{
					Name: "list",
					Usage: "Lists installed CFN macros",
					Action: listMacros(cfg),
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "region",
							Aliases: []string{"r"},
							EnvVars: []string{"AWS_REGION"},
							DefaultText: cfg.Region,
						},
					},
				},
				{
					Name: "delete",
					Usage: "Deletes CFN macros",
					Action: deleteMacros(cfg),
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "region",
							Aliases: []string{"r"},
							EnvVars: []string{"AWS_REGION"},
							DefaultText: cfg.Region,
						},
					},
				},
				{
					Name: "install",
					Usage: "Installs a new kilt CFN macro",
					Description: "Creates a CFN macro named MACRO_NAME that applies KILT_DEFINITION to CFN templates annotated with Transform: MACRO_NAME",
					ArgsUsage: "KILT_DEFINITION MACRO_NAME",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "opt-in",
							Aliases: []string{"o"},
							Usage: "Use opt-in logic instead of default opt-out",
						},
						&cli.StringFlag{
							Name: "region",
							Aliases: []string{"r"},
							Usage: "Specify region for installation",
							EnvVars: []string{"AWS_REGION"},
							DefaultText: cfg.Region,
							Required: true,
						},
						&cli.PathFlag{
							Name: "lambda-zip",
							Aliases: []string{"s"},
							Usage: "[Advanced] Specify an external lambda payload to use",
						},
						&cli.StringFlag{
							Name: "lambda-file-name",
							Aliases: []string{"n"},
							Usage: "[Advanced] Save lambda under a different name",
						},
						&cli.PathFlag{
							Name: "override-cfn-template",
							Usage: "[Advanced] Override the template used to deploy the macro",
						},
						&cli.StringFlag{
							Name: "kms-secret-arn",
							Usage: "[Advanced] Use a kms secret to pull images",
						},
					},
					Action: installMacro(cfg, hooks),
				},
			},
		},
	}
}


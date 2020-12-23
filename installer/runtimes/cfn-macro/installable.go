package cfn_macro

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	macroPrefix   = "KiltMacro"
	s3MacroPrefix = "cfn-macro-"
	macroS3Suffix = ".kilt.cfg"
)

type CfnMacroInstaller struct {
	awsConfig         aws.Config
	awsKiltBucketName string
}

type InstallationParameters struct {
	KiltDefinition io.ReadCloser
	MacroName      string
	OptIn          bool
	RecipeConfig   string

	LambdaZip          io.ReadCloser
	ZipDestinationName string

	CfnTemplate  io.ReadCloser
	ModelBuilder TemplateModelBuilderFunc
}

type TemplateDefaultModel struct {
	BucketName    string
	MacroName     string
	MacroFileName string
	OptIn         bool
	KiltZipPath   string
	RecipeConfig  string
}
type TemplateModelBuilderFunc func(input TemplateDefaultModel) interface{}

type Config struct {
	TemplateFile         string
	TemplateModelBuilder TemplateModelBuilderFunc
}

func NewCfnMacroRuntime(cfg aws.Config, bucketName string) *CfnMacroInstaller {
	return &CfnMacroInstaller{
		cfg,
		bucketName,
	}
}

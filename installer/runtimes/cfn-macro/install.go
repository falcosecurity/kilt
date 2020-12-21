package cfn_macro

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/kilt/installer"
	"github.com/falcosecurity/kilt/installer/util"
	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

func buildMacro(model TemplateDefaultModel) interface{} {
	return model
}

func installMacro(cfg aws.Config, hooks *installer.Hooks) cli.ActionFunc {
	pkger.Include("/runtimes/cfn-macro/kilt.yaml")
	var installerConfig *Config
	switch v := hooks.OverrideConfig.(type) {
	case Config:
		installerConfig = &v
	default:
		installerConfig = &Config{
			TemplateFile: "/runtimes/cfn-macro/kilt.yaml",
			TemplateModelBuilder: buildMacro,
		}
	}
	

	return func(c *cli.Context) error {
		if c.Args().Len() != 2 {
			 return cli.Exit("unexpected number of arguments", 1)
		}
		macroName := c.Args().Get(1)
		kiltDefinition, err := os.Open(c.Args().Get(0))
		if err != nil {
			return cli.Exit("cannot open file " +  c.Args().Get(0) + ": " + err.Error(), 1)
		}
		defer kiltDefinition.Close()

		region := c.String("region")
		if region != "" {
			cfg.Region = region
		}

		if hooks.PreInstall != nil {
			err = hooks.PreInstall(c)
			if err != nil {
				return cli.Exit("could not execute pre-install hook: " + err.Error(), 1)
			}
		}

		s3c := s3.NewFromConfig(cfg)
		cfnc := cloudformation.NewFromConfig(cfg)
		awsUtil := util.New(cfg)

		lambdaFileName := "kilt.zip"
		if c.IsSet("lambda-file-name"){
			lambdaFileName = c.String("lambda-file-name")
		}

		fmt.Printf("AWS configuration loaded account %s in %s\n", awsUtil.AwsAccountID, awsUtil.AwsRegionName)
		bucket, err := awsUtil.GetOrCreateKiltS3Bucket(s3c)
		if err != nil {
			return cli.Exit("could not find aws s3 bucket: " + err.Error(), 2)
		}

		fmt.Printf("Using AWS S3 bucket %s\n", bucket)

		uploader := manager.NewUploader(s3c)
		var f io.ReadCloser
		if c.IsSet("lambda-zip") {
			f, err = os.Open(c.Path("lambda-zip"))
		}else{
			f, err = pkger.Open("/artifacts/kilt-cfn-macro.zip")
		}
		if err!=nil {
			return cli.Exit("could not open embedded cfn-macro lambda: " + err.Error(), 3)
		}
		defer f.Close()

		fmt.Printf("Uploading lambda code. This might take a while...")
		_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key: aws.String(lambdaFileName),
			Body: f,
		})
		if err != nil {
			return cli.Exit("could not upload macro to s3: "+ err.Error(), 4)
		}

		fmt.Printf("DONE\n")

		kiltDefinitionFile := s3MacroPrefix + macroName + macroS3Suffix
		fmt.Printf("Uploading kilt definition %s...", kiltDefinitionFile)
		_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key: aws.String(kiltDefinitionFile),
			Body: kiltDefinition,
		})
		if err != nil {
			return cli.Exit("could not upload kilt definition to s3: " + err.Error(), 4)
		}
		fmt.Printf("DONE\n")

		var k io.ReadCloser
		if c.IsSet("override-cfn-template") {
			k, err = os.Open(c.Path("override-cfn-template"))
		}else{
			k, err = pkger.Open(installerConfig.TemplateFile)
		}
		if err != nil {
			return cli.Exit("could not read kilt definition: " + err.Error(), 5)
		}
		defer k.Close()

		data, err := ioutil.ReadAll(k)
		if err != nil {
			return cli.Exit("could not read kilt definition: " + err.Error(), 5)
		}

		t, err := template.New("kilt-definition").Parse(string(data))
		if err != nil {
			return cli.Exit("could not parse kilt template: " + err.Error(), 5)
		}


		var buf bytes.Buffer
		err = t.Execute(&buf, installerConfig.TemplateModelBuilder(TemplateDefaultModel{
			BucketName:    bucket,
			MacroName:     macroName,
			MacroFileName: kiltDefinitionFile,
			OptIn:         c.Bool("opt-in"),
			KiltZipPath:   lambdaFileName,
			KiltKmsSecret: c.String("kms-secret-arn"),
		}))
		if err != nil {
			return cli.Exit("could not compute macro template: " + err.Error(), 6)
		}

		stackName := macroPrefix + macroName
		fmt.Printf("Submitted cloudformation stack '%s'. Follow progress in AWS console\n", stackName)
		_, err = cfnc.CreateStack(context.Background(), &cloudformation.CreateStackInput{
			StackName: &stackName,
			TemplateBody: aws.String(buf.String()),
			Capabilities: []types.Capability{
				types.CapabilityCapabilityIam,
			},
		})
		if err != nil {
			return cli.Exit("could not create cfn macro: "+ err.Error(), 7)
		}

		if hooks.PostInstall != nil {
			err = hooks.PostInstall(c)
			if err != nil {
				return cli.Exit("could not execute post-install hook: " + err.Error(), 1)
			}
		}

		return nil
	}
}

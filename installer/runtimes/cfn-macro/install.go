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
	"io/ioutil"
	"text/template"
)

func (r *CfnMacroInstaller) InstallMacro(params *InstallationParameters) error {
	s3c := s3.NewFromConfig(r.awsConfig)
	cfnc := cloudformation.NewFromConfig(r.awsConfig)

	uploader := manager.NewUploader(s3c)

	fmt.Printf("Uploading lambda code. This might take a while...")
	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &r.awsKiltBucketName,
		Key: &params.ZipDestinationName,
		Body: params.LambdaZip,
	})
	if err != nil {
		fmt.Printf("ERROR!\n")
		return fmt.Errorf("could not upload macro to s3: %w", err)
	}
	fmt.Printf("DONE\n")

	kiltDefinitionFile := s3MacroPrefix + params.MacroName + macroS3Suffix
	fmt.Printf("Uploading kilt definition %s...", kiltDefinitionFile)
	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &r.awsKiltBucketName,
		Key: aws.String(kiltDefinitionFile),
		Body: params.KiltDefinition,
	})
	if err != nil {
		return fmt.Errorf("could not upload kilt definition to s3: %w", err)
	}
	fmt.Printf("DONE\n")

	data, err := ioutil.ReadAll(params.CfnTemplate)
	if err != nil {
		return fmt.Errorf("could not read CFN Template: %w", err)
	}

	t, err := template.New("kilt-definition").Parse(string(data))
	if err != nil {
		return fmt.Errorf("could not parse CFN template: %w" , err)
	}


	var buf bytes.Buffer
	err = t.Execute(&buf, params.ModelBuilder(TemplateDefaultModel{
		BucketName:    r.awsKiltBucketName,
		MacroName:     params.MacroName,
		MacroFileName: kiltDefinitionFile,
		OptIn:         params.OptIn,
		KiltZipPath:   params.ZipDestinationName,
		RecipeConfig:  params.RecipeConfig,
	}))
	if err != nil {
		return fmt.Errorf("could not compute macro template: %w" , err)
	}

	stackName := macroPrefix + params.MacroName
	_, err = cfnc.CreateStack(context.Background(), &cloudformation.CreateStackInput{
		StackName: &stackName,
		TemplateBody: aws.String(buf.String()),
		Capabilities: []types.Capability{
			types.CapabilityCapabilityIam,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create cfn macro: %w", err)
	}
	fmt.Printf("Submitted cloudformation stack '%s'. Follow creation progress in AWS console\n", stackName)
	return nil
}

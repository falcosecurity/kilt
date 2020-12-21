package cfn_macro

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/kilt/installer/util"
	"github.com/urfave/cli/v2"
)

func deleteMacros(cfg aws.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		region := c.String("region")
		if region != "" {
			cfg.Region = region
		}
		awsUtil := util.New(cfg)

		s3c := s3.NewFromConfig(cfg)
		cfnc := cloudformation.NewFromConfig(cfg)

		bucket, err := awsUtil.GetOrCreateKiltS3Bucket(s3c)
		if err != nil {
			return fmt.Errorf("error creating/retrieving s3 bucket: %w", err)
		}

		for _, arg := range c.Args().Slice() {
			deleteKey := s3MacroPrefix + arg + macroS3Suffix
			_, err = s3c.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &deleteKey,
			})
			if err != nil {
				fmt.Printf("Warning: could not delete s3://%s/%s: %s\n", bucket, deleteKey, err.Error())
			}

			stackName := macroPrefix + arg
			_, err = cfnc.DeleteStack(context.Background(), &cloudformation.DeleteStackInput{
				StackName: &stackName,
			})
			if err != nil {
				fmt.Printf("Warning: could not delete cloud formation stack '%s'\n", stackName)
			}
			fmt.Printf("Done\n")
		}
		return nil
	}
}
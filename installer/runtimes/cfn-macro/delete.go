package cfn_macro

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *CfnMacroInstaller) Delete(macros []string) error {
	s3c := s3.NewFromConfig(r.awsConfig)
	cfnc := cloudformation.NewFromConfig(r.awsConfig)

	for _, arg := range macros {
		deleteKey := s3MacroPrefix + arg + macroS3Suffix
		_, err := s3c.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: &r.awsKiltBucketName,
			Key:    &deleteKey,
		})
		if err != nil {
			fmt.Printf("Warning: could not delete s3://%s/%s: %s\n", r.awsKiltBucketName, deleteKey, err.Error())
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

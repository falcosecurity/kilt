package util

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const bucketTemplate = "kilt-{{ .AwsAccountID }}-{{ .AwsRegionName }}"

func GetAwsAccountName(cfg aws.Config, stsc *sts.Client) (string, error) {
	if stsc == nil {
		stsc = sts.NewFromConfig(cfg)
	}
	output, err := stsc.GetCallerIdentity(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot get own account id: %w", err)
	}
	return *output.Account, nil
}

func GetRegion(cfg aws.Config) (string, error) {
	if cfg.Region == "" {
		return "", fmt.Errorf("region is empty: please configure the region using AWS_REGION environment variable or awscli config")
	}
	return cfg.Region, nil
}

func GetBucketName(accountId string, region string) (string, error) {
	var buf bytes.Buffer
	t, err := template.New("kilt-bucket").Parse(bucketTemplate)
	if err != nil {
		return "", fmt.Errorf("cannot parse const bucket template: %w", err)
	}
	err = t.Execute(&buf, struct {
		AwsAccountID  string
		AwsRegionName string
	}{
		accountId,
		region,
	})
	if err != nil {
		return "", fmt.Errorf("cannot execute const bucket template: %w", err)
	}
	return buf.String(), nil
}

func EnsureBucketExists(bucketName string, s3Client *s3.Client) error {
	_, err := s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		_, err = s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: &bucketName,
		})
		if err != nil {
			return fmt.Errorf("could not create S3 bucket %s: %w", bucketName, err)
		}
		be := s3.NewBucketExistsWaiter(s3Client)
		err = be.Wait(context.Background(), &s3.HeadBucketInput{
			Bucket: &bucketName,
		}, 10 * time.Second)
		if err != nil {
			return fmt.Errorf("timed out while waiting the bucket %s to be created: %w", bucketName, err)
		}
	}
	return nil
}

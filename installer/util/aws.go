package util

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const bucketTemplate = "kilt-{{ .AwsAccountID }}-{{ .AwsRegionName }}"

type AwsUtil struct {
	config aws.Config

	AwsAccountID  string
	AwsRegionName string
	kiltBucket    string
}

func New(config aws.Config) AwsUtil {
	r := AwsUtil{
		config: config,
	}
	r.GetAwsAccountName()
	r.GetRegion()

	return r
}

func (r *AwsUtil) GetAwsAccountName() (string, error) {
	if r.AwsAccountID == "" {
		stsc := sts.NewFromConfig(r.config)
		output, err := stsc.GetCallerIdentity(context.TODO(), nil)
		if err != nil {
			return "", fmt.Errorf("cannot get own account id: %w", err)
		}
		r.AwsAccountID = *output.Account
	}

	return r.AwsAccountID, nil
}

func (r *AwsUtil) GetRegion() (string, error) {
	if r.AwsRegionName == "" {
		if r.config.Region == "" {
			return "", fmt.Errorf("region is empty: please configure the region using AWS_REGION environment variable or awscli config")
		}
		r.AwsRegionName = r.config.Region
	}

	return r.AwsRegionName, nil
}

func (r *AwsUtil) GetOrCreateKiltS3Bucket(s3Client *s3.Client) (string, error) {
	if r.kiltBucket == "" {
		var buf bytes.Buffer
		if s3Client == nil {
			s3Client = s3.NewFromConfig(r.config)
		}

		// template is a const so it cannot fail
		t, err := template.New("kilt-bucket").Parse(bucketTemplate)
		if err != nil {
			return "", fmt.Errorf("cannot parse const bucket template: %w", err)
		}
		err = t.Execute(&buf, r)
		if err != nil {
			return "", fmt.Errorf("cannot execute const bucket template: %w", err)
		}

		bucket := buf.String()

		_, err = s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
			Bucket: &bucket,
		})
		if err != nil {
			_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
				Bucket: &bucket,
			})
			if err != nil {
				return "", fmt.Errorf("could not create S3 bucket %s: %w", bucket, err)
			}
		}
		r.kiltBucket = bucket
	}

	return r.kiltBucket, nil
}

package cfn_macro

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/kilt/installer/util"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

func listMacros(cfg aws.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		region := c.String("region")
		if region != "" {
			cfg.Region = region
		}

		s3c := s3.NewFromConfig(cfg)
		awsUtil := util.New(cfg)

		bucket, err := awsUtil.GetOrCreateKiltS3Bucket(s3c)
		prefix := s3MacroPrefix
		if err != nil {
			return fmt.Errorf("error while creating/retrieving s3 bucket: %w", err)
		}

		p := s3.NewListObjectsV2Paginator(s3c, &s3.ListObjectsV2Input{
			Bucket: &bucket,
			Prefix: &prefix,
		})

		for p.HasMorePages() {
			page, err := p.NextPage(context.Background())
			if err != nil {
				return fmt.Errorf("cannot list s3 objects in page: %w", err)
			}
			for _, obj := range page.Contents {
				macroName := strings.TrimPrefix(*obj.Key, prefix)
				macroName = strings.TrimSuffix(macroName, macroS3Suffix)
				fmt.Printf("%s - Last Modified: %s\n", macroName, obj.LastModified.Format(time.RFC3339))
			}
		}


		return nil
	}
}

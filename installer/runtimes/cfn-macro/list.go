package cfn_macro

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
	"time"
)

func (r *CfnMacroInstaller) List() error {
	s3c := s3.NewFromConfig(r.awsConfig)
	prefix := s3MacroPrefix

	p := s3.NewListObjectsV2Paginator(s3c, &s3.ListObjectsV2Input{
		Bucket: &r.awsKiltBucketName,
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

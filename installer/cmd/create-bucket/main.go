package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/kilt/installer/util"
)

func main() {
	c, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Errorf("could not read aws config: %w", err))
	}

	s3c := s3.NewFromConfig(c)

	err = util.EnsureBucketExists(os.Args[1], os.Args[2], s3c)
	if err != nil {
		panic(err)
	}
}

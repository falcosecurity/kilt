package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/falcosecurity/kilt/utilities/logshipper/cwlogger"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(fmt.Errorf("could not load aws config: %w", err))
	}
	cw := cloudwatchlogs.New(cfg)
	logger, err := cwlogger.New(&cwlogger.Config{
		Client: cw,
		LogGroupName: os.Getenv("__CW_LOG_GROUP"),
	})
	if err != nil {
		panic(fmt.Errorf("could not open cw logs: %w", err))
	}

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for scanner.Scan() {
		logger.Log(time.Now(), scanner.Text())
	}
}

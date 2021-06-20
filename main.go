package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-2"))
	if err != nil {
		panic("aws configuration error, " + err.Error())
	}

	iamFinder := NewIAMFinderFromAWSConfig(cfg)
	oldAccessKeyList := iamFinder.ListAllOldAccessKey(5 * 24)

	sendOldAccessKeys(oldAccessKeyList, "https://hooks.slack.com/services/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")
}


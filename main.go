package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-2"))
	if err != nil {
		panic("aws configuration error, " + err.Error())
	}

	iamFinder := NewIAMFinderFromAWSConfig(cfg)
	oldAccessKeyList := iamFinder.ListAllOldAccessKey(10)

	for _, accessKeyMetadata := range oldAccessKeyList {
		d, _ := json.MarshalIndent(accessKeyMetadata, "", "")
		fmt.Printf("%s\n", d)
	}

}


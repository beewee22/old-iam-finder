package main

import (
	"context"
	env "github.com/Netflix/go-env"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Environment struct {
	AWSAccessKey string `env:"AWS_ACCESS_KEY_ID,required=true"`
	AWSSecretKey string `env:"AWS_SECRET_ACCESS_KEY,required=true"`
	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL,required=true"`
	ExpireHour int `env:"EXPIRE_HOUR,required=true"`

	Extras env.EnvSet
}

func main() {
	var environment Environment
	extras, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		panic("environment initialize error, " + err.Error())
	}
	environment.Extras = extras

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-2"))
	if err != nil {
		panic("aws configuration error, " + err.Error())
	}

	iamFinder := NewIAMFinderFromAWSConfig(cfg)
	oldAccessKeyList := iamFinder.ListAllOldAccessKey(environment.ExpireHour)

	if oldAccessKeyList != nil {
		sendOldAccessKeys(oldAccessKeyList, environment.SlackWebhookURL)
	}
}


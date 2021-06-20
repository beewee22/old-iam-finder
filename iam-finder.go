package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMFinder struct {
	iam *iam.Client
}

//NewIAMFinderFromAWSConfig Create new IAMFinder from config
func NewIAMFinderFromAWSConfig(cfg aws.Config) *IAMFinder {
	IAMClient := &IAMFinder{
		iam: iam.NewFromConfig(cfg),
	}

	return IAMClient
}

//ListAllIAMUsers List all IAM users from account
func (f IAMFinder) ListAllIAMUsers() *iam.ListUsersOutput {
	userList, err := f.iam.ListUsers(context.TODO(), &iam.ListUsersInput{})

	if err != nil {
		fmt.Println("Got an error retrieving users:")
		panic(err)
	}

	return userList
}

//ListAllAccessKeyMetadata List all IAM access key metadata
func (f IAMFinder) ListAllAccessKeyMetadata() (accessKeys []types.AccessKeyMetadata) {
	IAMList := f.ListAllIAMUsers()

	for _, user := range IAMList.Users {
		accessKey, err := f.iam.ListAccessKeys(
			context.TODO(),
			&iam.ListAccessKeysInput{UserName: user.UserName},
		)
		if err != nil {
			fmt.Println("Got an error retrieving access key from user" + *user.UserName)
			fmt.Printf("%v", err)
			continue
		}
		if accessKey != nil {
			accessKeys = append(accessKeys, accessKey.AccessKeyMetadata...)
		}
	}
	return
}

type OldAccessKeyInfo struct {
	UserName string
	AccessKeyId string
	LastUsed time.Time
	IsNotUsed bool
}

const AccesskeyPrefix = "AKIA"

func (f IAMFinder) ListAllOldAccessKey(expireHours int) (oldAccessKeyList []OldAccessKeyInfo) {
	// calc expected at least time
	expireDuration := time.Duration(expireHours) * time.Hour * -1
	expectedLastUsed := time.Now().Add(expireDuration)

	accessKeyMetadataList := f.ListAllAccessKeyMetadata()

	for _, accessKeyMetadata := range accessKeyMetadataList {
		// access key must starts with "AKIA"
		if strings.HasPrefix(*accessKeyMetadata.AccessKeyId, AccesskeyPrefix) {
			continue
		}
		last, err := f.iam.GetAccessKeyLastUsed(context.TODO(), &iam.GetAccessKeyLastUsedInput{
			AccessKeyId: accessKeyMetadata.AccessKeyId,
		})
		if err != nil {
			fmt.Println("Got an error retrieving access key last used time:")
			fmt.Printf("%v", accessKeyMetadata)
		}
		if last != nil {
			var lastUsedTime *time.Time
			isNotUsed := false
			// check if access key was not used
			if last.AccessKeyLastUsed.LastUsedDate == nil {
				// if the access key was not used at all, check with created date
				lastUsedTime = accessKeyMetadata.CreateDate
				isNotUsed = true
			} else {
				lastUsedTime = last.AccessKeyLastUsed.LastUsedDate
			}
			// check last used before expected time
			if lastUsedTime.Before(expectedLastUsed) {
				oldAccessKeyList = append(oldAccessKeyList, OldAccessKeyInfo{
					UserName:  *last.UserName,
					AccessKeyId: *accessKeyMetadata.AccessKeyId,
					LastUsed: *lastUsedTime,
					IsNotUsed: isNotUsed,
				})
			}
		}
	}

	return
}


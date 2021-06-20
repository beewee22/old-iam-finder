package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCreateSlackMessageFromAccessKeyList(t *testing.T) {
	type args struct {
		metadataList []OldAccessKeyInfo
	}
	mockOldAccessKeys := []OldAccessKeyInfo{
		{UserName: "Test1", AccessKeyId: "ACCESS_KEY_ID_1", LastUsed: time.Now(), IsNotUsed: false},
		{UserName: "Test2", AccessKeyId: "ACCESS_KEY_ID_2", LastUsed: time.Now().Add(-1 * time.Hour), IsNotUsed: true},
		{UserName: "Test3", AccessKeyId: "ACCESS_KEY_ID_3", LastUsed: time.Now().Add(-2 * time.Hour), IsNotUsed: false},
	}
	expectedSlackLayout := SlackAccessKeyMessageLayout{Blocks: []Block{
		{Type: "section", Text: &Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n*Username*: %s\n*AccessKey*: %s\n*Last used*: %s\n", mockOldAccessKeys[0].UserName, mockOldAccessKeys[0].AccessKeyId, mockOldAccessKeys[0].LastUsed.Format("2006-01-02 15:04:05")),
		}, Elements: nil}, {Type: "section", Text: &Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n*Username*: %s\n*AccessKey*: %s\n*Last used*: %s\n", mockOldAccessKeys[1].UserName, mockOldAccessKeys[1].AccessKeyId, mockOldAccessKeys[1].LastUsed.Format("2006-01-02 15:04:05")),
		}, Elements: nil}, {Type: "section", Text: &Text{
			Type: "mrkdwn",
			Text: fmt.Sprintf("\n*Username*: %s\n*AccessKey*: %s\n*Last used*: %s\n", mockOldAccessKeys[2].UserName, mockOldAccessKeys[2].AccessKeyId, mockOldAccessKeys[2].LastUsed.Format("2006-01-02 15:04:05")),
		}, Elements: nil},
	}}
	tests := []struct {
		name              string
		args              args
		wantMessageLayout SlackAccessKeyMessageLayout
	}{
		{name: "should return expected slack message layout", args: struct{ metadataList []OldAccessKeyInfo }{metadataList: mockOldAccessKeys}, wantMessageLayout: expectedSlackLayout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMessageLayout := CreateSlackMessageFromAccessKeyList(tt.args.metadataList); !reflect.DeepEqual(gotMessageLayout, tt.wantMessageLayout) {
				got, _ := json.MarshalIndent(gotMessageLayout, "", "  ")
				expect, _ := json.MarshalIndent(tt.wantMessageLayout, "", "  ")
				t.Errorf("CreateSlackMessageFromAccessKeyList() = %+v, want %+v", string(got), string(expect))
			}
		})
	}
}

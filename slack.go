package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackAccessKeyMessageLayout struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string    `json:"type"`
	Text     *Text     `json:"text,omitempty"`
	Elements *[]Element `json:"elements,omitempty"`
}

type Element struct {
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
	Text  string `json:"text"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

//CreateSlackMessageFromAccessKeyList Create Slack message layout from old access key metadata
func CreateSlackMessageFromAccessKeyList(metadataList []OldAccessKeyInfo) (messageLayout SlackAccessKeyMessageLayout){
	var blocks []Block
	for _, metadata := range metadataList {
		text := fmt.Sprintf("\n*Username*: %s\n*AccessKey*: %s\n*Last used*: %s\n", metadata.UserName, metadata.AccessKeyId, metadata.LastUsed.Format("2006-01-02 15:04:05"))
		blocks = append(blocks, Block{
			Type:     "section",
			Text:     &Text{
				Type: "mrkdwn",
				Text: text,
			},
			Elements: nil,
		})
	}
	
	messageLayout.Blocks = blocks

	return
}

//sendOldAccessKeys Send old access key to slack
func sendOldAccessKeys(metadata []OldAccessKeyInfo, webhookURL string) {
	// create message payload
	slackMessageLayout := CreateSlackMessageFromAccessKeyList(metadata)
	payload, err := json.Marshal(slackMessageLayout)
	if err != nil {
		println("slack payload json marshal error")
		fmt.Printf("%+v\n", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		println("slack http request error on creation")
		fmt.Printf("%+v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("slack http request failed")
		panic(err)
	}

	_ = resp.Body.Close()
}
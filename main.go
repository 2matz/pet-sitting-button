package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/2matz/soracom-ltem-button-go-handler/oneclick"
	"github.com/2matz/soracom-ltem-button-go-handler/slack"
	"github.com/aws/aws-lambda-go/lambda"
)

var logger log.Logger

// HandleRequest ...
func HandleRequest(ctx context.Context, event oneclick.Event) (result string, err error) {
	var slackWebhookURL *url.URL
	attributes := event.GetPlacementAttributes()
	u := attributes.(map[string]interface{})["slack_webhook_url"].(string)
	if u == "" {
		result = ""
		err = fmt.Errorf("%s", "Slack Webhook URL is not defined.")
		return
	}
	slackWebhookURL, err = url.ParseRequestURI(u)
	if err != nil {
		result = ""
		err = fmt.Errorf("%s", "Slack Webhook URL is incorrect format.")
		return
	}

	clickType, err := event.GetClickType()
	if err != nil {
		result = ""
		err = fmt.Errorf("%s", "Click type is not defined.")
		return
	}

	userName := event.GetProjectName()
	if userName == "" {
		userName = "SORACOM LTE-M Button"
	}

	var message string

	switch clickType {
	case oneclick.SingleClick:
		message = "ご飯の用意ができましたよ〜:rice::heart_eyes_cat:"
	case oneclick.DoubleClick:
		message = "トイレが綺麗になりました:smile_cat:"
	case oneclick.LongClick:
		message = "緊急事態です!:scream_cat::rotating_light:"
	}

	slackWebHookClient := slack.NewSlackWebhookClient(
		ctx,
		slackWebhookURL.String(),
		message,
		userName,
		"",
		"",
		"",
	)
	result, err = slackWebHookClient.Post()

	log.Printf("result: %v, error: %v", result, err)

	return
}

func main() {
	lambda.Start(HandleRequest)
}

package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/2matz/soracom-ltem-button-go-handler/oneclick"
	"github.com/2matz/soracom-ltem-button-go-handler/slack"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	singleClickMessage string
	doubleClickMessage string
	longClickMessage   string
)

const (
	defaultSingleClickMessage   = "single clicked!"
	defaultDoubleClickMessage   = "double clicked!"
	defaultLongClickMessage     = "long clicked!"
	unsupportedClickTypeMessage = "unsupported click type"
)

// HandleRequest ...
func HandleRequest(ctx context.Context, event oneclick.Event) {
	var slackWebhookURL *url.URL
	var err error
	attributes := event.GetPlacementAttributes().(map[string]interface{})

	u := attributes["slack_webhook_url"].(string)
	if u == "" {
		log.Println("Slack Webhook URL is not defined.")
		return
	}
	slackWebhookURL, err = url.ParseRequestURI(u)
	if err != nil {
		log.Printf("Slack Webhook URL is incorrect format. %v", err)
		return
	}

	if singleClickMessage == "" {
		singleClickMessage = attributes["single_click_message"].(string)
		if singleClickMessage == "" {
			singleClickMessage = defaultSingleClickMessage
		}
	}
	if doubleClickMessage == "" {
		doubleClickMessage = attributes["double_click_message"].(string)
		if doubleClickMessage == "" {
			doubleClickMessage = defaultDoubleClickMessage
		}
	}
	if longClickMessage == "" {
		longClickMessage = attributes["long_click_message"].(string)
		if longClickMessage == "" {
			longClickMessage = defaultLongClickMessage
		}
	}
	clickType, err := event.GetClickType()
	if err != nil {
		log.Printf("Click type is not defined. %v", err)
		return
	}

	userName := event.GetProjectName()
	if userName == "" {
		userName = "SORACOM LTE-M Button"
	}

	var message string

	switch clickType {
	case oneclick.SingleClick:
		message = singleClickMessage
	case oneclick.DoubleClick:
		message = doubleClickMessage
	case oneclick.LongClick:
		message = longClickMessage
	default:
		message = unsupportedClickTypeMessage
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
	_, err = slackWebHookClient.Post()
	if err != nil {
		log.Println(err)
	}
	return
}

func init() {
	singleClickMessage = os.Getenv("single_click_message")
	doubleClickMessage = os.Getenv("double_click_message")
	longClickMessage = os.Getenv("long_click_message")
}

func main() {
	lambda.Start(HandleRequest)
}

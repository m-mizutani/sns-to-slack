package main

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	lambda.Start(handler)

}

type slackMessage struct {
	Attachments []slackAttachment `json:"attachments,omitempty"`
	Text        string            `json:"text"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	UserName    string            `json:"username,omitempty"`
}

type slackAttachment struct {
	Actions        []slackAction `json:"actions,omitempty"`
	AttachmentType string        `json:"attachment_type,omitempty"`
	AuthorIcon     string        `json:"author_icon,omitempty"`
	AuthorName     string        `json:"author_name,omitempty"`
	CallbackID     string        `json:"callback_id,omitempty"`
	Color          string        `json:"color,omitempty"`
	Fallback       string        `json:"fallback,omitempty"`
	Fields         []slackField  `json:"fields,omitempty"`
	ImageURL       string        `json:"image_url,omitempty"`
	Text           string        `json:"text,omitempty"`
	Title          string        `json:"title,omitempty"`
}

type slackAction struct {
	Name  string `json:"name,omitempty"`
	Text  string `json:"text,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type slackField struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
	logger.WithField("event", snsEvent).Info("Start handler")

	incomingURL := os.Getenv("WEBHOOK_URL")
	if incomingURL == "" {
		return errors.New("WEBHOOK_URL is required")
	}

	for _, record := range snsEvent.Records {
		logger.WithField("record", record).Info("Start record")

		var msg slackMessage
		msg.Text = record.SNS.Message
		buf, err := json.Marshal(msg)
		if err != nil {
			return errors.Wrap(err, "Fail to marshal slack message")
		}

		resp, err := http.Post(incomingURL, "application/json", bytes.NewReader(buf))
		if err != nil {
			return errors.Wrap(err, "Fail to post message")
		}

		logger.WithField("resp", resp).Debug("done HTTP request")
	}

	logger.Info("Exit handler")
	return nil
}

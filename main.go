package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	lambda.Start(handler)

}

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
	logger.WithField("event", snsEvent).Info("Start handler")

	incomingURL := os.Getenv("WEBHOOK_URL")
	if incomingURL == "" {
		return errors.New("WEBHOOK_URL is required")
	}

	logger.Info("Exit handler")
	return nil
}

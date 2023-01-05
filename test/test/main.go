package main

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	mailgun "github.com/mailgun/mailgun-go/v3"

	"log"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.SQSEvent) error {

	for _, message := range request.Records {

		err := SendEmail(context.Background(), message.Body)

		if err != nil {
			log.Printf("Failed to send email. Error: %v\n", err)
			return err
		}
	}

	return nil
}

// Send email via mailgun
func SendEmail(cont context.Context, emailInfo string) error {

	// Load env variables for mailgun
	domain := os.Getenv("EMAIL_DOMAIN")
	apiKey := os.Getenv("EMAIL_APIKEY")
	senderEmail := os.Getenv("EMAIL_SENDER")

	MailGunInstanceCloud := mailgun.NewMailgun(domain, apiKey)
	MailGunInstanceCloud.SetAPIBase("https://api.eu.mailgun.net/v3")

	message := MailGunInstanceCloud.NewMessage(senderEmail, "test email lambda", "", emailInfo)
	message.SetTemplate(os.Getenv("EMAIL_TEMPLATE"))

	ctx, cancel := context.WithTimeout(cont, time.Second*time.Duration(10))
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := MailGunInstanceCloud.Send(ctx, message)

	if err != nil {
		log.Printf("Unable to send an email: %v\n", err)
		return err
	}

	log.Printf("Email Sent to %s: ID: [%s] Response: %s\n", emailInfo, id, resp)

	return nil
}

// sam local invoke -e events/event.json
// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

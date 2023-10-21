package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
)

var (
	sesSession      *ses.SES
	sentFromAddress string
)

func init() {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	sentFromAddress := os.Getenv("SENT_FROM_ADDRESS")

	if accessKey == "" || secretKey == "" || region == "" || sentFromAddress == "" {
		log.Fatal("AWS credentials or region are not set")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}
	sesSession = ses.New(sess)
}

func main() {
	r := gin.Default()

	r.POST("/v1/send", sendEmail)

	r.Run(":8080")
}

func sendEmail(c *gin.Context) {
	var jsonInput struct {
		Recipient     string `json:"recipient"`
		Subject       string `json:"subject"`
		Base64Content string `json:"base64_content"`
	}
	if err := c.BindJSON(&jsonInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	decodedContent, err := base64.StdEncoding.DecodeString(jsonInput.Base64Content)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid base64 content"})
		return
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(jsonInput.Recipient)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(jsonInput.Subject),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(string(decodedContent)),
				},
			},
		},
		Source: aws.String(sentFromAddress),
	}

	_, err = sesSession.SendEmail(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Email sent successfully"})
}

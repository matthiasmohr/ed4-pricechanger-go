package instance

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/matthiasmohr/ed4-pricechanger-go/utils/logger"
)

func GetConnection() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		logger.PANIC("Fehler in der AWS-Konfiguration", err)
	}
	return dynamodb.New(sess)
}

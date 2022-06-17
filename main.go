package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type DemoApp struct {
	ctx       context.Context
	client    *sqs.Client
	region    string
	queueName string
}

var app DemoApp

func init() {
	app = DemoApp{
		ctx:       context.TODO(),
		queueName: "my-queue",
		region:    "eu-central-1",
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == sqs.ServiceID && region == app.region {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: app.region,
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(app.ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("localstack", "localstack", "session")),
	)
	logError(err)

	app.client = sqs.NewFromConfig(cfg)
}

func main() {

	app.listQueues()

	app.createQueue()
	app.listQueues()

	app.sendMessage("test message")

	time.Sleep(3 * time.Second)
	app.receiveMessage()

	app.deleteQueue()
	app.listQueues()

	fmt.Println("SQS")
}

func (app *DemoApp) listQueues() {
	fmt.Println("-------------- ListQueues")
	input := &sqs.ListQueuesInput{}

	result, err := app.client.ListQueues(app.ctx, input)
	logError(err)

	for i, qURL := range result.QueueUrls {
		fmt.Printf("%d: %s\n", i+1, qURL)
	}
}

func (app *DemoApp) createQueue() {
	fmt.Println("-------------- CreateQueue")

	input := &sqs.CreateQueueInput{
		QueueName: &app.queueName,
		Attributes: map[string]string{
			"DelaySeconds":           "60",
			"MessageRetentionPeriod": "3600",
		},
	}

	result, err := app.client.CreateQueue(app.ctx, input)
	logError(err)

	fmt.Println("Queue URL:", *result.QueueUrl)
}

func (app *DemoApp) getQueueURL(queueName string) (string, error) {
	result, err := app.client.GetQueueUrl(app.ctx, &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Println("Error getting queue URL", err)
		return "", err
	}

	return *result.QueueUrl, nil
}

func (app *DemoApp) sendMessage(message string) {
	fmt.Println("-------------- SendMessage")

	queueURL, err := app.getQueueURL(app.queueName)
	if err != nil {
		return
	}

	input := &sqs.SendMessageInput{
		DelaySeconds: 5,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Artist": {
				DataType:    aws.String("String"),
				StringValue: aws.String("RMHighlander"),
			},
			"SongTitle": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Odyssey"),
			},
		},
		MessageBody: aws.String(message),
		QueueUrl:    &queueURL,
	}

	resp, err := app.client.SendMessage(app.ctx, input)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Println("MessageId:" + *resp.MessageId)
}

func (app *DemoApp) receiveMessage() {
	fmt.Println("-------------- ReceiveMessage")

	queueURL, err := app.getQueueURL(app.queueName)
	if err != nil {
		return
	}

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   int32(10),
		QueueUrl:            &queueURL,
	}

	result, err := app.client.ReceiveMessage(app.ctx, input)
	if err != nil {
		fmt.Println("Receive message error:", err)
		return
	}

	fmt.Println(result.Messages)

	// fmt.Println("Message ID:     " + *result.Messages[0].MessageId)
	// fmt.Println("Message Handle: " + *result.Messages[0].ReceiptHandle)
}

func (app *DemoApp) deleteMessage(messageHandle string) {
	fmt.Println("-------------- DeleteMessage")

	queueURL, err := app.getQueueURL(app.queueName)
	if err != nil {
		return
	}

	input := &sqs.DeleteMessageInput{
		ReceiptHandle: &messageHandle,
		QueueUrl:      &queueURL,
	}
	_, err = app.client.DeleteMessage(app.ctx, input)
	if err != nil {
		fmt.Println("Error delete message:", err)
		return
	}

	fmt.Println("Message successfully deleted")
}

func (app *DemoApp) deleteQueue() {
	fmt.Println("-------------- DeleteQueue")

	queueURL, err := app.getQueueURL(app.queueName)
	if err != nil {
		return
	}

	input := &sqs.DeleteQueueInput{
		QueueUrl: &queueURL,
	}
	_, err = app.client.DeleteQueue(app.ctx, input)
	if err != nil {
		fmt.Println("Error deleting the queue:", err)
		return
	}

	fmt.Println("Deleted queue with URL " + queueURL)
}

func logError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

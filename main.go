package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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
	fmt.Println("-------------- ListQueues")

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

func logError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

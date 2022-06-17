SQS Kinesis Demo Project
-----

## SQS

Run SQS on LocalStack:
```
docker-compose -f docker-compose-sqs.yaml up -d
```

Health-check:
http://localhost:4566/health

Build and run the app:
```
go mod tidy

go run main.go
```


## Brief Overview

Queue attributes:
+ `QueueName`: name of the queue
+ Attributes:
    - `DelaySeconds`: delay of the message before being delivered
    - `VisibilityTimeout`: time in seconds where a message is only visible to a single consumer
+ `QueueUrl`: the URL of the queue to send message to
+ `MessageBody`: should be a string
+ `MaxNumberOfMessages`: the maximum number of messages to retrieve

## SQS CLI

* List queues:
```
aws sqs list-queues \
    --endpoint-url=http://localhost:4566
```

* Create queue:
```
aws sqs create-queue --queue-name my-queue \
    --endpoint-url=http://localhost:4566
```

Output example:
```
{
    "QueueUrl": "http://localhost:4566/000000000000/my-queue"
}
```

* Send message:
```
aws sqs send-message --message-body 'Test SQS message' \
    --queue-url http://localhost:4566/000000000000/my-queue \
    --endpoint-url=http://localhost:4566
```
Output:
```
{
    "MD5OfMessageBody": "108461c64e49b6406c003ee0534d44f0",
    "MessageId": "8f827c22-4ae9-3873-45a0-197e3dafa8c3"
}
```

* Consume message:
```
aws sqs receive-message \
    --queue-url http://localhost:4566/000000000000/my-queue \
    --endpoint-url=http://localhost:4566
```
Output:
```
{
    "Messages": [
        {
            "MessageId": "8f827c22-4ae9-3873-45a0-197e3dafa8c3",
            "ReceiptHandle": "pyzwqmgxbastsmohifkjrroezeppvuqnxfqitdpqhjzjjoaxnlvxssklhhouwexxpjkbondwqfyafxecyfwdznibltkyrfgkpggskqqoylyjeeqlwtlmmyplbrzcfiqitavbvwiciqxbnwgsxosqdnstwycuxrmuqekniwzsuphrvmvplaydrujha",
            "MD5OfBody": "108461c64e49b6406c003ee0534d44f0",
            "Body": "Test SQS message"
        }
    ]
}
```

* Delete queue:
```
aws sqs delete-queue \
    --queue-url http://localhost:4566/000000000000/my-queue \
    --endpoint-url=http://localhost:4566
```

## Resources

https://hub.docker.com/r/localstack/localstack

https://docs.localstack.cloud/aws/sqs/

AWS SQS CLI
https://docs.aws.amazon.com/cli/latest/reference/sqs/

AWS SDK V2
https://aws.github.io/aws-sdk-go-v2/docs/code-examples/sqs/
package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"
	"log"
	"sync"
)

type SqsConsumer struct {
	QueueURL            string
	MaxNumberOfMessages int64
	WaitTimeSeconds     int64
	SleepDuration       time.Duration

	client *sqs.SQS
}

type ReceiveMessageChan chan *sqs.Message
type DeleteMessageChan chan *sqs.Message

func NewSqsConsumer(queueURL string) *SqsConsumer {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))

	return &SqsConsumer{
		QueueURL:            queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
		SleepDuration:       1 * time.Second,
		client:              sqs.New(sess),
	}
}

func (c *SqsConsumer) Consume(messages ReceiveMessageChan) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.QueueURL),
		AttributeNames:        []*string{aws.String("All")},
		MaxNumberOfMessages:   aws.Int64(c.MaxNumberOfMessages),
		MessageAttributeNames: []*string{aws.String("All")},
		WaitTimeSeconds:       aws.Int64(c.WaitTimeSeconds),
	}

	for {
		received, err := c.client.ReceiveMessage(params)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Printf("Error reading queue: %s", awsErr)
				time.Sleep(c.SleepDuration)
			}
		} else {
			if len(received.Messages) == 0 {
				time.Sleep(c.SleepDuration)
			} else {
				for _, m := range received.Messages {
					messages <- m
				}
			}
		}
	}
}

func (c *SqsConsumer) Delete(messages DeleteMessageChan, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for m := range messages {
		_, err := c.client.DeleteMessage(
			&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(c.QueueURL),
				ReceiptHandle: aws.String(*m.ReceiptHandle),
			},
		)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Printf("Error deleting message: %s", awsErr)
			}
		} else {
			log.Printf("Delete message: %s", *m.MessageId)
		}
	}
}

package main

import (
	"github.com/shin1x1/go-sqsd/config"
	"github.com/shin1x1/go-sqsd/consumer"
	"github.com/shin1x1/go-sqsd/task"
	"github.com/shin1x1/go-sqsd/worker"
	"sync"
)

func main() {
	cfg, err := config.LoadConfig("./go-sqsd.conf")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	receiveMessages := make(consumer.ReceiveMessageChan)
	deleteMessages := make(consumer.DeleteMessageChan)

	for i := 1; i <= cfg.Worker.Workers; i++ {
		t := &task.HttpTask{
			No:  i,
			Url: cfg.Worker.Url,
		}
		go (&worker.Worker{}).Run(receiveMessages, deleteMessages, t, &wg)
	}

	c := consumer.NewSqsConsumer(cfg.Sqs.QueueUrl)
	go c.Delete(deleteMessages, &wg)
	c.Consume(receiveMessages)

	close(receiveMessages)

	wg.Wait()
}

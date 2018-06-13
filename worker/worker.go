package worker

import (
	"github.com/shin1x1/go-sqsd/task"
	"sync"
	"github.com/shin1x1/go-sqsd/consumer"
)

type Worker struct {
}

func (w *Worker) Run(c consumer.ReceiveMessageChan, d consumer.DeleteMessageChan, t task.Task, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for m := range c {
		if err := t.Run(m); err == nil {
			d <- m
		}
	}
}

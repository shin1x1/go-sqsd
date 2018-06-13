package task

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Task interface {
	Run(m *sqs.Message) error
}

type PrintTask struct {
	No int
}

func (t *PrintTask) Run(m *sqs.Message) error {
	fmt.Printf("%d: %s\n", t.No, m)

	return nil
}

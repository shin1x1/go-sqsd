package task

import (
	"fmt"
	"net/http"
	"bytes"
	"log"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"
	"strconv"
)

type HttpTask struct {
	No          int
	Url         string
	ContentType string
}

func (t *HttpTask) Run(m *sqs.Message) error {
	b := bytes.NewBuffer([]byte(*m.Body))

	req, _ := http.NewRequest("POST", t.Url, b)

	setHeaders(req, t, m)

	fmt.Printf("HTTP: %d: %s %s\n", t.No, t.Url, *m.Body)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error delivering message: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		return fmt.Errorf("Error: %s", m.MessageId)
	}

	return nil
}

func setHeaders(req *http.Request, t *HttpTask, m *sqs.Message) {
	req.Header.Set("Content-Type", t.ContentType)
	req.Header.Add("User-Agent", "go-sqsd")
	req.Header.Add("X-Aws-Sqsd-Msgid", *m.MessageId)

	if v, ok := m.Attributes["ApproximateFirstReceiveTimestamp"]; ok {
		if v1, err := convertUnixTimeToRFC3339(v); err == nil {
			req.Header.Add("X-Aws-Sqsd-First-Received-At", v1)
		}
	}

	if v, ok := m.Attributes["ApproximateReceiveCount"]; ok {
		req.Header.Add("X-Aws-Sqsd-Receive-Count", *v)
	}

	if v, ok := m.Attributes["SentTimestamp"]; ok {
		if v1, err := convertUnixTimeToRFC3339(v); err == nil {
			req.Header.Add("X-Aws-Sqsd-Sent-Timestamp", v1)
		}
	}
}

func convertUnixTimeToRFC3339(s * string) (string, error) {
	v, err := strconv.Atoi(*s);
	if err != nil {
		return "", err
	}

	v = v / 1000
	return time.Unix(int64(v), 0).UTC().Format(time.RFC3339), nil
}

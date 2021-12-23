package healthcheck

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type SlackIncomingWebhookConfigJSON struct {
	URL string
}

// https://api.slack.com/messaging/webhooks
// https://golangcode.com/send-slack-messages-without-a-library/
func (s *SlackIncomingWebhookConfigJSON) fire(text string) error {
	body, err := json.Marshal(map[string]interface{}{
		"text": text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.URL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from slack")
	}

	return nil
}

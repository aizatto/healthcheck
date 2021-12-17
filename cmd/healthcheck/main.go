package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aizatto/healthcheck/clients"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func main() {
	clients.Init()
	errs := healthcheck()
	for _, err := range errs {
		log.Println(err)
	}
}

func healthcheck() []error {
	body, err := ioutil.ReadFile(os.Getenv("CONFIG_JSON5"))
	if err != nil {
		return []error{err}
	}

	var config ConfigJSON
	err = json5.Unmarshal(body, &config)
	if err != nil {
		return []error{err}
	}

	targets := make([]*Target, len(config.Targets))
	for i := range targets {
		targets[i] = config.Targets[i].toTarget()
	}

	for true {
		for _, target := range targets {
			err := target.healthcheck()
			if err != nil {
				log.Printf("%+v", target)
				if target.Online {
					// send alert service is offline
					config.offline(target, err)
					target.Online = false
				}
				log.Println(err)
				continue
			}

			if !target.Online {
				// send alert back online
				config.online(target)
				target.Online = true
			}
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (c ConfigJSON) offline(t *Target, targeterr error) {
	for _, alert := range c.Alerts {
		var err error
		switch alert.Type {
		case "slack-incoming-webhook":
			err = alert.SlackIncomingWebhookConfig.fire(
				fmt.Sprintf(
					"%s: %s is offline. %s",
					c.Name,
					t.Name,
					targeterr,
				),
			)
		}

		if err != nil {
			log.Println(err)
			// TODO: consider disabling alert?
		}
	}
}

func (c ConfigJSON) online(t *Target) {
	for _, alert := range c.Alerts {
		var err error
		switch alert.Type {
		case "slack-incoming-webhook":
			err = alert.SlackIncomingWebhookConfig.fire(
				fmt.Sprintf("%s: %s is online", c.Name, t.Name),
			)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

// https://api.slack.com/messaging/webhooks
// https://golangcode.com/send-slack-messages-without-a-library/
func (s SlackIncomingWebhookConfigJSON) fire(text string) error {
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}

	return nil
}

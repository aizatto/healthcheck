package healthcheck

import (
	"fmt"
	"log"
)

type ConfigInterface interface {
	offline(t TargetInterface, err error)
	online(t TargetInterface)
}

type ConfigJSON struct {
	Name    string
	Alerts  []AlertJSON
	Targets []TargetJSON
}

func (c ConfigJSON) offline(t TargetInterface, targeterr error) {
	for _, alert := range c.Alerts {
		message := fmt.Sprintf(
			"%s: %s is offline. %s",
			c.Name,
			t.name(),
			targeterr,
		)

		switch alert.Type {
		case "slack-incoming-webhook":
			err := alert.SlackIncomingWebhookConfig.fire(message)

			if err != nil {
				log.Println(err)
			}
		case "stdout":
			log.Println(message)
		default:
			log.Printf("unsupported alert type: %s", alert.Type)
		}
	}
}

func (c ConfigJSON) online(t TargetInterface) {
	for _, alert := range c.Alerts {
		message := fmt.Sprintf("%s: %s is online", c.Name, t.name())

		switch alert.Type {
		case "slack-incoming-webhook":
			err := alert.SlackIncomingWebhookConfig.fire(message)

			if err != nil {
				log.Println(err)
			}
		case "stdout":
			log.Println(message)
		default:
			log.Printf("unsupported alert type: %s", alert.Type)
		}
	}
}

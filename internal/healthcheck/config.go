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
		var err error
		switch alert.Type {
		case "slack-incoming-webhook":
			err = alert.SlackIncomingWebhookConfig.fire(
				fmt.Sprintf(
					"%s: %s is offline. %s",
					c.Name,
					t.name(),
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

func (c ConfigJSON) online(t TargetInterface) {
	for _, alert := range c.Alerts {
		var err error
		switch alert.Type {
		case "slack-incoming-webhook":
			err = alert.SlackIncomingWebhookConfig.fire(
				fmt.Sprintf("%s: %s is online", c.Name, t.name()),
			)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

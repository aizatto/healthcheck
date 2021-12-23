package healthcheck

import (
	"os"

	"github.com/rollbar/rollbar-go"
)

type RollbarConfigJSON struct {
	Environment RollbarConfigValue `json:"environment"`
	Token       RollbarConfigValue `json:"token"`
	Level       RollbarConfigValue `json:"level"`
}

type RollbarConfigValue struct {
	Value     string `json:"value"`
	ValueFrom string `json:"valueFrom"`
	Key       string `json:"key"`
}

// https://api.slack.com/messaging/webhooks
// https://golangcode.com/send-slack-messages-without-a-library/
func (r *RollbarConfigJSON) fire(text string) error {
	token := r.Token.value()
	if token == "" {
		return nil
	}

	rollbar.SetEnabled(true)
	rollbar.SetToken(token)

	if environment := r.Token.value(); environment != "" {
		rollbar.SetEnvironment(environment)
	}

	level := r.Level.value()
	if level == "" {
		level = rollbar.ERR
	}

	rollbar.Log(level, text)
	return nil
}

func (c *RollbarConfigValue) value() string {
	switch c.ValueFrom {
	case "env":
		return os.Getenv(c.Key)
	default:
		return c.Value
	}
}

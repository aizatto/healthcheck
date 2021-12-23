package healthcheck

type AlertJSON struct {
	Type                       string
	RollbarConfig              *RollbarConfigJSON
	SlackIncomingWebhookConfig *SlackIncomingWebhookConfigJSON
}

type AlertInterface interface {
	fire(string) error
}

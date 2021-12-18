package healthcheck

type AlertJSON struct {
	Type                       string
	SlackIncomingWebhookConfig *SlackIncomingWebhookConfigJSON
}

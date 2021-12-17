package healthcheck

type AlertJSON struct {
	Type                       string
	SlackIncomingWebhookConfig *SlackIncomingWebhookConfigJSON
}

type TargetJSON struct {
	Alerts    []AlertJSON
	ValueFrom string
	Key       string
	Name      string
	URL       string
}

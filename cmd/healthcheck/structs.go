package main

type Target struct {
	Name   string
	URL    string
	Online bool
}

type AlertJSON struct {
	Type                       string
	SlackIncomingWebhookConfig SlackIncomingWebhookConfigJSON
}

type SlackIncomingWebhookConfigJSON struct {
	URL string
}

type TargetJSON struct {
	Alerts    []AlertJSON
	ValueFrom string
	Key       string
	Name      string
	URL       string
}

type ConfigJSON struct {
	Name    string
	Alerts  []AlertJSON
	Targets []TargetJSON
}

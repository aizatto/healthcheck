package main

type Target struct {
	Name   string
	URL    string
	Online bool
}

type AlertJSON struct {
	Type string
}

type TargetJSON struct {
	Alerts    []AlertJSON
	ValueFrom string
	Key       string
	Name      string
	URL       string
}

type ConfigJSON struct {
	Alerts  []AlertJSON
	Targets []TargetJSON
}

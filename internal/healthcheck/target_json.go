package healthcheck

import (
	"net/http"
	"os"
)

type TargetJSON struct {
	Alerts            []AlertJSON        `json:"alerts"`
	ValueFrom         string             `json:"valueFrom"`
	Key               string             `json:"key"`
	Name              string             `json:"name"`
	URL               string             `json:"url"`
	HttpRequestConfig *HttpRequestConfig `json:"httpRequestConfig"`
}

type HttpRequestConfig struct {
	Method               string `json:"method"`
	ContentType          string `json:"contentType"`
	Body                 string `json:"body"`
	ExpectedResponseCode int    `json:"expectedResponseCode"`
}

func (t *TargetJSON) ToTarget() *Target {
	return &Target{
		Config: t,
		Online: true, // asume everything is online by default
	}
}

func (t *TargetJSON) name() string {
	if len(t.Name) > 0 {
		return t.Name
	}

	switch t.ValueFrom {
	case "env":
		return t.Key
	}

	return t.URL
}

func (t *TargetJSON) url() string {
	url := ""
	switch t.ValueFrom {
	case "env":
		url = os.Getenv(t.Key)
	default:
		url = t.URL
	}

	return url
}

func (t *TargetJSON) httpRequestConfig() *HttpRequestConfig {
	if t.HttpRequestConfig == nil {
		return &HttpRequestConfig{
			Method:               http.MethodGet,
			ExpectedResponseCode: http.StatusOK,
		}
	}

	if t.HttpRequestConfig.ExpectedResponseCode == 0 {
		t.HttpRequestConfig.ExpectedResponseCode = http.StatusOK
	}

	return t.HttpRequestConfig
}

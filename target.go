package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func (t TargetJSON) toTarget() *Target {
	return &Target{
		Name:   t.name(),
		URL:    t.url(),
		Online: false,
	}
}

func (t TargetJSON) name() string {
	if len(t.Name) > 0 {
		return t.Name
	}

	switch t.ValueFrom {
	case "env":
		return t.Key
	}

	return t.URL
}

func (t TargetJSON) url() string {
	url := ""
	switch t.ValueFrom {
	case "env":
		url = os.Getenv(t.Key)
	default:
		url = t.URL
	}

	return url
}

func (t *Target) healthcheck() error {
	resp, err := http.Get(t.URL)
	if err != nil {
		return errors.Wrapf(err, "%s Failed: %s", t.Name, t.URL)
	}

	log.Printf("%s (%d): %s\n", t.Name, resp.StatusCode, t.URL)

	return nil
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aizatto/healthcheck/clients"
	"github.com/pkg/errors"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Target struct {
	ValueFrom string
	Key       string
	Name      string
	URL       string
}

type Config struct {
	Targets []Target
}

func main() {
	clients.Init()
	errs := healthcheck()
	for _, err := range errs {
		log.Println(err)
	}
}

func healthcheck() []error {
	body, err := ioutil.ReadFile(os.Getenv("CONFIG_JSON5"))
	if err != nil {
		return []error{err}
	}

	var config Config
	err = json5.Unmarshal(body, &config)
	if err != nil {
		return []error{err}
	}

	for true {
		for _, target := range config.Targets {
			err := target.healthcheck()
			if err != nil {
				log.Println(err)
				continue
			}
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (t Target) name() string {
	if len(t.Name) > 0 {
		return t.Name
	}

	switch t.ValueFrom {
	case "env":
		return t.Key
	}

	return ""
}

func (t Target) url() string {
	url := ""
	switch t.ValueFrom {
	case "env":
		url = os.Getenv(t.Key)
	default:
		url = t.URL
	}

	return url
}

func (t Target) healthcheck() error {
	url := t.url()
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "Failed: %s", url)
	}

	log.Printf("%s (%d): %s\n", t.Key, resp.StatusCode, url)

	return nil
}

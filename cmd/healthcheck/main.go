package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aizatto/healthcheck/clients"
	"github.com/aizatto/healthcheck/internal/healthcheck"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func main() {
	clients.Init()
	errs := run()
	for _, err := range errs {
		log.Println(err)
	}
}

func run() []error {
	body, err := ioutil.ReadFile(os.Getenv("CONFIG_JSON5"))
	if err != nil {
		return []error{err}
	}

	var config healthcheck.ConfigJSON
	if err := json5.Unmarshal(body, &config); err != nil {
		return []error{err}
	}

	targets := make([]healthcheck.TargetInterface, len(config.Targets))
	for i := range targets {
		targets[i] = config.Targets[i].ToTarget()
	}

	healthcheck.HealthcheckTargets(&config, targets)
	return []error{}
}

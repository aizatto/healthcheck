package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aizatto/healthcheck/clients"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

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

	var config ConfigJSON
	err = json5.Unmarshal(body, &config)
	if err != nil {
		return []error{err}
	}

	for true {
		for _, target := range config.Targets {
			err := target.toTarget().healthcheck()
			if err != nil {
				log.Println(err)
				continue
			}
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

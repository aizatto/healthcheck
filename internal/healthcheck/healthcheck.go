package healthcheck

import (
	"log"
	"time"
)

func HealthcheckTargets(config ConfigInterface, targets []TargetInterface) {
	for {
		for _, target := range targets {
			log.Printf("%+v", target)
			HealthcheckTarget(config, target)
		}

		time.Sleep(5 * time.Second)
	}
}

func HealthcheckTarget(config ConfigInterface, target TargetInterface) error {
	if err := target.healthcheck(); err != nil {
		if target.isOnline() {
			// send alert service is offline
			config.offline(target, err)
			target.setOnline(false)
		}

		return err
	}

	// online := target.isOnline()
	if !target.isOnline() {
		// send alert back online
		config.online(target)
		target.setOnline(true)
	}
	return nil
}

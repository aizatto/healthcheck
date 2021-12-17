package healthcheck

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type TargetInterface interface {
	name() string
	healthcheck() error
	isOnline() bool
	setOnline(bool) TargetInterface
}

type Target struct {
	Name   string
	URL    string
	Online bool
}

func (t *Target) name() string {
	return t.Name
}

func (t *Target) healthcheck() error {
	resp, err := http.Get(t.URL)
	if err != nil {
		return errors.Wrapf(err, "%s Failed: %s", t.Name, t.URL)
	}

	log.Printf("%s (%d): %s\n", t.Name, resp.StatusCode, t.URL)

	return nil
}

func (t *Target) isOnline() bool {
	return t.Online
}

func (t *Target) setOnline(status bool) TargetInterface {
	t.Online = status
	return t
}

package healthcheck

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type TargetInterface interface {
	name() string
	healthcheck() error
	isOnline() bool
	setOnline(bool) TargetInterface
}

type Target struct {
	Config *TargetJSON
	Online bool
}

func (t *Target) name() string {
	return t.Config.name()
}

func (t *Target) healthcheck() error {
	var resp *http.Response
	var err error

	requestConfig := t.Config.httpRequestConfig()

	switch requestConfig.Method {
	case http.MethodPost:
		resp, err = http.Post(
			t.Config.url(),
			requestConfig.ContentType,
			strings.NewReader(requestConfig.Body),
		)
	default:
		resp, err = http.Get(t.Config.url())
	}

	if err != nil {
		return errors.Wrapf(err, "%s Failed: %s", t.name(), t.Config.URL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != t.Config.HttpRequestConfig.ExpectedResponseCode {
		return fmt.Errorf(
			"%s Failed: %s.\nUnexpected HTTP response code %d, expecting %d.",
			t.name(),
			t.Config.URL,
			resp.StatusCode,
			t.Config.HttpRequestConfig.ExpectedResponseCode,
		)
	}

	log.Printf("%s (%d): %s\n", t.name(), resp.StatusCode, t.Config.url())

	return nil
}

func (t *Target) isOnline() bool {
	return t.Online
}

func (t *Target) setOnline(status bool) TargetInterface {
	t.Online = status
	return t
}

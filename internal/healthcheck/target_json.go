package healthcheck

import "os"

func (t TargetJSON) ToTarget() *Target {
	return &Target{
		Name:   t.name(),
		URL:    t.url(),
		Online: true, // asume everything is online by default
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

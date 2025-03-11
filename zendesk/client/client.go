package client

import (
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/nukosuke/go-zendesk/zendesk"
)

type (
	Client struct {
		// use struct embedding for extension
		zendesk.Client
	}
)

// addOptions build query string
func addOptions(s string, opts interface{}) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

package webmention

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type Webmention struct {
	Src  string
	Dest string
}

const (
	InvalidWebmentionEndpointURLErrorPrefix = "Webmention URL cannot be parsed"
	InvalidSourceEndpointURLErrorPrefix = "Source URL cannot be parsed"
	InvalidDestEndpointURLErrorPrefix = "Destination URL cannot be parsed"
	HTTPRequestConstructionError = "Failed to construct HTTP Request. Please report this as an issue."
)

func (w *Webmention) ToHTTPRequest (webmentionEndpointURL string) (*http.Request, error) {
	notify, err := url.Parse(webmentionEndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, InvalidWebmentionEndpointURLErrorPrefix)
	}

	source, err := url.Parse(w.Src)
	if err != nil {
		return nil, errors.Wrap(err, InvalidSourceEndpointURLErrorPrefix)
	}

	target, err := url.Parse(w.Dest)
	if err != nil {
		return nil, errors.Wrap(err, InvalidDestEndpointURLErrorPrefix)
	}

	form := url.Values{}
	form.Add("source", source.String())
	form.Add("target", target.String())

	req, err := http.NewRequest(http.MethodPost, notify.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, HTTPRequestConstructionError)
	}

	req.Header = map[string][]string{ "Content-Type": []string{contentTypeFormURLEncoded} }
	req.Form = form

	return req, nil
}

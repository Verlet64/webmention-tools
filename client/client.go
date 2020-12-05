package client

import (
	"errors"
	"fmt"
	"net/http"
	"webmention-tools/parse"
)

type WebmentionClient struct {
	client *http.Client
}

const (
	WebmentionClientFailedFetchDiscoveryDocument = "Failed to fetch discovery document"
)

func NewWebmentionClient() *WebmentionClient {
	defaultHttpClient := http.DefaultClient

	return &WebmentionClient{
		client: defaultHttpClient,
	}
}

func (wm *WebmentionClient) DiscoverWebmentionEndpointFromURL(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := wm.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf(WebmentionClientFailedFetchDiscoveryDocument+"[Status %v]", res.StatusCode))
	}

	endpoint, err := parse.ParseWebmentionURLFromHTML(res.Body)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}

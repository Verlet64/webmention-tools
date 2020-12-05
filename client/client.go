package client

import (
	"errors"
	"fmt"
	"net/http"
	"webmention-tools/parse"
	"webmention-tools/webmention"
)

type WebmentionClient struct {
	client *http.Client
}

const (
	WebmentionClientFailedFetchDiscoveryDocument = "Failed to fetch parse document"
	WebmentionSendFailure                        = "Failed to dispatch webmention"
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

func (wm *WebmentionClient) SendWebmention(url string, source string, target string) error {
	mention := webmention.Webmention{Dest: target, Src: source}
	req, err := mention.ToHTTPRequest(url)
	if err != nil {
		return err
	}

	res, err := wm.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("%s [Status=%v]", WebmentionSendFailure, res.StatusCode))
	}

	return nil
}

package webmention_test

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"webmention-tools/webmention"
)

const (
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

func TestWebmention_ToHTTPRequestReturnsValidWebmentionRequest(t *testing.T) {
	webmentionEndpoint, err := url.Parse("https://example.com/webmention")
	if err != nil {
		t.Fatalf("Failed to construct test webmention URL")
	}

	wr := webmention.Webmention{Dest: "https://example.com/foobar", Src: "https://example.com/testbar"}

	data := url.Values{}
	data.Add("source", wr.Src)
	data.Add("target", wr.Dest)

	expected := &http.Request{
		Method: http.MethodPost,
		Header: map[string][]string{
			"Content-Type": []string{contentTypeFormURLEncoded},
		},
		Host: "example.com",
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL: webmentionEndpoint,
		Form: data,
	}
	expected = expected.WithContext(context.Background())

	got, err := wr.ToHTTPRequest(webmentionEndpoint.String())

	if err != nil {
		t.Fatalf("Expected no errors, found %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected %v to exactly match %v", expected, got)
	}
}

func TestWebmention_ToHTTPRequestFailsWithNonHTTPWebmentionEndpoint(t *testing.T) {
	wr := webmention.Webmention{Dest: "https://example.com/foobar", Src: "https://example.com/testbar"}

	got, gotErr := wr.ToHTTPRequest(":example.com/webmention")
	if got != nil {
		t.Fatalf("Error should not return any non-nil zero value, found %v", got)
	}

	if gotErr == nil {
		t.Fatalf("Expected an error, found nil")
	}

	errPrefix := strings.Split(gotErr.Error(), ":")[0]
	if errPrefix != webmention.InvalidWebmentionEndpointURLErrorPrefix {
		t.Fatalf("Expected %v, got %v", webmention.InvalidWebmentionEndpointURLErrorPrefix, errPrefix)
	}
}

func TestWebmention_ToHTTPRequestFailsWithNonHTTPSourceURL(t *testing.T) {
	webmentionEndpoint, err := url.Parse("https://example.com/webmention")
	if err != nil {
		t.Fatalf("Failed to construct test webmention URL")
	}

	wr := webmention.Webmention{Dest: "https://example.com/foobar", Src: ":example.com/testbar"}

	got, gotErr := wr.ToHTTPRequest(webmentionEndpoint.String())
	if got != nil {
		t.Fatalf("Error should not return any non-nil zero value, found %v", got)
	}

	if gotErr == nil {
		t.Fatalf("Expected an error, found nil")
	}

	errPrefix := strings.Split(gotErr.Error(), ":")[0]
	if errPrefix != webmention.InvalidSourceEndpointURLErrorPrefix {
		t.Fatalf("Expected %v, got %v", webmention.InvalidSourceEndpointURLErrorPrefix, errPrefix)
	}
}

func TestWebmention_ToHTTPRequestFailsWithNonHTTPDestURL(t *testing.T) {
	webmentionEndpoint, err := url.Parse("https://example.com/webmention")
	if err != nil {
		t.Fatalf("Failed to construct test webmention URL")
	}

	wr := webmention.Webmention{Dest: ":example.com/foobar", Src: "https://example.com/testbar"}

	got, gotErr := wr.ToHTTPRequest(webmentionEndpoint.String())
	if got != nil {
		t.Fatalf("Error should not return any non-nil zero value, found %v", got)
	}

	if gotErr == nil {
		t.Fatalf("Expected an error, found nil")
	}

	errPrefix := strings.Split(gotErr.Error(), ":")[0]
	if errPrefix != webmention.InvalidDestEndpointURLErrorPrefix {
		t.Fatalf("Expected %v, got %v", webmention.InvalidDestEndpointURLErrorPrefix, errPrefix)
	}
}


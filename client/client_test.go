package client_test

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"webmention-tools/client"
)

const (
	testTemplate = `
<!doctype html>
<html>
	<head>
		{{ if .Endpoint }} 
			<link href="{{ .Endpoint }}" rel="webmention" />
		{{ end }}
	</head>
	<body>
    	<span> Some Content </span>
	</body>
</html>
`
)

func TestWebmentionClient_SendWebmentionSendsValidWebmentionPayload(t *testing.T) {
	source := "http://test.example.com/src"
	dest := "http://test.example.com/dest"

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("Expected a POST request, got %v", request.Method)
		}

		if request.Proto != "HTTP/1.1" {
			t.Fatalf("Expected a HTTP/1.1 request")
		}

		if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Fatalf("Unexpected content type %v", request.Header.Get("Content-Type"))
		}

		expectedPayload := url.Values{}
		expectedPayload.Add("source", source)
		expectedPayload.Add("target", dest)

		writer.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	wm := client.NewWebmentionClient()
	err := wm.SendWebmention(srv.URL, source, dest)
	if err != nil {
		t.Fatalf("Encountered error sending webmention: %v", err)
	}
}

func TestWebmentionClient_SendWebmentionErrorsOnNon2xx(t *testing.T) {
	source := "http://test.example.com/src"
	dest := "http://test.example.com/dest"

	var status int

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(status)
	}))
	defer srv.Close()

	testCases := []struct {
		responseCode int
		errorPrefix  string
	}{
		{http.StatusPermanentRedirect, client.WebmentionSendFailure},
		{http.StatusInternalServerError, client.WebmentionSendFailure},
		{http.StatusForbidden, client.WebmentionSendFailure},
	}

	for _, tcs := range testCases {
		status = tcs.responseCode
		wm := client.NewWebmentionClient()
		err := wm.SendWebmention(srv.URL, source, dest)
		if err == nil {
			t.Fatalf("Expected an error from SendWebmention")
		}

		if err.Error() != fmt.Sprintf("%v [Status=%v]", client.WebmentionSendFailure, status) {
			t.Fatalf("Encountered error sending webmention: %v", err)
		}
	}

}

func TestWebmentionClient_DiscoverWebmentionEndpointFromHTMLReturnsWebmentionURL(t *testing.T) {
	var html []byte

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write(html)
	}))
	defer srv.Close()

	tests := []struct {
		Endpoint          string
		ShowWebmentionURL bool
		ParsedURL         string
	}{
		{Endpoint: srv.URL + "/webmention", ParsedURL: fmt.Sprintf("%v/webmention", srv.URL)},
		{Endpoint: "", ParsedURL: ""},
	}

	temp, err := template.New("client-test").Parse(testTemplate)
	if err != nil {
		t.Fatalf("Failed to construct client response template")
	}

	for idx, tc := range tests {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)

		err = temp.Execute(w, tc)
		err = w.Flush()
		if err != nil {
			t.Fatalf("Failed to insert template values with error: %v", err)
		}

		html = b.Bytes()

		c := client.NewWebmentionClient()
		parsed, err := c.DiscoverWebmentionEndpointFromURL(fmt.Sprintf("%v", srv.URL))
		if err != nil {
			t.Fatalf("Encountered error when discovering Webmention URL for provided document URL: %v", err)
		}

		if parsed != tc.ParsedURL {
			t.Fatalf("[%v] Expected Webmention URL %v, found %v", idx, tc.ParsedURL, parsed)
		}
	}
}

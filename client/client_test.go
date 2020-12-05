package client_test

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
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

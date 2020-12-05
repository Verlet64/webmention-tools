package discovery_test

import (
	"bufio"
	"bytes"
	"html/template"
	"strings"
	"testing"
	"webmention-tools/discovery"
)

const (
	testTemplate = `
<!doctype html>
<html>
	<head>
		{{ if .Endpoint }} 
			<link href="{{ .Endpoint }}" rel="webmention" />
		{{ else }}
		{{ end }}
	</head>
	<body>
    	<span> Some Content </span>
	</body>
</html>
`
)

func TestParseWebmentionURLFromHTML(t *testing.T) {
	tests := []struct {
		Endpoint          string
		ShowWebmentionURL bool
		ParsedURL         string
		ParseError        error
	}{
		{Endpoint: "http://example.com/webmention", ParsedURL: "http://example.com/webmention", ParseError: nil},
		{Endpoint: "", ParsedURL: "", ParseError: nil},
		{Endpoint: ":example.com/webmention", ParsedURL: "", ParseError: nil},
	}

	temp, err := template.New("parse-webmention-test-template").Parse(testTemplate)
	if err != nil {
		t.Fatalf("Failed to construct test template with error: %v", err)
	}

	for idx, testCase := range tests {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)

		err = temp.Execute(w, testCase)
		err = w.Flush()
		if err != nil {
			t.Fatalf("Failed to insert template values with error: %v", err)
		}

		html := string(b.Bytes())

		got, err := discovery.ParseWebmentionURLFromHTML(&b)
		if got != testCase.ParsedURL || err != testCase.ParseError {
			t.Fatalf(`
Failed Scenario %v
Expected result %v, error %v, got result %v, error %v
Template: %v
`, idx, testCase.ParsedURL, testCase.ParseError, got, err, html)
		}

	}

}

func TestParseWebmentionURLFromHTML_InvalidHTMLContent(t *testing.T) {
	var b bytes.Buffer
	_, err := b.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Failed to construct test input")
	}

	got, err := discovery.ParseWebmentionURLFromHTML(&b)
	if got != "" {
		t.Fatalf("Parsed URL from invalid HTML")
	}

	if err != nil && !strings.HasPrefix(err.Error(), discovery.HTMLParseFailurePrefix) {
		t.Fatalf("Expected error to start with %v, found %v", discovery.HTMLParseFailurePrefix, err)
	}
}

package parse

import (
	"io"
	"net/url"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

const (
	HTMLParseFailurePrefix    = "Failed to parse HTML"
	WebmentionURLParseFailure = "Failed to extract URL for webmention endpoint"
)

func ParseWebmentionURLFromHTML(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", errors.Wrap(err, HTMLParseFailurePrefix)
	}

	var endpoint *url.URL

	doc.Find("link[rel='webmention']").Last().Each(func(i int, s *goquery.Selection) {
		raw, exists := s.Attr("href")
		if !exists {
			err = errors.New(WebmentionURLParseFailure)
		}

		endpoint, err = url.Parse(raw)
		if err != nil || endpoint.Scheme != "http" {
			endpoint = nil
		}
	})

	if endpoint != nil {
		return endpoint.String(), err
	}

	return "", err
}

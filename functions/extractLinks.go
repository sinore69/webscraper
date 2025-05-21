package functions

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(baseURL string, body io.Reader) ([]string, error) {
	links := []string{}
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return links, nil
			}
			return nil, z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						link := resolveURL(baseURL, attr.Val)
						if link != "" && strings.HasPrefix(link, "http") {
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

func resolveURL(base, href string) string {
	baseParsed, err := url.Parse(base)
	if err != nil {
		return ""
	}
	hrefParsed, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return baseParsed.ResolveReference(hrefParsed).String()
}

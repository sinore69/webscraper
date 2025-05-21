package functions

import (
	"fmt"
	"io"
	"net/http"
)

func FetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", pageURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status for %s: %s", pageURL, resp.Status)
	}
	return resp.Body, nil
}

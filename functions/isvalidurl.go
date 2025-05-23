package functions

import "net/url"

func IsValidURL(s string) bool {
	parsed, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return parsed.Scheme != "" && parsed.Host != ""
}

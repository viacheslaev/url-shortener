package link

import "net/url"

func validateURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil &&
		(u.Scheme == "http" || u.Scheme == "https") &&
		u.Host != ""
}

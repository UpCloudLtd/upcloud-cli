package format

import (
	"net/url"
	"strings"
)

// HiddenPlaceholder replaces a secret value in human-readable output.
const HiddenPlaceholder = "hidden (--show-secrets)"

// Secret returns value when show is true. Otherwise it returns HiddenPlaceholder,
// or an empty string when value is empty.
func Secret(value string, show bool) string {
	if show || value == "" {
		return value
	}
	return HiddenPlaceholder
}

// SecretURI returns raw with the password of its user info replaced by
// "xxxxx", unless show is true. The whole value is hidden when it does not
// parse as a URL, or when a known secret is still in the redacted value (for
// example a password in the query string or in an opaque URI).
func SecretURI(raw string, show bool, knownSecrets ...string) string {
	if show || raw == "" {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return HiddenPlaceholder
	}
	redacted := u.Redacted()
	for _, secret := range knownSecrets {
		if secret == "" {
			continue
		}
		if strings.Contains(redacted, secret) ||
			strings.Contains(redacted, url.QueryEscape(secret)) ||
			strings.Contains(redacted, url.PathEscape(secret)) {
			return HiddenPlaceholder
		}
	}
	return redacted
}

package format

import "net/url"

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
// "xxxxx", unless show is true. A value that does not parse as a URL is
// hidden completely, because it can still contain a password.
func SecretURI(raw string, show bool) string {
	if show || raw == "" {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return HiddenPlaceholder
	}
	return u.Redacted()
}

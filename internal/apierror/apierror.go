// Package apierror provides the CLI's own error type for UpCloud API errors, independent of the SDK version.
package apierror

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Problem is an UpCloud API error. Its Error and ErrorCode output match the v8 SDK's upcloud.Problem.
type Problem struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
	CorrelationID string         `json:"correlation_id,omitempty"`
	Status        int            `json:"status"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func (p *Problem) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "%s (type=%s, status=%d", p.Title, p.Type, p.Status)
	if p.CorrelationID != "" {
		_, _ = fmt.Fprintf(&sb, ", correlation_id=%s", p.CorrelationID)
	}
	for _, ip := range p.InvalidParams {
		_, _ = fmt.Fprintf(&sb, ", invalid_params_%s='%s'", ip.Name, ip.Reason)
	}
	sb.WriteString(")")
	return sb.String()
}

// ErrorCode returns the error code, e.g. SERVICE_NOT_FOUND, also when Type is a documentation URL.
func (p *Problem) ErrorCode() string {
	parsedURL, err := url.Parse(p.Type)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return p.Type
	}
	return strings.Replace(parsedURL.Fragment, "ERROR_", "", 1)
}

// FromResponse converts an unsuccessful API response into a *Problem. It accepts problem+json and legacy error bodies.
func FromResponse(status int, body []byte) error {
	if status >= 200 && status < 300 {
		return fmt.Errorf("unexpected API response (status %d)", status)
	}

	prob := &Problem{}
	if json.Unmarshal(body, prob) == nil && prob.Title != "" {
		if prob.Status == 0 {
			prob.Status = status
		}
		return prob
	}

	legacy := struct {
		Error struct {
			ErrorCode    string `json:"error_code"`
			ErrorMessage string `json:"error_message"`
		} `json:"error"`
	}{}
	if json.Unmarshal(body, &legacy) == nil && legacy.Error.ErrorCode != "" {
		return &Problem{Type: legacy.Error.ErrorCode, Title: legacy.Error.ErrorMessage, Status: status}
	}

	title := strings.TrimSpace(string(body))
	if title == "" {
		title = http.StatusText(status)
	}
	return &Problem{Title: title, Status: status}
}

func IsNotFound(err error) bool {
	var prob *Problem
	return errors.As(err, &prob) && prob.Status == http.StatusNotFound
}

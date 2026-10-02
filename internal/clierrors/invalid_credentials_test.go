package clierrors

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestCheckAuthenticationFailed(t *testing.T) {
	for _, test := range []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "v8 authentication failed", err: &upcloud.Problem{Type: "https://developers.upcloud.com/1.3/errors#ERROR_AUTHENTICATION_FAILED"}, expected: true},
		{name: "v8 invalid credentials", err: &upcloud.Problem{Type: "INVALID_CREDENTIALS"}, expected: true},
		{name: "v9 authentication failed", err: &apierror.Problem{Type: "https://developers.upcloud.com/1.3/errors#ERROR_AUTHENTICATION_FAILED"}, expected: true},
		{name: "wrapped v9 invalid credentials", err: fmt.Errorf("resolving: %w", &apierror.Problem{Type: "INVALID_CREDENTIALS"}), expected: true},
		{name: "v8 other error", err: &upcloud.Problem{Type: "SERVICE_NOT_FOUND"}},
		{name: "v9 other error", err: &apierror.Problem{Type: "SERVICE_NOT_FOUND"}},
		{name: "plain error", err: fmt.Errorf("authentication failed")},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, CheckAuthenticationFailed(test.err))
		})
	}
}

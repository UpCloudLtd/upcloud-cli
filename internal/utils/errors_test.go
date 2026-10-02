package utils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestIsNotFoundError(t *testing.T) {
	for _, test := range []struct {
		name     string
		err      error
		expected bool
	}{
		{name: "v8 not found", err: &upcloud.Problem{Status: http.StatusNotFound}, expected: true},
		{name: "wrapped v8 not found", err: fmt.Errorf("wrapped: %w", &upcloud.Problem{Status: http.StatusNotFound}), expected: true},
		{name: "v9 not found", err: &apierror.Problem{Status: http.StatusNotFound}, expected: true},
		{name: "wrapped v9 not found", err: fmt.Errorf("wrapped: %w", &apierror.Problem{Status: http.StatusNotFound}), expected: true},
		{name: "v8 other status", err: &upcloud.Problem{Status: http.StatusBadRequest}},
		{name: "v9 other status", err: &apierror.Problem{Status: http.StatusBadRequest}},
		{name: "plain error", err: fmt.Errorf("not found")},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsNotFoundError(test.err))
		})
	}
}

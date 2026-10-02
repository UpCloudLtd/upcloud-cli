package apierror

import (
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromResponse(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   int
		body     string
		expected Problem
	}{
		{
			name:     "problem details",
			status:   http.StatusBadRequest,
			body:     `{"type": "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST", "title": "Invalid request", "status": 400, "correlation_id": "abc", "invalid_params": [{"name": "plan", "reason": "unknown"}]}`,
			expected: Problem{Type: "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST", Title: "Invalid request", Status: 400, CorrelationID: "abc", InvalidParams: []InvalidParam{{Name: "plan", Reason: "unknown"}}},
		},
		{
			name:     "problem details without status",
			status:   http.StatusConflict,
			body:     `{"type": "conflict", "title": "Conflict"}`,
			expected: Problem{Type: "conflict", Title: "Conflict", Status: http.StatusConflict},
		},
		{
			name:     "legacy error",
			status:   http.StatusNotFound,
			body:     `{"error": {"error_code": "SERVICE_NOT_FOUND", "error_message": "Service not found"}}`,
			expected: Problem{Type: "SERVICE_NOT_FOUND", Title: "Service not found", Status: http.StatusNotFound},
		},
		{
			name:     "plain text",
			status:   http.StatusBadGateway,
			body:     "bad gateway\n",
			expected: Problem{Title: "bad gateway", Status: http.StatusBadGateway},
		},
		{
			name:     "empty body",
			status:   http.StatusServiceUnavailable,
			expected: Problem{Title: "Service Unavailable", Status: http.StatusServiceUnavailable},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := FromResponse(test.status, []byte(test.body))
			var problem *Problem
			require.ErrorAs(t, err, &problem)
			assert.Equal(t, test.expected, *problem)
			assert.Equal(t, test.status == http.StatusNotFound, IsNotFound(err))
		})
	}

	assert.EqualError(t, FromResponse(http.StatusOK, nil), "unexpected API response (status 200)")
}

// Users must see the same messages and error codes as with the v8 SDK error type.
func TestProblem_MatchesV8(t *testing.T) {
	for _, p := range []Problem{
		{Type: "https://developers.upcloud.com/1.3/errors#ERROR_AUTHENTICATION_FAILED", Title: "Authentication failed", Status: 401},
		{Type: "SERVICE_NOT_FOUND", Title: "Service not found", Status: 404, CorrelationID: "abc"},
		{Type: "invalid", Title: "Invalid", Status: 400, InvalidParams: []InvalidParam{{Name: "plan", Reason: "unknown"}, {Name: "zone", Reason: "missing"}}},
	} {
		v8 := upcloud.Problem{Type: p.Type, Title: p.Title, Status: p.Status, CorrelationID: p.CorrelationID}
		for _, ip := range p.InvalidParams {
			v8.InvalidParams = append(v8.InvalidParams, upcloud.ProblemInvalidParam{Name: ip.Name, Reason: ip.Reason})
		}
		assert.Equal(t, v8.Error(), p.Error())
		assert.Equal(t, v8.ErrorCode(), p.ErrorCode())
	}
}

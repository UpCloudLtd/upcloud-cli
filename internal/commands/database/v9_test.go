package database

import (
	"context"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForDatabaseState(t *testing.T) {
	id := uuid.MustParse("5c2f1f0e-7d1a-4f7e-9f0b-2a6c1d3e4b5a")
	path := "GET /1.3/database/" + id.String()
	serverError := fakeResponse{status: http.StatusBadGateway, contentType: "text/plain", body: "bad gateway"}
	rebuilding := fakeResponse{status: http.StatusOK, body: `{"state": "rebuilding"}`}
	running := fakeResponse{status: http.StatusOK, body: `{"state": "running"}`}

	for _, test := range []struct {
		name     string
		sequence []fakeResponse
		error    string
		notFound bool
		polls    int
	}{
		{
			name:     "returns when the state is reached",
			sequence: []fakeResponse{rebuilding, rebuilding, running},
			polls:    3,
		},
		{
			name:     "retries up to three consecutive server errors",
			sequence: []fakeResponse{serverError, serverError, serverError, running},
			polls:    4,
		},
		{
			name:     "fails on the fourth consecutive server error",
			sequence: []fakeResponse{serverError, serverError, serverError, serverError},
			error:    "bad gateway (type=, status=502)",
			polls:    4,
		},
		{
			name:     "fails when the database enters the error state",
			sequence: []fakeResponse{rebuilding, {status: http.StatusOK, body: `{"state": "error", "state_error": {"init": "failed", "backup": "missing"}}`}},
			error:    "managed database entered error state: missing (backup), failed (init)",
			polls:    2,
		},
		{
			name:     "fails when the database is not found",
			sequence: []fakeResponse{{status: http.StatusNotFound, contentType: "application/problem+json", body: `{"type": "SERVICE_NOT_FOUND", "title": "Service not found", "status": 404}`}},
			error:    "Service not found (type=SERVICE_NOT_FOUND, status=404)",
			notFound: true,
			polls:    1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newFakeAPI(t, map[string]fakeResponse{path: {sequence: test.sequence}})

			err := waitForDatabaseState(context.Background(), api.client, id, databaseStateRunning)

			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.error)
				assert.Equal(t, test.notFound, apierror.IsNotFound(err))
			}
			assert.Len(t, api.all(), test.polls)
		})
	}
}

func TestPropertySchemaAllowsString(t *testing.T) {
	assert.True(t, propertySchema{Type: "string"}.allowsString())
	assert.True(t, propertySchema{Type: []any{"string", "null"}}.allowsString())
	assert.False(t, propertySchema{Type: "integer"}.allowsString())
	assert.False(t, propertySchema{Type: []any{"integer", "null"}}.allowsString())
	assert.False(t, propertySchema{}.allowsString())
}

func TestWaitForDatabaseState_Cancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	api := newFakeAPI(t, nil)

	err := waitForDatabaseState(ctx, api.client, uuid.MustParse("5c2f1f0e-7d1a-4f7e-9f0b-2a6c1d3e4b5a"), databaseStateRunning)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, api.all())
}

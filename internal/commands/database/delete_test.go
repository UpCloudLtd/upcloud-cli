package database

import (
	"context"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deletedDatabaseUUID = "27fbd082-30b0-11eb-adc1-0242ac120004"

func TestDeleteCommand(t *testing.T) {
	path := "/1.3/database/" + deletedDatabaseUUID
	notFound := fakeResponse{
		status:      http.StatusNotFound,
		contentType: "application/problem+json",
		body:        `{"type": "https://developers.upcloud.com/1.3/errors#ERROR_SERVICE_NOT_FOUND", "title": "Service not found", "status": 404}`,
	}

	for _, test := range []struct {
		name      string
		args      []string
		responses map[string]fakeResponse
		calls     []string
		error     string
		notFound  bool
	}{
		{
			name:      "delete with UUID",
			args:      []string{deletedDatabaseUUID},
			responses: map[string]fakeResponse{"DELETE " + path: {status: http.StatusNoContent}},
			calls:     []string{"DELETE"},
		},
		{
			name: "disable termination protection before deleting",
			args: []string{deletedDatabaseUUID, "--disable-termination-protection"},
			responses: map[string]fakeResponse{
				"PATCH " + path:  {status: http.StatusOK, body: `{"uuid": "` + deletedDatabaseUUID + `", "termination_protection": false}`},
				"DELETE " + path: {status: http.StatusNoContent},
			},
			calls: []string{"PATCH", "DELETE"},
		},
		{
			name: "failed termination protection update prevents deletion",
			args: []string{deletedDatabaseUUID, "--disable-termination-protection"},
			responses: map[string]fakeResponse{
				"PATCH " + path: {status: http.StatusConflict, body: `{"type":"MODIFICATION_FAILED","title":"Could not disable termination protection","status":409}`},
			},
			calls: []string{"PATCH"},
			error: "Could not disable termination protection (type=MODIFICATION_FAILED, status=409)",
		},
		{
			name: "wait until the database is deleted",
			args: []string{deletedDatabaseUUID, "--wait"},
			responses: map[string]fakeResponse{
				"DELETE " + path: {status: http.StatusNoContent},
				"GET " + path:    notFound,
			},
			calls: []string{"DELETE", "GET"},
		},
		{
			name: "wait while database still exists",
			args: []string{deletedDatabaseUUID, "--wait"},
			responses: map[string]fakeResponse{
				"DELETE " + path: {status: http.StatusNoContent},
				"GET " + path: {sequence: []fakeResponse{
					{status: http.StatusOK, body: `{"uuid":"` + deletedDatabaseUUID + `","state":"running"}`},
					notFound,
				}},
			},
			calls: []string{"DELETE", "GET", "GET"},
		},
		{
			name:      "not found is reported as a not found error",
			args:      []string{deletedDatabaseUUID},
			responses: map[string]fakeResponse{"DELETE " + path: notFound},
			calls:     []string{"DELETE"},
			error:     "Service not found (type=https://developers.upcloud.com/1.3/errors#ERROR_SERVICE_NOT_FOUND, status=404)",
			notFound:  true,
		},
		{
			name:  "invalid UUID",
			args:  []string{"not-a-uuid"},
			error: `invalid database UUID "not-a-uuid": invalid UUID length: 10`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newFakeAPI(t, test.responses)
			conf := config.New()
			c := commands.BuildCommand(DeleteCommand(), nil, conf)
			c.Cobra().SetArgs(test.args)

			_, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
				assert.Equal(t, test.notFound, apierror.IsNotFound(err))
			} else {
				require.NoError(t, err)
			}

			var calls []string
			for _, r := range api.all() {
				calls = append(calls, r.Method)
			}
			assert.Equal(t, test.calls, calls)
			if len(test.calls) > 0 && test.calls[0] == "PATCH" {
				assert.JSONEq(t, `{"termination_protection": false}`, string(api.all()[0].Body))
			}
		})
	}
}

func TestWaitUntilDatabaseDeleted_Cancelled(t *testing.T) {
	api := newFakeAPI(t, nil)
	conf := config.New()
	exec := commands.NewExecutor(conf, &smock.Service{}, conf.NewLogger("test")).WithV9Client(api.client)
	conf.Cancel()

	err := waitUntilDatabaseDeleted(exec, api.client, uuid.MustParse(deletedDatabaseUUID))

	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, api.all())
}

package databaseindex

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteManagedDatabaseIndex(t *testing.T) {
	serviceUUID := "898c4cf0-524c-4fc1-9c47-8cc697ed2d52"

	for _, test := range []struct {
		name     string
		args     []string
		expected request.DeleteManagedDatabaseIndexRequest
		errorMsg string
	}{
		{
			name:     "no args",
			args:     []string{serviceUUID},
			errorMsg: `required flag(s) "name" not set`,
		},
		{
			name: "delete success",
			args: []string{
				serviceUUID,
				"--name", ".index-to-delete",
			},
			expected: request.DeleteManagedDatabaseIndexRequest{
				ServiceUUID: serviceUUID,
				IndexName:   ".index-to-delete",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := DeleteCommand()
			mService := new(smock.Service)

			mService.On("DeleteManagedDatabaseIndex", &test.expected).Return(nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "DeleteManagedDatabaseIndex", 1)
			}
		})
	}
}

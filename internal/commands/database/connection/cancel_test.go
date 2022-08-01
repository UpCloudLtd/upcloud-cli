package databaseconnection

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CancelManagedDatabaseConnection"
	uuid := "0fa980c4-0e4f-460b-9869-11b7bd62b833"
	for _, test := range []struct {
		name     string
		args     []string
		error    string
		expected request.CancelManagedDatabaseConnection
	}{
		{
			name:  "no process id",
			args:  []string{},
			error: `required flag(s) "pid" not set`,
		},
		{
			name: "soft cancel",
			args: []string{"--pid", "123456"},
			expected: request.CancelManagedDatabaseConnection{
				UUID:      uuid,
				Pid:       123456,
				Terminate: false,
			},
		},
		{
			name: "terminate",
			args: []string{"--pid", "987654", "--terminate"},
			expected: request.CancelManagedDatabaseConnection{
				UUID:      uuid,
				Pid:       987654,
				Terminate: true,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CancelCommand()
			mService := new(smock.Service)

			mService.On(targetMethod, &test.expected).Return(nil)

			c := commands.BuildCommand(testCmd, nil, config.New())

			c.Cobra().SetArgs(append(test.args, uuid))
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}

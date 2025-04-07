package all

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestPurgeCommand(t *testing.T) {
	for _, test := range []struct {
		name   string
		args   []string
		called [][]interface{}
	}{
		{
			name: "purge non-persistent tf-acc-test resources",
			args: []string{"--include", "*tf-acc-test*", "--exclude", "*persistent*"},
			called: [][]interface{}{
				{"DeleteManagedObjectStorage", &request.DeleteManagedObjectStorageRequest{
					UUID: objectStorages[0].UUID,
				}},
				{"DeleteNetwork", &request.DeleteNetworkRequest{
					UUID: networks.Networks[0].UUID,
				}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			text.DisableColors()

			mService := smock.Service{}
			mockListResponses(&mService)
			for _, call := range test.called {
				method := call[0].(string)
				args := call[1:]
				mService.On(method, args...).Return(nil)
			}

			conf := config.New()
			command := commands.BuildCommand(PurgeCommand(), nil, conf)

			command.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(command, &mService, conf)

			assert.NoError(t, err)
			// We only check expected calls here, as the mock service will error with "Unexpected Method Call" on any unexpected calls.
			for _, call := range test.called {
				method := call[0].(string)
				args := call[1:]
				mService.AssertCalled(t, method, args...)
			}
		})
	}
}

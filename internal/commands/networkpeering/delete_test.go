package networkpeering

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCommand(t *testing.T) {
	targetMethod := "DeleteNetworkPeering"

	peering := upcloud.NetworkPeering{
		Name: "test-peering",
		UUID: "9cb62e7d-e95f-4eaa-9c8b-9c6f5e2a66db",
	}

	for _, test := range []struct {
		name  string
		arg   string
		error string
		req   request.DeleteNetworkPeeringRequest
	}{
		{
			name: "delete with UUID",
			arg:  peering.UUID,
			req:  request.DeleteNetworkPeeringRequest{UUID: peering.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.req
			mService.On(targetMethod, &req).Return(nil)

			conf := config.New()
			command := commands.BuildCommand(DeleteCommand(), nil, conf)

			command.Cobra().SetArgs([]string{test.arg})
			_, err := mockexecute.MockExecute(command, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}

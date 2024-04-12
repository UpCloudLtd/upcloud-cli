package networkinterface

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
	targetMethod := "DeleteNetworkInterface"

	server := upcloud.Server{UUID: "97fbd082-30b0-11eb-adc1-0242ac120002", Title: "test-server"}
	servers := upcloud.Servers{Servers: []upcloud.Server{server}}

	for _, test := range []struct {
		name  string
		arg   string
		flags []string
		error string
		req   request.DeleteNetworkInterfaceRequest
	}{
		{
			name:  "delete network interface with UUID",
			arg:   server.UUID,
			flags: []string{"--index", "4"},
			req:   request.DeleteNetworkInterfaceRequest{ServerUUID: server.UUID, Index: 4},
		},
		{
			name:  "index is missing",
			arg:   server.UUID,
			error: `required flag(s) "index" not set`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.req
			mService.On(targetMethod, &req).Return(nil)

			mService.On("GetServers").Return(&servers, nil)
			conf := config.New()

			c := commands.BuildCommand(DeleteCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.flags, test.arg))
			_, err := mockexecute.MockExecute(c, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}

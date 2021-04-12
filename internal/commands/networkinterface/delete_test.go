package networkinterface

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
			name:  "server is missing",
			flags: []string{"--index", "4"},
			error: "at least one server uuid is required",
		},
		{
			name:  "index is missing",
			arg:   server.UUID,
			error: "index is required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(nil)

			mService.On("GetServers").Return(&servers, nil)
			conf := config.New()

			c := commands.BuildCommand(DeleteCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}

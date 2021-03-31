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

	s := upcloud.Server{UUID: "97fbd082-30b0-11eb-adc1-0242ac120002", Title: "test-server"}
	servers := upcloud.Servers{Servers: []upcloud.Server{s}}

	for _, test := range []struct {
		name  string
		args  []string
		flags []string
		error string
		req   request.DeleteNetworkInterfaceRequest
	}{
		{
			name:  "delete network interface with UUID",
			args:  []string{s.UUID},
			flags: []string{"--index", "4"},
			req:   request.DeleteNetworkInterfaceRequest{ServerUUID: s.UUID, Index: 4},
		},
		{
			name:  "delete network interface with title",
			args:  []string{s.Title},
			flags: []string{"--index", "4"},
			req:   request.DeleteNetworkInterfaceRequest{ServerUUID: s.UUID, Index: 4},
		},
		{
			name:  "server is missing",
			flags: []string{"--index", "4"},
			error: "at least one server uuid is required",
		},
		{
			name:  "index is missing",
			args:  []string{s.UUID},
			error: "index is required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(nil)

			mService.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(DeleteCommand(&mService, &mService), nil, config.New())
			err := c.SetFlags(test.flags)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()(test.args)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}

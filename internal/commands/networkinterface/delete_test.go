package networkinterface

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteCommand(t *testing.T) {
	methodName := "DeleteNetworkInterface"

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
			server.CachedServers = nil
			mns := MockNetworkService{}
			mns.On(methodName, &test.req).Return(nil)

			mss := server.MockServerService{}
			mss.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(DeleteCommand(&mns, &mss), nil, config.New(viper.New()))
			c.SetFlags(test.flags)

			_, err := c.MakeExecuteCommand()(test.args)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mns.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}

}

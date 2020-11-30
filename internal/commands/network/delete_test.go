package network

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
	methodName := "DeleteNetwork"

	n := upcloud.Network{UUID: "0a30b5ca-d0e3-4f7c-81d0-f77d42ea6366", Name: "test-network"}

	for _, test := range []struct {
		name  string
		args  []string
		flags []string
		error string
		req   request.DeleteNetworkRequest
	}{
		{
			name: "delete network with UUID",
			args: []string{n.UUID},
			req:  request.DeleteNetworkRequest{UUID: n.UUID},
		},
		{
			name: "delete network with name",
			args: []string{n.Name},
			req:  request.DeleteNetworkRequest{UUID: n.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedNetworks = nil
			server.CachedServers = nil
			mns := MockNetworkService{}
			mns.On(methodName, &test.req).Return(nil)
			mns.On("GetNetworks").Return(&upcloud.Networks{Networks: []upcloud.Network{n}}, nil)

			c := commands.BuildCommand(DeleteCommand(&mns), nil, config.New(viper.New()))
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
